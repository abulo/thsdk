/*
 * thsdk native bridge — 在专用非 Go 线程中完成 dlopen 与全部 CallAlloc/FreeResult 调用。
 *
 * 背景：原生运行组件（hq.dylib / hq.so / hq.dll）是 Go 语言构建的 c-shared 库。
 * 若在 Go 进程的普通线程上 dlopen（触发其构造函数初始化第二个 Go 运行时）
 * 或直接调用其导出函数，会与宿主 Go 运行时的线程本地存储（TLS）冲突，
 * 导致宿主堆损坏（bad sweepgen / SIGSEGV）。因此这里用一条独立的原生线程
 * 承载全部原生调用，宿主只与这条线程交换 C 字符串。
 */

#include <stdlib.h>
#include <string.h>

/* 原生组件两个导出符号的函数指针类型。 */
typedef void* (*callalloc_t)(const char*);
typedef void (*freeresult_t)(void*);

/* 全局状态：组件句柄与两个导出符号地址（由 worker 线程初始化）。 */
static void* g_handle = NULL;
static callalloc_t g_callalloc = NULL;
static freeresult_t g_freeresult = NULL;

#ifdef _WIN32
#ifndef _WIN32_WINNT
#define _WIN32_WINNT 0x0600
#endif
#include <windows.h>

/* 线程同步：CRITICAL_SECTION + 两个事件。
 * 说明：CONDITION_VARIABLE 在部分 mingw-w64 SDK（如 win32 threads 版本）
 * 中缺失，这里使用事件（Event）实现作业-完成同步，兼容性最好。 */
static CRITICAL_SECTION g_lock;
static HANDLE g_job_event = NULL;   /* 有作业待处理时置位 */
static HANDLE g_done_event = NULL;  /* 作业完成时置位 */
static const char* g_path = NULL;   /* 组件路径（bridge_init 时 strdup） */
static char* g_request = NULL;      /* 待处理请求（strdup 副本） */
static char* g_response = NULL;     /* 完成后的响应（strdup 副本） */
static int g_job_pending = 0;       /* 是否有未处理作业 */
static int g_shutdown = 0;          /* 是否要求线程退出 */

/* worker：专用原生线程主循环（Windows 分支）。
 * 职责：加载组件 → 解析符号 → 循环消费作业：取请求 → CallAlloc →
 * strdup 响应 → FreeResult → 置位完成事件。 */
static DWORD WINAPI worker(LPVOID arg) {
    (void)arg;
    /* 第一步：加载组件并解析导出符号（只执行一次）。 */
    g_handle = LoadLibraryA(g_path);
    if (!g_handle) return 1;
    g_callalloc = (callalloc_t)(void*)GetProcAddress(g_handle, "CallAlloc");
    g_freeresult = (freeresult_t)(void*)GetProcAddress(g_handle, "FreeResult");
    if (!g_callalloc || !g_freeresult) return 2;
    /* 第二步：循环处理作业。 */
    for (;;) {
        /* 等待作业事件置位。 */
        WaitForSingleObject(g_job_event, INFINITE);
        EnterCriticalSection(&g_lock);
        /* 退出信号优先处理。 */
        if (g_shutdown) { LeaveCriticalSection(&g_lock); break; }
        /* 取出请求并清除作业标记。 */
        char* req = g_request;
        g_job_pending = 0;
        LeaveCriticalSection(&g_lock);

        /* 在锁外执行原生调用（耗时操作不阻塞其他桥调用）。 */
        char* resp = NULL;
        if (req) {
            void* p = g_callalloc(req);
            if (p) {
                /* strdup 复制响应，随后立即释放原生侧返回对象。 */
                resp = _strdup((char*)p);
                g_freeresult(p);
            }
        }

        /* 写回响应并置位完成事件唤醒调用方。 */
        EnterCriticalSection(&g_lock);
        g_response = resp;
        LeaveCriticalSection(&g_lock);
        SetEvent(g_done_event);
    }
    return 0;
}

#else
#include <dlfcn.h>
#include <pthread.h>

/* 线程同步：互斥锁 + 两个条件变量（作业条件 / 完成条件）。 */
static pthread_mutex_t g_lock = PTHREAD_MUTEX_INITIALIZER;
static pthread_cond_t g_job_cond = PTHREAD_COND_INITIALIZER;
static pthread_cond_t g_done_cond = PTHREAD_COND_INITIALIZER;
static const char* g_path = NULL;   /* 组件路径（bridge_init 时 strdup） */
static char* g_request = NULL;      /* 待处理请求（strdup 副本） */
static char* g_response = NULL;     /* 完成后的响应（strdup 副本） */
static int g_job_pending = 0;       /* 是否有未处理作业 */
static int g_job_done = 0;          /* 当前作业是否已完成 */
static int g_shutdown = 0;          /* 是否要求线程退出 */

/* worker：专用原生线程主循环（POSIX 分支）。
 * 职责：加载组件 → 解析符号 → 循环消费作业：取请求 → CallAlloc →
 * strdup 响应 → FreeResult → 唤醒等待调用方。 */
static void* worker(void* arg) {
    (void)arg;
    /* 第一步：加载组件并解析导出符号（只执行一次）。 */
    g_handle = dlopen(g_path, RTLD_NOW);
    if (!g_handle) return (void*)1;
    g_callalloc = (callalloc_t)dlsym(g_handle, "CallAlloc");
    g_freeresult = (freeresult_t)dlsym(g_handle, "FreeResult");
    if (!g_callalloc || !g_freeresult) return (void*)2;
    /* 第二步：循环处理作业。 */
    for (;;) {
        /* 等待有作业或退出信号。 */
        pthread_mutex_lock(&g_lock);
        while (!g_job_pending && !g_shutdown) pthread_cond_wait(&g_job_cond, &g_lock);
        if (g_shutdown) { pthread_mutex_unlock(&g_lock); break; }
        /* 取出请求并清除作业标记。 */
        char* req = g_request;
        g_job_pending = 0;
        pthread_mutex_unlock(&g_lock);

        /* 在锁外执行原生调用。 */
        char* resp = NULL;
        if (req) {
            void* p = g_callalloc(req);
            if (p) {
                /* strdup 复制响应，随后立即释放原生侧返回对象。 */
                resp = strdup((char*)p);
                g_freeresult(p);
            }
        }

        /* 写回响应并唤醒调用方。 */
        pthread_mutex_lock(&g_lock);
        g_response = resp;
        g_job_done = 1;
        pthread_cond_signal(&g_done_cond);
        pthread_mutex_unlock(&g_lock);
    }
    return NULL;
}
#endif

/* bridge_init 启动原生线程并完成 dlopen/dlsym。
 * 入参：path 为组件绝对路径。
 * 出参：0 成功；3/4 为初始化失败。
 * 说明：调用方需要保证进程生命周期内只调用一次。 */
int bridge_init(const char* path) {
#ifdef _WIN32
    /* 初始化同步原语与事件句柄。 */
    InitializeCriticalSection(&g_lock);
    g_job_event = CreateEvent(NULL, FALSE, FALSE, NULL);
    g_done_event = CreateEvent(NULL, FALSE, FALSE, NULL);
    if (!g_job_event || !g_done_event) return 3;
    /* 复制组件路径并创建专用线程。 */
    g_path = _strdup(path);
    HANDLE t = CreateThread(NULL, 0, worker, NULL, 0, NULL);
    if (!t) return 4;
    CloseHandle(t);
    return 0;
#else
    /* 复制组件路径并创建专用线程。 */
    g_path = strdup(path);
    pthread_t tid;
    if (pthread_create(&tid, NULL, worker, NULL) != 0) return 3;
    return 0;
#endif
}

/* bridge_call 提交一次请求并等待响应。
 * 入参：request 为 JSON 请求字符串。
 * 出参：strdup 复制出的响应字符串（调用方负责 bridge_free_string 释放），
 *       失败时返回 NULL。
 * 说明：串行执行——同一时刻只处理一个请求；调用方需要保证
 *       bridge_init 已成功。 */
char* bridge_call(const char* request) {
    char* resp = NULL;
#ifdef _WIN32
    /* 写请求并唤醒 worker（作业-完成事件模型）。 */
    EnterCriticalSection(&g_lock);
    if (g_request) { free(g_request); g_request = NULL; }
    g_request = _strdup(request);
    g_job_pending = 1;
    LeaveCriticalSection(&g_lock);
    SetEvent(g_job_event);
    /* 等待完成事件。 */
    WaitForSingleObject(g_done_event, INFINITE);
    /* 取出响应。 */
    EnterCriticalSection(&g_lock);
    resp = g_response;
    g_response = NULL;
    LeaveCriticalSection(&g_lock);
#else
    /* 写请求并唤醒 worker（互斥锁+条件变量模型）。 */
    pthread_mutex_lock(&g_lock);
    /* 等待上一个作业被消费完（防止覆盖未处理请求）。 */
    while (g_job_pending) pthread_cond_wait(&g_job_cond, &g_lock);
    if (g_request) { free(g_request); g_request = NULL; }
    g_request = strdup(request);
    g_job_pending = 1;
    g_job_done = 0;
    pthread_cond_signal(&g_job_cond);
    /* 等待完成信号。 */
    while (!g_job_done) pthread_cond_wait(&g_done_cond, &g_lock);
    resp = g_response;
    g_response = NULL;
    pthread_mutex_unlock(&g_lock);
#endif
    return resp;
}

/* bridge_free_string 释放 bridge_call 返回的字符串。
 * 入参：s 为 bridge_call 的返回值；传 NULL 安全。 */
void bridge_free_string(char* s) {
    if (s) free(s);
}

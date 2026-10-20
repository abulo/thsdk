/*
 * thsdk native bridge — 在专用非 Go 线程中完成 dlopen 与全部 CallAlloc/FreeResult 调用。
 *
 * 背景：hq.dylib / hq.so / hq.dll 是 Go 语言构建的 c-shared 库。
 * 若在 Go 进程的普通线程上 dlopen（触发其构造函数初始化第二个 Go 运行时）
 * 或直接调用其导出函数，会与宿主 Go 运行时的线程本地存储（TLS）冲突，
 * 导致宿主堆损坏（bad sweepgen / SIGSEGV）。因此这里用一条独立的原生线程
 * 承载全部原生调用，宿主只与这条线程交换 C 字符串。
 */

#include <stdlib.h>
#include <string.h>

typedef void* (*callalloc_t)(const char*);
typedef void (*freeresult_t)(void*);

static void* g_handle = NULL;
static callalloc_t g_callalloc = NULL;
static freeresult_t g_freeresult = NULL;

#ifdef _WIN32
#ifndef _WIN32_WINNT
#define _WIN32_WINNT 0x0600
#endif
#include <windows.h>

/* 线程同步：CRITICAL_SECTION + 两个事件。CONDITION_VARIABLE 在部分
 * mingw-w64 SDK（如 win32 threads 版本）中缺失，事件兼容性最好。 */
static CRITICAL_SECTION g_lock;
static HANDLE g_job_event = NULL;
static HANDLE g_done_event = NULL;
static const char* g_path = NULL;
static char* g_request = NULL;
static char* g_response = NULL;
static int g_job_pending = 0;
static int g_shutdown = 0;

static DWORD WINAPI worker(LPVOID arg) {
    (void)arg;
    g_handle = LoadLibraryA(g_path);
    if (!g_handle) return 1;
    g_callalloc = (callalloc_t)(void*)GetProcAddress(g_handle, "CallAlloc");
    g_freeresult = (freeresult_t)(void*)GetProcAddress(g_handle, "FreeResult");
    if (!g_callalloc || !g_freeresult) return 2;
    for (;;) {
        WaitForSingleObject(g_job_event, INFINITE);
        EnterCriticalSection(&g_lock);
        if (g_shutdown) { LeaveCriticalSection(&g_lock); break; }
        char* req = g_request;
        g_job_pending = 0;
        LeaveCriticalSection(&g_lock);

        char* resp = NULL;
        if (req) {
            void* p = g_callalloc(req);
            if (p) {
                resp = _strdup((char*)p);
                g_freeresult(p);
            }
        }

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

static pthread_mutex_t g_lock = PTHREAD_MUTEX_INITIALIZER;
static pthread_cond_t g_job_cond = PTHREAD_COND_INITIALIZER;
static pthread_cond_t g_done_cond = PTHREAD_COND_INITIALIZER;
static const char* g_path = NULL;
static char* g_request = NULL;
static char* g_response = NULL;
static int g_job_pending = 0;
static int g_job_done = 0;
static int g_shutdown = 0;

static void* worker(void* arg) {
    (void)arg;
    g_handle = dlopen(g_path, RTLD_NOW);
    if (!g_handle) return (void*)1;
    g_callalloc = (callalloc_t)dlsym(g_handle, "CallAlloc");
    g_freeresult = (freeresult_t)dlsym(g_handle, "FreeResult");
    if (!g_callalloc || !g_freeresult) return (void*)2;
    for (;;) {
        pthread_mutex_lock(&g_lock);
        while (!g_job_pending && !g_shutdown) pthread_cond_wait(&g_job_cond, &g_lock);
        if (g_shutdown) { pthread_mutex_unlock(&g_lock); break; }
        char* req = g_request;
        g_job_pending = 0;
        pthread_mutex_unlock(&g_lock);

        char* resp = NULL;
        if (req) {
            void* p = g_callalloc(req);
            if (p) {
                resp = strdup((char*)p);
                g_freeresult(p);
            }
        }

        pthread_mutex_lock(&g_lock);
        g_response = resp;
        g_job_done = 1;
        pthread_cond_signal(&g_done_cond);
        pthread_mutex_unlock(&g_lock);
    }
    return NULL;
}
#endif

/* bridge_init 启动原生线程并完成 dlopen/dlsym；0 成功，非 0 失败。 */
int bridge_init(const char* path) {
#ifdef _WIN32
    InitializeCriticalSection(&g_lock);
    g_job_event = CreateEvent(NULL, FALSE, FALSE, NULL);
    g_done_event = CreateEvent(NULL, FALSE, FALSE, NULL);
    if (!g_job_event || !g_done_event) return 3;
    g_path = _strdup(path);
    HANDLE t = CreateThread(NULL, 0, worker, NULL, 0, NULL);
    if (!t) return 4;
    CloseHandle(t);
    return 0;
#else
    g_path = strdup(path);
    pthread_t tid;
    if (pthread_create(&tid, NULL, worker, NULL) != 0) return 3;
    return 0;
#endif
}

/* bridge_call 提交一次请求并等待响应；返回 strdup 的响应字符串（调用方负责释放），失败返回 NULL。 */
char* bridge_call(const char* request) {
    char* resp = NULL;
#ifdef _WIN32
    EnterCriticalSection(&g_lock);
    if (g_request) { free(g_request); g_request = NULL; }
    g_request = _strdup(request);
    g_job_pending = 1;
    LeaveCriticalSection(&g_lock);
    SetEvent(g_job_event);
    WaitForSingleObject(g_done_event, INFINITE);
    EnterCriticalSection(&g_lock);
    resp = g_response;
    g_response = NULL;
    LeaveCriticalSection(&g_lock);
#else
    pthread_mutex_lock(&g_lock);
    while (g_job_pending) pthread_cond_wait(&g_job_cond, &g_lock);
    if (g_request) { free(g_request); g_request = NULL; }
    g_request = strdup(request);
    g_job_pending = 1;
    g_job_done = 0;
    pthread_cond_signal(&g_job_cond);
    while (!g_job_done) pthread_cond_wait(&g_done_cond, &g_lock);
    resp = g_response;
    g_response = NULL;
    pthread_mutex_unlock(&g_lock);
#endif
    return resp;
}

/* bridge_free_string 释放 bridge_call 返回的字符串。 */
void bridge_free_string(char* s) {
    if (s) free(s);
}

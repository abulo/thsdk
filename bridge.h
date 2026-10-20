#ifndef THSDK_BRIDGE_H
#define THSDK_BRIDGE_H

/* bridge_init 启动原生线程并完成 dlopen/dlsym；0 成功，非 0 失败。 */
int bridge_init(const char* path);

/* bridge_call 提交一次请求并等待响应；返回 strdup 的响应字符串（调用方负责释放），失败返回 NULL。 */
char* bridge_call(const char* request);

/* bridge_free_string 释放 bridge_call 返回的字符串。 */
void bridge_free_string(char* s);

#endif /* THSDK_BRIDGE_H */

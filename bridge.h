#ifndef THSDK_BRIDGE_H
#define THSDK_BRIDGE_H

/* bridge_init 启动原生线程并完成 dlopen/dlsym。
 * 入参：path 为运行组件绝对路径（C 字符串）。
 * 出参：0 表示成功，非 0 表示失败。
 * 说明：必须在首次 bridge_call 之前调用且只调用一次。 */
int bridge_init(const char* path);

/* bridge_call 提交一次请求并等待响应。
 * 入参：request 为 JSON 请求字符串（C 字符串，可复用）。
 * 出参：strdup 复制出的响应字符串（调用方负责用 bridge_free_string
 *       释放）；失败时返回 NULL。 */
char* bridge_call(const char* request);

/* bridge_free_string 释放 bridge_call 返回的字符串。
 * 入参：s 为 bridge_call 的返回值；传 NULL 安全。 */
void bridge_free_string(char* s);

#endif /* THSDK_BRIDGE_H */

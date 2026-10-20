@echo off
rem thsdk-go Windows 真机一键构建与测试脚本（在 Windows 上运行）。
rem 要求：Go 1.21+（含 C 编译器，安装 Go 时勾选或配置 mingw-w64）。
rem 用法：双击或在 cmd/PowerShell 中执行 build_all_windows.bat

setlocal
cd /d "%~dp0\.."
echo ==^> go build ./... (本机)
go build ./... || exit /b 1

echo ==^> go vet ./...
go vet ./... || exit /b 1

echo ==^> go test -count=1 .
go test -count=1 . || exit /b 1

echo ==^> 构建示例
go build -o smoke.exe ./examples/smoke || exit /b 1
go build -o demo.exe ./examples/demo || exit /b 1

echo ==^> 运行组件握手冒烟（加载 libs\windows\amd64\hq.dll）
smoke.exe || exit /b 1

echo.
echo 全部通过。可运行 demo.exe 体验真实行情/认证接口。
endlocal

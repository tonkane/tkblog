@echo off
SETLOCAL EnableDelayedExpansion

:: 设置版本包路径
SET VERSION_PACKAGE=github.com/tkane/tkblog/pkg/version

:: 获取最新的标签版本号
for /f "tokens=1" %%i in ('git describe --tags --always --match=v*') do SET VERSION=%%i

echo version is : %VERSION%

:: 检查工作目录状态是否干净
git status --porcelain 2>nul | findstr /r /c:"^ M" >nul || SET GIT_TREE_STATE=dirty
if not exist "%~dp0git status --porcelain 2>nul" SET GIT_TREE_STATE=clean

:: 获取当前的 Git 提交哈希值
for /f "tokens=1" %%i in ('git rev-parse HEAD') do SET GIT_COMMIT=%%i

:: 获取当前的 UTC 时间并格式化
:: 使用PowerShell获取UTC时间并格式化为ISO 8601格式
FOR /F "delims=" %%a IN ('powershell -Command "(Get-Date -UFormat "%%Y-%%m-%%dT%%H:%%M:%%SZ").ToString()"') DO (
    SET BUILD_DATE=%%a
)

:: 显示BUILD_DATE环境变量的值
echo The BUILD_DATE is: %BUILD_DATE%

:: 构建 GO_LDFLAGS 字符串，用于设置编译时的版本信息
SET GO_LDFLAGS=-X %VERSION_PACKAGE%.GitVersion=%VERSION% -X %VERSION_PACKAGE%.GitCommit=%GIT_COMMIT% -X %VERSION_PACKAGE%.GitTreeState=%GIT_TREE_STATE% -X %VERSION_PACKAGE%.BuildDate=%BUILD_DATE%

:: 打印 GO_LDFLAGS，用于调试
echo %GO_LDFLAGS%

:: 使用指定的 -ldflags 编译 Go 程序
go build -v -ldflags "%GO_LDFLAGS%" -o _output/tkblog.exe cmd/tkblog/main.go

:: 结束脚本
ENDLOCAL
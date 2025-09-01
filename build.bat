@echo off
REM LDAPS连接工具Windows编译脚本

echo 开始编译LDAPS连接工具...

REM 设置程序名称
set APP_NAME=ldaps-tool

REM 创建输出目录
if not exist build mkdir build

REM 编译Windows 64位版本
echo 正在编译Windows 64位版本...
set GOOS=windows
set GOARCH=amd64
go build -o build\%APP_NAME%-windows-amd64.exe main.go
if %errorlevel% equ 0 (
    echo ✅ Windows 64位版本编译成功: build\%APP_NAME%-windows-amd64.exe
) else (
    echo ❌ Windows 64位版本编译失败
    exit /b 1
)

REM 编译Windows 32位版本
echo 正在编译Windows 32位版本...
set GOOS=windows
set GOARCH=386
go build -o build\%APP_NAME%-windows-386.exe main.go
if %errorlevel% equ 0 (
    echo ✅ Windows 32位版本编译成功: build\%APP_NAME%-windows-386.exe
) else (
    echo ❌ Windows 32位版本编译失败
    exit /b 1
)

REM 编译Linux 64位版本
echo 正在编译Linux 64位版本...
set GOOS=linux
set GOARCH=amd64
go build -o build\%APP_NAME%-linux-amd64 main.go
if %errorlevel% equ 0 (
    echo ✅ Linux 64位版本编译成功: build\%APP_NAME%-linux-amd64
) else (
    echo ❌ Linux 64位版本编译失败
    exit /b 1
)

REM 编译macOS 64位版本
echo 正在编译macOS 64位版本...
set GOOS=darwin
set GOARCH=amd64
go build -o build\%APP_NAME%-darwin-amd64 main.go
if %errorlevel% equ 0 (
    echo ✅ macOS 64位版本编译成功: build\%APP_NAME%-darwin-amd64
) else (
    echo ❌ macOS 64位版本编译失败
    exit /b 1
)

REM 编译macOS ARM64版本 (Apple Silicon)
echo 正在编译macOS ARM64版本...
set GOOS=darwin
set GOARCH=arm64
go build -o build\%APP_NAME%-darwin-arm64 main.go
if %errorlevel% equ 0 (
    echo ✅ macOS ARM64版本编译成功: build\%APP_NAME%-darwin-arm64
) else (
    echo ❌ macOS ARM64版本编译失败
    exit /b 1
)

echo.
echo 🎉 所有版本编译完成！
echo.
echo 编译结果:
dir build\
echo.
echo Windows用户使用方法:
echo   build\%APP_NAME%-windows-amd64.exe -server ldap.example.com -basedn "dc=example,dc=com" -username "cn=admin,dc=example,dc=com" -password "yourpassword"
echo.

pause
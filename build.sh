#!/bin/bash

# LDAPS连接工具编译脚本
# 支持多平台编译

echo "开始编译LDAPS连接工具..."

# 设置程序名称
APP_NAME="ldaps-tool"

# 创建输出目录
mkdir -p build

# 编译Windows 64位版本
echo "正在编译Windows 64位版本..."
GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}-windows-amd64.exe main.go
if [ $? -eq 0 ]; then
    echo "✅ Windows 64位版本编译成功: build/${APP_NAME}-windows-amd64.exe"
else
    echo "❌ Windows 64位版本编译失败"
    exit 1
fi

# 编译Windows 32位版本
echo "正在编译Windows 32位版本..."
GOOS=windows GOARCH=386 go build -o build/${APP_NAME}-windows-386.exe main.go
if [ $? -eq 0 ]; then
    echo "✅ Windows 32位版本编译成功: build/${APP_NAME}-windows-386.exe"
else
    echo "❌ Windows 32位版本编译失败"
    exit 1
fi

# 编译Linux 64位版本
echo "正在编译Linux 64位版本..."
GOOS=linux GOARCH=amd64 go build -o build/${APP_NAME}-linux-amd64 main.go
if [ $? -eq 0 ]; then
    echo "✅ Linux 64位版本编译成功: build/${APP_NAME}-linux-amd64"
else
    echo "❌ Linux 64位版本编译失败"
    exit 1
fi

# 编译macOS 64位版本
echo "正在编译macOS 64位版本..."
GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}-darwin-amd64 main.go
if [ $? -eq 0 ]; then
    echo "✅ macOS 64位版本编译成功: build/${APP_NAME}-darwin-amd64"
else
    echo "❌ macOS 64位版本编译失败"
    exit 1
fi

# 编译macOS ARM64版本 (Apple Silicon)
echo "正在编译macOS ARM64版本..."
GOOS=darwin GOARCH=arm64 go build -o build/${APP_NAME}-darwin-arm64 main.go
if [ $? -eq 0 ]; then
    echo "✅ macOS ARM64版本编译成功: build/${APP_NAME}-darwin-arm64"
else
    echo "❌ macOS ARM64版本编译失败"
    exit 1
fi

echo ""
echo "🎉 所有版本编译完成！"
echo ""
echo "编译结果:"
ls -la build/
echo ""
echo "Windows用户使用方法:"
echo "  build/${APP_NAME}-windows-amd64.exe -server ldap.example.com -basedn \"dc=example,dc=com\" -username \"cn=admin,dc=example,dc=com\" -password \"yourpassword\""
echo ""
echo "Linux用户使用方法:"
echo "  ./build/${APP_NAME}-linux-amd64 -server ldap.example.com -basedn \"dc=example,dc=com\" -username \"cn=admin,dc=example,dc=com\" -password \"yourpassword\""
echo ""
echo "macOS用户使用方法:"
echo "  ./build/${APP_NAME}-darwin-amd64 -server ldap.example.com -basedn \"dc=example,dc=com\" -username \"cn=admin,dc=example,dc=com\" -password \"yourpassword\""
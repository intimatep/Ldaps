# Go LDAPS 连接示例

这个项目演示了如何使用Go语言连接LDAPS服务器。

## 功能特性

- ✅ LDAPS安全连接 (TLS加密)
- ✅ 用户认证
- ✅ 用户信息查询
- ✅ 搜索功能
- ✅ 连接池管理
- ✅ 错误处理

## 安装依赖

```bash
go mod tidy
```

## 使用方法

### 命令行参数

程序支持通过命令行参数指定LDAPS连接信息：

```bash
go run main.go -server <服务器> -basedn <基础DN> -username <用户名> -password <密码> [选项]
```

#### 必需参数：
- `-server`: LDAP服务器地址
- `-basedn`: 基础DN (如: dc=example,dc=com)
- `-username`: 用户名 (如: cn=admin,dc=example,dc=com)
- `-password`: 密码

#### 可选参数：
- `-port`: 端口号 (默认: 636)
- `-help`: 显示帮助信息

### 使用示例

```bash
# 基本连接
go run main.go -server ldap.example.com -basedn "dc=example,dc=com" -username "cn=admin,dc=example,dc=com" -password "mypassword"

# 指定端口
go run main.go -server ldap.example.com -port 636 -basedn "dc=example,dc=com" -username "cn=admin,dc=example,dc=com" -password "mypassword"

# 显示帮助
go run main.go -help
```

### 交互式搜索

连接成功后，程序会进入交互式搜索模式，支持以下命令：

```bash
ldap> search (objectClass=person)     # 搜索所有人员
ldap> user john                       # 搜索用户john
ldap> group admin                     # 搜索组admin
ldap> all                            # 列出所有条目
ldap> help                           # 显示帮助
ldap> quit                           # 退出
```

### 搜索示例

```bash
# 搜索所有人员
ldap> search (objectClass=person)

# 搜索包含admin的条目
ldap> search (cn=*admin*)

# 搜索有邮箱的人员
ldap> search (&(objectClass=person)(mail=*))

# 搜索特定用户
ldap> user john

# 搜索特定组
ldap> group developers
```
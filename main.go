package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type LDAPSConfig struct {
	Server   string
	Port     int
	BaseDN   string
	Username string
	Password string
	UseTLS   bool
}

func main() {
	// 命令行参数
	var (
		server   = flag.String("server", "", "LDAP服务器地址 (必需)")
		port     = flag.Int("port", 636, "LDAP端口 (默认636)")
		baseDN   = flag.String("basedn", "", "基础DN (必需)")
		username = flag.String("username", "", "用户名 (必需)")
		password = flag.String("password", "", "密码 (必需)")
		help     = flag.Bool("help", false, "显示帮助信息")
	)
	
	flag.Parse()

	// 显示帮助信息
	if *help {
		showHelp()
		return
	}

	// 验证必需参数
	if *server == "" || *baseDN == "" || *username == "" || *password == "" {
		fmt.Println("错误: 缺少必需参数")
		fmt.Println("使用 -help 查看帮助信息")
		showUsage()
		os.Exit(1)
	}

	// 创建配置
	config := LDAPSConfig{
		Server:   *server,
		Port:     *port,
		BaseDN:   *baseDN,
		Username: *username,
		Password: *password,
		UseTLS:   true,
	}

	fmt.Printf("正在连接LDAPS服务器: %s:%d\n", config.Server, config.Port)
	fmt.Printf("基础DN: %s\n", config.BaseDN)
	fmt.Printf("用户名: %s\n", config.Username)
	fmt.Println("---")

	// 连接LDAPS服务器
	conn, err := connectLDAPS(config)
	if err != nil {
		log.Fatalf("连接LDAPS失败: %v", err)
	}
	defer conn.Close()

	// 启动交互式搜索
	startInteractiveSearch(conn, config.BaseDN)
}

func showHelp() {
	fmt.Println("LDAPS连接工具")
	fmt.Println("用法:")
	fmt.Println("  go run main.go -server <服务器> -basedn <基础DN> -username <用户名> -password <密码> [选项]")
	fmt.Println("")
	fmt.Println("必需参数:")
	fmt.Println("  -server     LDAP服务器地址")
	fmt.Println("  -basedn     基础DN (如: dc=example,dc=com)")
	fmt.Println("  -username   用户名 (如: cn=admin,dc=example,dc=com)")
	fmt.Println("  -password   密码")
	fmt.Println("")
	fmt.Println("可选参数:")
	fmt.Println("  -port       端口号 (默认: 636)")
	fmt.Println("  -help       显示此帮助信息")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  go run main.go -server ldap.example.com -basedn \"dc=example,dc=com\" -username \"cn=admin,dc=example,dc=com\" -password \"mypassword\"")
}

func showUsage() {
	fmt.Println("")
	fmt.Println("示例用法:")
	fmt.Println("go run main.go -server ldap.example.com -basedn \"dc=example,dc=com\" -username \"cn=admin,dc=example,dc=com\" -password \"mypassword\"")
}

// 连接LDAPS服务器
func connectLDAPS(config LDAPSConfig) (*ldap.Conn, error) {
	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // 开发环境设为true，生产环境应设为false
		ServerName:         config.Server,
	}

	// 连接LDAPS
	conn, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", config.Server, config.Port), tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("LDAPS连接失败: %w", err)
	}

	// 绑定认证
	err = conn.Bind(config.Username, config.Password)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("LDAP绑定失败: %w", err)
	}

	fmt.Println("LDAPS连接成功!")
	return conn, nil
}

// 启动交互式搜索
func startInteractiveSearch(conn *ldap.Conn, baseDN string) {
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Println("LDAP搜索工具 - 输入命令进行搜索")
	fmt.Println("可用命令:")
	fmt.Println("  search <过滤器>     - 搜索条目 (如: search (objectClass=person))")
	fmt.Println("  user <用户名>       - 搜索特定用户 (如: user john)")
	fmt.Println("  group <组名>        - 搜索特定组 (如: group admin)")
	fmt.Println("  all                 - 列出所有条目")
	fmt.Println("  help                - 显示帮助")
	fmt.Println("  quit                - 退出")
	fmt.Println("")

	for {
		fmt.Print("ldap> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := strings.ToLower(parts[0])

		switch command {
		case "quit", "exit", "q":
			fmt.Println("再见!")
			return
		case "help", "h":
			showSearchHelp()
		case "search", "s":
			if len(parts) < 2 {
				fmt.Println("用法: search <过滤器>")
				continue
			}
			filter := strings.Join(parts[1:], " ")
			searchWithFilter(conn, baseDN, filter)
		case "user", "u":
			if len(parts) < 2 {
				fmt.Println("用法: user <用户名>")
				continue
			}
			username := parts[1]
			searchUser(conn, baseDN, username)
		case "group", "g":
			if len(parts) < 2 {
				fmt.Println("用法: group <组名>")
				continue
			}
			groupname := parts[1]
			searchGroup(conn, baseDN, groupname)
		case "all", "a":
			searchAll(conn, baseDN)
		default:
			fmt.Printf("未知命令: %s\n", command)
			fmt.Println("输入 'help' 查看可用命令")
		}
	}
}

func showSearchHelp() {
	fmt.Println("搜索命令帮助:")
	fmt.Println("  search (objectClass=person)           - 搜索所有人员")
	fmt.Println("  search (objectClass=organizationalUnit) - 搜索组织单位")
	fmt.Println("  search (cn=*admin*)                   - 搜索包含admin的条目")
	fmt.Println("  search (&(objectClass=person)(mail=*)) - 搜索有邮箱的人员")
	fmt.Println("  user john                             - 搜索用户john")
	fmt.Println("  group admin                           - 搜索组admin")
	fmt.Println("  all                                   - 列出所有条目")
}

// 通用搜索函数
func searchWithFilter(conn *ldap.Conn, baseDN, filter string) {
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1000, // 限制结果数量
		0,
		false,
		filter,
		[]string{"*"}, // 返回所有属性
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 个条目:\n", len(sr.Entries))
	for i, entry := range sr.Entries {
		fmt.Printf("\n[%d] DN: %s\n", i+1, entry.DN)
		for _, attr := range entry.Attributes {
			if len(attr.Values) > 0 {
				fmt.Printf("  %s: %s\n", attr.Name, strings.Join(attr.Values, ", "))
			}
		}
	}
	fmt.Println()
}

// 搜索用户
func searchUser(conn *ldap.Conn, baseDN, username string) {
	filter := fmt.Sprintf("(|(uid=%s)(cn=%s)(sAMAccountName=%s))", 
		ldap.EscapeFilter(username), 
		ldap.EscapeFilter(username),
		ldap.EscapeFilter(username))
	
	fmt.Printf("搜索用户: %s\n", username)
	searchWithFilter(conn, baseDN, filter)
}

// 搜索组
func searchGroup(conn *ldap.Conn, baseDN, groupname string) {
	filter := fmt.Sprintf("(|(cn=%s)(name=%s))", 
		ldap.EscapeFilter(groupname), 
		ldap.EscapeFilter(groupname))
	
	fmt.Printf("搜索组: %s\n", groupname)
	searchWithFilter(conn, baseDN, filter)
}

// 搜索所有条目
func searchAll(conn *ldap.Conn, baseDN string) {
	fmt.Println("列出所有条目 (限制前50个):")
	
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		50, // 限制结果数量
		0,
		false,
		"(objectClass=*)",
		[]string{"cn", "ou", "objectClass", "mail", "uid"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 个条目:\n", len(sr.Entries))
	for i, entry := range sr.Entries {
		fmt.Printf("\n[%d] DN: %s\n", i+1, entry.DN)
		fmt.Printf("  对象类: %s\n", strings.Join(entry.GetAttributeValues("objectClass"), ", "))
		if cn := entry.GetAttributeValue("cn"); cn != "" {
			fmt.Printf("  CN: %s\n", cn)
		}
		if ou := entry.GetAttributeValue("ou"); ou != "" {
			fmt.Printf("  OU: %s\n", ou)
		}
		if mail := entry.GetAttributeValue("mail"); mail != "" {
			fmt.Printf("  邮箱: %s\n", mail)
		}
		if uid := entry.GetAttributeValue("uid"); uid != "" {
			fmt.Printf("  UID: %s\n", uid)
		}
	}
	fmt.Println()
}
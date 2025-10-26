package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	Emby struct {
		ServerURL string `mapstructure:"server_url"`
	} `mapstructure:"emby"`
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Domains []struct {
		Name string `mapstructure:"name"`
		URL  string `mapstructure:"url"`
	} `mapstructure:"domains"`
}

type ServerDomainsResponse struct {
	Data []DomainInfo `json:"data"`
	OK   bool         `json:"ok"`
}

type DomainInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var config Config

func main() {
	// 加载配置
	if err := loadConfig(); err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 设置 Gin 为生产模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// 注册路由
	r.GET("/emby/System/Ext/ServerDomains", handleServerDomains)

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	if err := r.Run(addr); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&config)
}

func handleServerDomains(c *gin.Context) {
	// 提取 token
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token not found",
			"ok":    false,
		})
		return
	}

	// 验证 token
	if !validateToken(token) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"ok":    false,
		})
		return
	}

	// 构造返回数据
	domains := make([]DomainInfo, len(config.Domains))
	for i, d := range config.Domains {
		domains[i] = DomainInfo{
			Name: d.Name,
			URL:  d.URL,
		}
	}

	response := ServerDomainsResponse{
		Data: domains,
		OK:   true,
	}

	c.JSON(http.StatusOK, response)
}

func validateToken(token string) bool {
	// 构造验证 URL
	url := fmt.Sprintf("%s/emby/System/Info?X-Emby-Token=%s", config.Emby.ServerURL, token)

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Request creation failed: %v\n", err)
		return false
	}

	// 添加浏览器相关的请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "*/*")

	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// 检查状态码
	// fmt.Printf("Response status: %d, URL: %s\n", resp.StatusCode, url)
	return resp.StatusCode == http.StatusOK
}

func extractToken(c *gin.Context) string {
	if t := c.Query("X-Emby-Token"); t != "" {
		return t
	}
	if t := c.GetHeader("X-Emby-Token"); t != "" {
		return t
	}
	if t := c.Query("api_key"); t != "" {
		return t
	}
	if t := c.GetHeader("X-Emby-Authorization"); t != "" {
		content := strings.TrimPrefix(t, "MediaBrowser ")
		return getTokenByStringSplit(content)
	}
	if t, err := c.Cookie("Authorization"); err == nil && t != "" {
		return t
	}
	if t := c.GetHeader("Authorization"); t != "" {
		return t
	}
	if t := c.Query("token"); t != "" {
		return t
	}
	if t, err := c.Cookie("token"); err == nil && t != "" {
		return t
	}
	if t := c.GetHeader("token"); t != "" {
		return t
	}
	return ""
}

func getTokenByStringSplit(mediaBrowserHeader string) string {
	// 按逗号分割各个字段
	parts := strings.Split(mediaBrowserHeader, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part) // 去除空白字符
		if strings.HasPrefix(part, "Token=") {
			// 去除 Token= 前缀和引号
			token := strings.TrimPrefix(part, "Token=")
			token = strings.Trim(token, `"`) // 去除引号
			return token
		}
	}
	return ""
}

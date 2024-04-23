package config

import (
	"encoding/json"
	"os"

	"github.com/koding/multiconfig"
)

var (
	// C 全局配置(需要先执行MustLoad,否则拿不到配置)
	C = new(Config)
)

func init() {
	loaders := []multiconfig.Loader{
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{},
	}

	loaders = append(loaders, &multiconfig.TOMLLoader{Path: "./config.toml"})

	m := multiconfig.DefaultLoader{
		Loader:    multiconfig.MultiLoader(loaders...),
		Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
	}

	m.MustLoad(C)

	printWithJSON()
}

// IsDebugMode 是否是debug模式
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// PrintWithJSON 基于JSON格式输出配置
func printWithJSON() {
	b, err := json.MarshalIndent(C, "", " ")
	if err != nil {
		os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
		return
	}
	os.Stdout.WriteString(string(b) + "\n")
}

// Config 配置参数
type Config struct {
	RunMode     string
	Swagger     bool
	HTTP        HTTP
	CORS        CORS
	GZIP        GZIP
	JWTAuth     JWTAuth
	Redis       Redis
	MINIO       MINIO
	Log         Log
	RateLimiter RateLimiter
}

type HTTP struct {
	Host             string
	Port             int
	CertFile         string
	KeyFile          string
	ShutdownTimeout  int
	MaxContentLength int64
	MaxLoggerLength  int `default:"4096"`
}

// CORS 跨域请求配置参数
type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

// GZIP gzip压缩
type GZIP struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

// Redis redis配置参数
type Redis struct {
	Addr     string
	Password string
}

// JWTAuth 用户认证
type JWTAuth struct {
	Enable        bool
	SigningMethod string
	SigningKey    string
	Expired       int
	Store         string
	FilePath      string
	RedisDB       int
	RedisPrefix   string
}

// MINIO 存储桶
type MINIO struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Debug           bool
	AutoMakeBucket  bool
}

// Log 日志配置参数
type Log struct {
	Level      int
	Format     string
	Output     string
	OutputFile string
}

// RateLimiter 请求频率限制配置参数
type RateLimiter struct {
	Enable  bool
	Count   int64
	RedisDB int
}

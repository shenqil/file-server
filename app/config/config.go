package config

import (
	"encoding/json"
	"github.com/koding/multiconfig"
	"os"
	"strings"
	"sync"
)

var (
	// C 全局配置(需要先执行MustLoad,否则拿不到配置)
	C    = new(Config)
	once sync.Once
)

func init() {
	MustLoad("./configs/config.toml")
	PrintWithJSON()
}

// MustLoad 加载配置
func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}

		m.MustLoad(C)
	})
}

// PrintWithJSON 基于JSON格式输出配置
func PrintWithJSON() {
	b, err := json.MarshalIndent(C, "", " ")
	if err != nil {
		os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
		return
	}
	os.Stdout.WriteString(string(b) + "\n")
}

// Config 配置参数
type Config struct {
	HTTP  HTTP
	Minio Minio
}

// HTTP http配置参数
type HTTP struct {
	Host             string
	Port             string
	CertFile         string
	KeyFile          string
	ShutdownTimeout  int
	MaxContentLength int64
	MaxLoggerLength  int `default:"4096"`
}

// Minio 文件对象存储
type Minio struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Bucket    string
	Region    string
}

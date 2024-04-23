package app

import (
	"context"
	"crypto/tls"
	"fileServer/util/config"
	"fileServer/util/logger"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Init 应用初始化
func Init(ctx context.Context) (func(), error) {
	logger.WithContext(ctx).Printf("服务启动，运行模式：%s，进程号：%d", config.C.RunMode, os.Getpid())

	// 初始化依赖注入器
	injector, injectorCleanFunc, err := BuildInjector()
	if err != nil {
		return nil, err
	}

	// 初始化HTTP服务
	httpServerCleanFunc := InitHTTPServer(ctx, injector.Engine)

	return func() {
		httpServerCleanFunc()
		injectorCleanFunc()
	}, nil
}

// InitHTTPServer 初始化http服务
func InitHTTPServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.WithContext(ctx).Printf("HTTP server is running at %s.", addr)

		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
	}
}

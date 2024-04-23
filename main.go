package main

import (
	"context"
	"fileServer/app"
	"fileServer/util/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// VERSION 版本号
var VERSION = "1.0.0"

func main() {
	logger.SetVersion(VERSION)

	ctx := logger.NewTagContext(context.Background(), "__main__")

	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 应用初始化
	cleanFunc, err := app.Init(ctx)
	if err != nil {
		logger.WithContext(ctx).Errorf("应用初始化失败, err=%s", err.Error())
		return
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("接收到信号[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.WithContext(ctx).Infof("服务退出")
	time.Sleep(time.Second)
	os.Exit(state)
}

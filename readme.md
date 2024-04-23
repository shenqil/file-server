# File Server

## 依赖工具

```
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/google/wire/cmd/wire@latest
```

## 生成`swagger`文档

```bash
swag init 
```
> 注意生成完成之后需要在`r_api.go`中添加`_ "fileServer/docs"`

## 重新生成依赖注入文件

```bash
wire gen ./app
```

## 运行

```bash
go run main.go
```

> 启动成功之后，可在浏览器中输入地址进行访问：[http://127.0.0.1:11888/swagger/index.html](http://127.0.0.1:11888/swagger/index.html)

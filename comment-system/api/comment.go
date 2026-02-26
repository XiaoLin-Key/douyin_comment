// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net/http"

	"comment-system/api/internal/config"
	"comment-system/api/internal/handler"
	"comment-system/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/comment-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 1. 允许跨域 (CORS)
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 2. 静态文件服务：将 /web 映射到本地文件
	// 这样你可以通过 http://localhost:8888/index.html 访问前端
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/",
		Handler: http.FileServer(http.Dir("../web")).ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/:file",
		Handler: http.FileServer(http.Dir("../web")).ServeHTTP,
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

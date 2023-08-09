package main

import (
        "github.com/micro/go-log"

        "github.com/micro/go-web"
        //"zhiliao_web/handler"
        "github.com/gin-gonic/gin"
        "zhiliao_web/router"
        "zhiliao_web/middle_ware"
)

func main() {

        router := gin.Default()
        // 使用全局中间件，跨域请求
        router.Use(middle_ware.CrosMiddleWare)

        // 注册路由组
        all_router.InitRouter(router)

        service := web.NewService(
                web.Name("zhiliao.web.web.zhiliao_web"),
                web.Version("latest"),
                web.Handler(router),
                web.Address(":8081"),
        )

	    // initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }


	    // run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}

package routes

import (
	"go-proj/app/web/controller"
	"go-proj/app/web/middleware"
	"go-proj/healthCheck/ginCheck"
	"go-proj/library/ginMonitor"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebRoute(router *gin.Engine) {
	//访问日志中间件处理
	logWare := &middleware.LogWare{}

	//对所有的请求进行性能监控，一般来说生产环境，可以对指定的接口做性能监控
	router.Use(logWare.Access(), logWare.Recover(), ginMonitor.Monitor())
	//router.Use(logWare.Access(), logWare.Recover())

	router.NoRoute(middleware.NotFoundHandler())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
		})
	})

	router.GET("/check", ginCheck.HealthCheck)

	homeCtrl := &controller.HomeController{}

	//对个别接口进行性能监控
	//router.GET("/index", ginMonitor.Monitor(), homeCtrl.Index)
	router.GET("/index", homeCtrl.Index)

	router.GET("/test", homeCtrl.Test)

	//定义api前缀分组
	v1 := router.Group("/v1")
	// http://localhost:1338/v1/info/123
	v1.GET("/info/:id", homeCtrl.Info)

	//http://localhost:1338/v1/data?id=456
	v1.GET("/data", homeCtrl.GetData)

	v1.GET("/set-data", homeCtrl.SetData)

	router.GET("/long-async", homeCtrl.LongAsync)

	// 测试将http 处理器和处理器函数包装为gin.Handler
	router.GET("/foo", WrapHttpHandler(FooHandler()))
	router.GET("/foo2", WrapHandlerFunc(FooHandlerFunc))
}

// WrapHttpHandler 将http handler包装为gin.HandlerFunc
func WrapHttpHandler(h http.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// FooHandler http.Handler
func FooHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
}

// WrapHandlerFunc 将http handlerFunc处理器函数包装为gin.HandlerFunc
// 由于http.Handler底层是一个interface上面有ServeHTTP方法
// 而http.HandlerFunc实现了http.Handler的ServeHTTP方法,就相当于实现了http.Handler接口
// gin.Context包含了w http.ResponseWriter, r *http.Request
// 所以调用h.ServeHTTP然后包ctx.Writer,ctx.Request传入就可以处理http请求
func WrapHandlerFunc(h http.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// FooHandlerFunc http处理器函数
func FooHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world,123"))
}

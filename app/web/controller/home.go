package controller

import (
	"context"
	"log"
	"time"

	"github.com/daheige/go-proj/app/logic"
	"github.com/daheige/go-proj/conf/grpcconf"
	"github.com/daheige/go-proj/library/logger"

	config "github.com/daheige/go-proj/conf"

	"github.com/gin-gonic/gin"

	pb "github.com/daheige/go-proj/pb"
)

// HomeController home ctrl.
type HomeController struct {
	BaseController
}

// Index index.
func (ctrl *HomeController) Index(ctx *gin.Context) {
	ctx.JSON(HTTPSuccess, gin.H{
		"code":    200,
		"message": "ok",
	})
}

// Test 测试panic
func (ctrl *HomeController) Test(ctx *gin.Context) {
	panic(11)

	ctx.JSON(HTTPSuccess, gin.H{
		"code":    0,
		"message": "ok",
		"data":    "this is test",
	})
}

// Info 接收url 路由参数id
func (ctrl *HomeController) Info(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(HTTPSuccess, gin.H{
		"code":    0,
		"message": "ok",
		"data":    "current id: " + id,
	})
}

// GetData get homeLogic GetData.
func (ctrl *HomeController) GetData(ctx *gin.Context) {
	homeLogic := logic.HomeLogic{}
	homeLogic.SetCtx(ctx.Request.Context())

	data := homeLogic.GetData()

	ctx.JSON(HTTPSuccess, gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

// SetData 模拟redis set data
func (ctrl *HomeController) SetData(ctx *gin.Context) {
	redisObj, err := config.GetRedisObj("default")
	if err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": "redis connection error",
		})

		return
	}

	//用完就需要释放连接，防止过多的连接导致redis连接过多而陷入长久等待，从而redis崩溃
	defer redisObj.Close()

	_, err = redisObj.Do("set", "myname", "daheige")
	if err != nil {
		log.Println("redis error:", err)
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": "set data error",
		})

		return
	}

	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "set data success",
	})
}

// LongAsync When starting new Goroutines inside a middleware or handler,
// you SHOULD NOT use the original context inside it,
// you have to use a read-only copy.
func (ctrl *HomeController) LongAsync(ctx *gin.Context) {
	// create copy to be used inside the goroutine
	cCp := ctx.Copy()
	go func() {
		// simulate a long task with time.Sleep(). 3 seconds
		time.Sleep(3 * time.Second)

		// note that you are using the copied context "cCp", IMPORTANT
		log.Println("Done! in path " + cCp.Request.URL.Path)
	}()

	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "ok",
	})
}

// SayHello 测试http request然后调用grpc service.
// http://localhost:1338/say-hello/daheige
func (ctrl *HomeController) SayHello(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.Param("name")
	if name == "" {
		ctrl.Error(c, 500, "server inner error", nil)
		return
	}

	// 调用hello rpc服务
	helloService := pb.NewGreeterServiceClient(grpcconf.HellServiceClient)

	grcCtx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	res, err := helloService.SayHello(grcCtx, &pb.HelloReq{
		Name: name,
	})

	if err != nil {
		logger.Error(ctx, "grpc call error", map[string]interface{}{
			"trace_error": err.Error(),
		})

		log.Println("grpc call error: ", err)

		ctrl.Error(c, APIError, "service error", EmptyObject{})
		return
	}

	ctrl.Success(c, "ok", res)
}

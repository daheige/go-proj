package controller

import (
	"go-proj/app/logic"
	"log"
	"time"

	config "go-proj/conf"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	BaseController
}

// action
func (ctrl *HomeController) Index(ctx *gin.Context) {
	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    200,
		"message": "ok",
	})
}

func (ctrl *HomeController) Test(ctx *gin.Context) {
	panic(11)

	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data":    "this is test",
	})
}

func (ctrl *HomeController) Info(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data":    "current id: " + id,
	})
}

func (ctrl *HomeController) GetData(ctx *gin.Context) {
	homeLogic := logic.HomeLogic{}
	homeLogic.SetCtx(ctx.Request.Context())

	data := homeLogic.GetData()

	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

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

// When starting new Goroutines inside a middleware or handler,
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

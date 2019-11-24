package job

import (
	"go-proj/library/helper"
	"log"
)

type TestJob struct {
	BaseJob
}

func (j *TestJob) Info() {
	id := helper.GetStringByCtx(j.ctx, "id")
	log.Println("current id: ", id)
}

package job

import (
	"log"

	"github.com/daheige/go-proj/library/helper"
)

type TestJob struct {
	BaseJob
}

func (j *TestJob) Info() {
	id := helper.GetStringByCtx(j.ctx, "id")
	log.Println("current id: ", id)
}

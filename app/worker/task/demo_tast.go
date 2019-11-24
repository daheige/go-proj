package task

import "log"

type TestTask struct {
	BaseTask
}

func (t *TestTask) Hello() {
	log.Println("hello world")
}

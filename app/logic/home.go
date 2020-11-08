package logic

import (
	"log"

	"github.com/daheige/go-proj/library/helper"
)

type HomeLogic struct {
	BaseLogic
}

func (h *HomeLogic) GetData() []string {
	log.Println(helper.GetStringByCtx(h.Ctx, "current_uid"))

	return []string{
		"golang",
		"php",
	}
}

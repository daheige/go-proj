package logic

import (
	"go-proj/library/helper"
	"log"
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

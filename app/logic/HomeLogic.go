package logic

import (
	"go-proj/app/helper"
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

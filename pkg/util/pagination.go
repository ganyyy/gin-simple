package util

import (
	"gin-simple/pkg/settings"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"log"
)

func GetPage(c *gin.Context) int {
	result := 0
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		log.Printf("error query page:%v", err)
		return 0
	}
	if page > 0 {
		result = (page - 1) * settings.PageSize
	}

	return result
}

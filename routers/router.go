package routers

import (
	"gin-simple/pkg/settings"
	v1 "gin-simple/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(settings.RunMode)

	r.GET("/test", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "ok",
		})
	})

	apiV1 := r.Group("/api/v1")
	{
		// Tags api
		apiV1.GET("/tags", v1.GetTags)
		apiV1.POST("/tags", v1.AddTag)
		apiV1.PUT("/tags/:id", v1.EditTag)
		apiV1.DELETE("/tags/:id", v1.DeleteTag)

		// Article api
		apiV1.GET("/articles", v1.GetArticles)
		apiV1.GET("/articles/:id", v1.GetArticle)
		apiV1.POST("/articles", v1.AddArticle)
		apiV1.PUT("/articles/:id", v1.EditArticle)
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}

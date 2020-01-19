package v1

import (
	"gin-simple/models"
	"gin-simple/pkg/err"
	"gin-simple/pkg/settings"
	"gin-simple/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"log"
	"net/http"
)

//获取单个文章
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := err.INVALID_PARAMS
	var data interface{}

	if !valid.HasErrors() {
		if models.ExistArticleById(id) {
			data = models.GetArticle(id)
			code = err.SUCCESS
		} else {
			code = err.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, e := range valid.Errors {
			log.Printf("e.key:%s, e.message:%s", e.Key, e.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": data,
	})
}

//获取多个文章
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state = -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")
	}

	var tagId int
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("TagID最小值是")
	}

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		code = err.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), settings.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		for k, v := range valid.ErrorsMap {
			data[k] = v
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": data,
	})
}

//新增文章
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID大于1")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("描述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只能为0或者1")

	resp := make(map[string]interface{})

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state
			models.AddArticle(data)
			code = err.SUCCESS
		} else {
			code = err.ERROR_NOT_EXIST_TAG
		}
	} else {
		for k, v := range valid.ErrorsMap {
			resp[k] = v
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": resp,
	})
}

//修改文章
func EditArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID不能小于1")
	valid.MaxSize(title, 100, "title").Message("标题长度不大于100")
	valid.MaxSize(desc, 255, "desc").Message("描述最多255")
	valid.MaxSize(content, 65535, "content").Message("内容最长65535")
	valid.Required(modifiedBy, "modified_by").Message("需要有修改人")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长字符100")
	valid.Range(state, 0, 1, "state").Message("状态只能为0或者1")

	resp := make(map[string]interface{})
	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleById(id) {
			if models.ExistTagByID(tagId) {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}
				data["modified_by"] = modifiedBy

				models.EditArticle(id, data)
				code = err.SUCCESS
			} else {
				code = err.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = err.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for k, v := range valid.ErrorsMap {
			resp[k] = v
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": resp,
	})
}

//删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	resp := make(map[string]interface{})
	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleById(id) {
			models.DeleteArticle(id)
			code = err.SUCCESS
		} else {
			code = err.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for k, v := range valid.ErrorsMap {
			resp[k] = v
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": resp,
	})
}

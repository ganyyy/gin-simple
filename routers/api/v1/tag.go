package v1

import (
	"gin-simple/models"
	"gin-simple/pkg/err"
	"gin-simple/pkg/settings"
	"gin-simple/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// 获取所有标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	// 数据查询
	maps := make(map[string]interface{})
	// 结果集
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := err.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(c), settings.PageSize, maps)
	data["total"] = models.GetTagsTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": data,
	})
}

// 添加一个新的标签
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长100字符")
	valid.Range(state, 0, 1, "state").Message("状态只能是0或者1")

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = err.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = err.ERROR_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": valid.ErrorsMap,
	})
}

// 修改一个标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只能0或者1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长100字符")

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		code = err.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		} else {
			code = err.ERROR_NOT_EXIST_TAG
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 删除一个标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := err.INVALID_PARAMS
	if !valid.HasErrors() {
		code = err.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = err.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  err.GetMsg(code),
		"data": make(map[string]string),
	})
}

package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	DeletedOn  int    `json:"deleted_on"`
	State      int    `json:"state"`
}

func (t *Tag) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("CreatedOn", time.Now().Unix())
}

func (t *Tag) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("ModifiedOn", time.Now().Unix())
}

func (t *Tag) BeforeDelete(scope *gorm.Scope) error {
	return scope.SetColumn("DeletedOn", time.Now().Unix())
}

func GetTags(pageNum, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

func GetTagsTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID == id {
		return true
	}
	return false
}

func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})
	return true
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})
	return true
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Update(data)
	return true
}

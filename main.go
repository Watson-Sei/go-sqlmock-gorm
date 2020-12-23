package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Tag struct {
	ID 		string	`gorm:"primary_key"`
	Name 	string
}

// Model
type ModelInterface interface {
	GetAllTag() (*[]Tag, error)
	GetByIdTag(id string) (*Tag, error)
	CreateTag(id, name string) (*Tag, error)
	UpdateTag(id, name string, tag *Tag) (*Tag, error)
	DeleteTag(id string) error
}

type Model struct {
	db *gorm.DB
}

func (m Model) GetAllTag() (*[]Tag, error) {
	var tags []Tag
	err := m.db.Find(&tags).Error
	return &tags, err
}

func (m Model) GetByIdTag(id string) (*Tag, error) {
	var tag Tag
	err := m.db.Where("id = ?", id).Find(&tag).Error
	return &tag, err
}

func (m Model) CreateTag(id, name string) (*Tag, error) {
	tx := m.db.Begin()
	tag := &Tag{
		ID: id,
		Name: name,
	}
	err := tx.Create(tag).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return tag, err
}

func (m Model) UpdateTag(id, name string, tag *Tag) (*Tag, error) {
	tx := m.db.Begin()
	err := tx.Model(&tag).Where("id = ?",id).Update("name", name).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return tag, err
}

func (m Model) DeleteTag(id string) error {
	var tag Tag
	tx := m.db.Begin()
	err := tx.Delete(&tag, id).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// Controller
type Controller struct {
	model ModelInterface
}

func (c Controller) GetAllTag() (*[]Tag, error) {
	tags, err := c.model.GetAllTag()
	return tags, err
}

func (c Controller) GetByIdTag(id string) (*Tag, error) {
	tag, err := c.model.GetByIdTag(id)
	return tag, err
}

func (c Controller) CreateTag(id, name string) (*Tag, error) {
	tag, err := c.model.CreateTag(id, name)
	return tag, err
}

func (c Controller) UpdateTag(id, name string) (*Tag, error) {
	tag, err := c.model.GetByIdTag(id)
	if err != nil {
		return nil, err
	}
	tag, err = c.model.UpdateTag(id, name, tag)
	return tag, err
}

func (c Controller) DeleteTag(id string) error {
	err := c.model.DeleteTag(id)
	return err
}


func main()  {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Tag{})

	m := Model{db: db}

	m.db.Create(&Tag{ID: "1", Name: "Google"})
	m.db.Create(&Tag{ID: "2", Name: "Facebook"})


	// コントローラ定義
	controller := Controller{model: m}

	// 全データを取得する処理
	res, err := controller.GetAllTag()
	fmt.Println(res)
	// 特定のデータを取得する処理
	res1, err := controller.GetByIdTag("1")
	fmt.Println(res1)
	// データを作成する処理
	res2, err := controller.CreateTag("3","Amazon")
	fmt.Println(res2)
	// データを更新する処理
	res3, err := controller.UpdateTag("2","Microsoft")
	fmt.Println(res3)
	// データを削除する処理
	err = controller.DeleteTag("2")
	fmt.Println(err)
}

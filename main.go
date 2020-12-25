package main
import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Tag struct {
	ID 	 string `gorm:"primary_key"`
	Name string
}

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

type Controller struct {
	model ModelInterface
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
	err := tx.Model(&tag).Where("id = ?", id).Update("name", name).Error
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

func (c Controller) GetAllTag() error {
	ret, err := c.model.GetAllTag()
	fmt.Println(ret)
	return err
}

func (c Controller) GetByIdTag(id string) error {
	ret, err := c.model.GetByIdTag(id)
	fmt.Println(ret)
	return err
}

func (c Controller) CreateTag(id, name string) error {
	ret, err := c.model.CreateTag(id, name)
	fmt.Println(ret)
	return err
}

func (c Controller) UpdateTag(id, name string) error {
	tag, err := c.model.GetByIdTag(id)
	if err != nil {
		return err
	}
	ret, err := c.model.UpdateTag(id, name, tag)
	fmt.Println(ret)
	return err
}

func (c Controller) DeleteTag(id string) error {
	err := c.model.DeleteTag(id)
	fmt.Println("data delete success")
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
	controller := Model{db: db}
	mycontroller := Controller{model: controller}
	mycontroller.GetAllTag()
	mycontroller.CreateTag("3", "Amazon")
	mycontroller.GetByIdTag("3")
	mycontroller.UpdateTag("3", "AWS")
	mycontroller.DeleteTag("3")
}
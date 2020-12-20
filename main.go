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

func GetTag(db *gorm.DB, id string) (*Tag, error) {
	var tag Tag
	err := db.Where("id = ?", id).Find(&tag).Error
	return &tag, err
}

func CreateTag(db *gorm.DB, id string, name string) (*Tag, error) {
	tx := db.Begin()
	tag := &Tag{
		ID: id,
		Name: name,
	}
	err := tx.Create(tag).Error
	if err != nil {
		tx.Rollback()
		return tag, err
	}
	tx.Commit()
	return tag, err
}

func main()  {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	db.AutoMigrate(&Tag{})
	
	//db.Create(&Tag{ID: "1", Name: "Google"})
	//db.Create(&Tag{ID: "2", Name: "FaceBook"})

	res, err := CreateTag(db, "1","Google")
	fmt.Println(res)
	res, err = CreateTag(db, "2", "FaceBook")
	fmt.Println(res)

	res, err = GetTag(db, "1")
	fmt.Println(res.ID, res.Name)
}

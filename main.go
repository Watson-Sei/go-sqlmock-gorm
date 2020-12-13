package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID string
	Name string
}

func CreateUser(db *gorm.DB, user *User) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func main()  {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	CreateUser(db, &User{ID: "1", Name: "Taro"})

	var user User
	db.First(&user, 1).Scan(&user)
	fmt.Println(user.ID, user.Name)
}

package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID           uint `gorm:"primaryKey; AUTO_INCREMENT;not null;"`
	Name 		 string `gorm:"unique;not null"`
}

func (b *User) TableName() string {
	return "users"
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

	CreateUser(db, &User{Name: "Yamada"})
	CreateUser(db, &User{Name: "Suzuki"})

	var users []User
	db.Find(&users).Scan(&users)
	fmt.Println(users)

}

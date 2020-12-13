package main

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type v2Suite struct {
	db *gorm.DB
	mock sqlmock.Sqlmock
	user User
}

func TestCreateUser(t *testing.T) {
	s := &v2Suite{}
	var (
		db *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	if db == nil {
		t.Error("mock db is null")
	}

	if s.mock == nil {
		t.Error("sqlmock is null")
	}


	s.db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Errorf("Failed to open gorm v2 db, got error: %v", err)
	}

	if s.db == nil {
		t.Error("gorm db is null")
	}

	s.user = User{
		ID: "2",
		Name: "Taro",
	}

	defer db.Close()

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("id","name") VALUES ($1,$2)`)).
		WithArgs(s.user.ID ,s.user.Name).WillReturnRows(sqlmock.NewRows([]string{"id","name"}).AddRow(s.user.ID, s.user.Name))
	s.mock.ExpectCommit()

	if err = CreateUser(s.db, &User{ID: "2", Name: "Taro"}); err != nil {
		t.Errorf("Failed to insert to gorm db, got error: %v", err)
	}
}
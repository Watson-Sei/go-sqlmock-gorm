package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)


// モックを作成
func GetNewDbMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, mock, err
	}

	gormDB, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{DriverName: "mysql", Conn: db, SkipInitializeWithVersion: true}}, &gorm.Config{})

	if err != nil {
		return gormDB, mock, err
	}

	return gormDB, mock, err
}

func TestGetAllTag(t *testing.T) {
	db, mock, err := GetNewDbMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `tags`")).
		WillReturnRows(sqlmock.NewRows([]string{"id","name"}).
			AddRow("1","Google").
			AddRow("2","FaceBook"))

	res, err := GetAllTag(db)
	if err != nil {
		t.Fatal(err)
	}

	if res == nil {
		t.Errorf("値が取得できていません %v", res)
	}


}

func TestGetTag(t *testing.T) {
	db, mock, err := GetNewDbMock()
	if err != nil {
		t.Errorf("Failed to initialize mock DB: %v", err)
	}

	id := "1"
	name := "Google"
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `tags` WHERE id = ?")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id","name"}).
			AddRow(id, name))

	//
	res, err := GetTag(db, id)
	if err != nil {
		t.Fatal(err)
	}

	if res.ID != id || res.Name != name {
		t.Errorf("取得結果不一致 %+v", res)
	}
}

func TestCreateTag(t *testing.T) {
	db, mock, err := GetNewDbMock()
	if err != nil {
		t.Fatal(err)
	}

	id := "1"
	name := "Google"

	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `tags` (`id`,`name`) VALUES (?,?)")).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(1,1))
	mock.ExpectCommit()

	_, err = CreateTag(db, id, name)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTag(t *testing.T)  {
	db, mock, err := GetNewDbMock()
	if err != nil {
		t.Fatal(err)
	}

	id := "1"

	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM `tags` WHERE `tags`.`id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1,1))
	mock.ExpectCommit()

	if err = DeleteTag(db, id); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
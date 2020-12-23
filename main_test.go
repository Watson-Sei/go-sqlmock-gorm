package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

// tesitfy mock を使った Model の Mock
type ModelMock struct {
	mock.Mock
}

// 難しいが、メソッドをモックするときの書き方
func (m *ModelMock) GetAllTag() (*[]Tag, error) {
	args := m.Called()
	return args.Get(0).(*[]Tag), args.Error(1)
}

func TestControllerGetAllTag (t *testing.T) {
	// Model のモックを作成
	modelMock := new(ModelMock)
	modelMock.On("GetAllTag").Return(&[]Tag{}, nil)

	// モックを入れてテスト
	controller := Controller{model: modelMock}
	tags, err := controller.GetAllTag()


	assert.Nil(t, err)
	assert.Equal(t, tags, &[]Tag{})
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

	m := Model{db: db}
	res, err := m.GetAllTag()

	assert.Nil(t, err)

	// want は期待する結果
	want := &[]Tag {
		{"1", "Google"},
		{"2", "FaceBook"},
	}
	assert.Equal(t, want, res)
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

func TestUpdateTag(t *testing.T) {
	db, mock, err := GetNewDbMock()
	if err != nil {
		t.Fatal(err)
	}

	id := "1"
	name := "Google21"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE `tags`")).
		WillReturnResult(sqlmock.NewResult(1,1))
	mock.ExpectCommit()

	_, err = UpdateTag(db, id, name, &Tag{})
	if err != nil {
		t.Fatal(err)
	}
}

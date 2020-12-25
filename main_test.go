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

// ControllerTest

type ModelMock struct {
	mock.Mock
}

func (m *ModelMock) GetAllTag() (*[]Tag, error) {
	args := m.Called()
	return args.Get(0).(*[]Tag), args.Error(1)
}

func (m *ModelMock) GetByIdTag(id string) (*Tag, error) {
	args := m.Called(id)
	return args.Get(0).(*Tag), args.Error(1)
}

func (m *ModelMock) CreateTag(id, name string) (*Tag, error) {
	args := m.Called(id, name)
	return args.Get(0).(*Tag), args.Error(1)
}

func (m *ModelMock) UpdateTag(id, name string, tag *Tag) (*Tag, error) {
	args := m.Called(id, name, tag)
	return args.Get(0).(*Tag), args.Error(1)
}

func (m *ModelMock) DeleteTag(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestController_GetAllTag(t *testing.T) {
	testTag := &[]Tag{
		{ID: "1", Name: "Google"},
		{ID: "2", Name: "Facebook"},
		{ID: "3", Name: "Amazon"},
	}

	modelMock := new(ModelMock)
	modelMock.On("GetAllTag").Return(testTag, nil)

	controller := Controller{model: modelMock}
	ret, err := controller.model.GetAllTag()

	assert.Nil(t, err)
	assert.Equal(t, testTag, ret)
}

func TestController_GetByIdTag(t *testing.T) {
	testTag := &Tag{ID: "1", Name: "Google"}

	modelMock := new(ModelMock)
	modelMock.On("GetByIdTag", testTag.ID).Return(testTag, nil)

	controller := Controller{model: modelMock}
	ret, err := controller.model.GetByIdTag(testTag.ID)

	assert.Nil(t, err)
	assert.Equal(t, testTag, ret)
}

func TestController_CreateTag(t *testing.T) {
	testTag := &Tag{ID: "1", Name: "Google"}

	modelMock := new(ModelMock)
	modelMock.On("CreateTag", testTag.ID, testTag.Name).Return(testTag, nil)

	controller := Controller{model: modelMock}
	ret, err := controller.model.CreateTag(testTag.ID, testTag.Name)

	assert.Nil(t, err)
	assert.Equal(t, testTag, ret)
}

func TestController_UpdateTag(t *testing.T) {
	testTag := &Tag{ID: "1", Name: "Google"}

	modelMock := new(ModelMock)
	modelMock.On("GetByIdTag", testTag.ID).Return(testTag, nil)
	modelMock.On("UpdateTag", testTag.ID, "Amazon", testTag).Return(&Tag{ID: "1", Name: "Amazon"}, nil)

	controller := Controller{model: modelMock}
	tag, err := controller.model.GetByIdTag(testTag.ID)
	ret, err := controller.model.UpdateTag(testTag.ID, "Amazon", tag)

	assert.Nil(t, err)
	assert.Equal(t, "Amazon", ret.Name)
}

func TestController_DeleteTag(t *testing.T) {
	modelMock := new(ModelMock)
	modelMock.On("DeleteTag", "1").Return(nil)

	controller := Controller{model: modelMock}
	err := controller.model.DeleteTag("1")

	assert.Nil(t, err)
}

// ModelTest
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

func TestModel_GetAllTag(t *testing.T) {
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

func TestModel_GetByIdTag(t *testing.T) {
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

	m := Model{db: db}
	ret, err := m.GetByIdTag(id)

	want := &Tag{ID: id, Name: name}

	assert.Nil(t, err)
	assert.Equal(t, want, ret)
}

func TestModel_CreateTag(t *testing.T) {
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

	m := Model{db: db}
	ret, err := m.CreateTag(id, name)

	want := &Tag{ID: id, Name: name}

	assert.Nil(t, err)
	assert.Equal(t, want, ret)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestModel_UpdateTag(t *testing.T) {
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

	m := Model{db: db}
	ret, err := m.UpdateTag(id, name, &Tag{ID: id, Name: name})

	want := &Tag{ID: id, Name: name}

	assert.Nil(t, err)
	assert.Equal(t, want, ret)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestModel_DeleteTag(t *testing.T) {
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

	m := Model{db: db}
	err = m.DeleteTag(id)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
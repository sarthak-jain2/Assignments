package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sarthak-jain2/go-bookstore/pkg/config"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `gorm:""json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

// function to connect with mysql Database
func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

// function to create a book which takes and returns book
func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

// function which gives all the books in the form of slices
func GetBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getbook Book
	db := db.Where("Id=?", Id).Find(&getbook)
	return &getbook, db
}

func DeleteBook(Id int64) Book {
	var delbook Book
	db.Where("Id=?", Id).Delete(delbook)
	return delbook

}

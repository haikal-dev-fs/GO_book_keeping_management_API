package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// migrate database

	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

}

func CreateBook(c *gin.Context) {
	var book Book
	// bind request body

	if err := c.ShouldBindJSON(&book); err != nil {
		ResponseJson(c, http.StatusBadRequest, "Invalid Input", nil)
		return
	}

	DB.Create(&book)
	ResponseJson(c, http.StatusCreated, "Book Created", book)
}

func GetBooks(c *gin.Context) {
	var books []Book
	DB.Find(&books)
	ResponseJson(c, http.StatusOK, "Books Found", books)
}

func GetBook(c *gin.Context) {
	var book Book
	if err := DB.First(&book, c.Param("id")).Error; err != nil {
		ResponseJson(c, http.StatusNotFound, "Book Not Found", nil)
		return
	}
	ResponseJson(c, http.StatusOK, "Book Retrieved Successfully", book)
}

func UpdateBook(c *gin.Context) {
	var book Book

	if err := DB.First(&book, c.Param("id")).Error; err != nil {
		ResponseJson(c, http.StatusNotFound, "Book Not Found", nil)
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		ResponseJson(c, http.StatusBadRequest, "Invalid Input", nil)
		return
	}

	DB.Save(&book)
	ResponseJson(c, http.StatusOK, "Book Updated Successfully", book)

}

func DeleteBook(c *gin.Context) {
	var book Book
	if err := DB.Delete(&book, c.Param("id")).Error; err != nil {
		ResponseJson(c, http.StatusNotFound, "Book Not Found", nil)
		return
	}

	ResponseJson(c, http.StatusOK, "Book Deleted Successfully", nil)
}

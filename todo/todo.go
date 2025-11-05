package todo

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	Title     string `json:"text" binding:"required"`
	gorm.Model 
}

type TodoHandler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}


// Define the table name for the Todo model
func (Todo) TableName() string {
	return "todos"
} 

func (t *TodoHandler) NewTask(c *gin.Context) {

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
	}	

	if todo.Title == "sleep" {
		transactionId := c.Request.Header.Get("Transaction-Id")
		aud, _ := c.Get("aud")
		log.Println(transactionId, aud, "sleep is not allowed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "sleep is not allowed"})
		return
	}

	r := t.db.Create(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": todo.Model.ID})
	return

}


func (t *TodoHandler) List(c *gin.Context) {
	var todos []Todo
	r := t.db.Find(&todos)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
	return
}

func (t *TodoHandler) Remove(c *gin.Context) {
	idParam := c.Param("id")
	id,err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	r := t.db.Delete(&Todo{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
	return
}
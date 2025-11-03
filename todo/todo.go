package todo

import (
	"net/http"

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

func (Todo) TableName() string {
	return "todos"
} // กำหนดชื่อ table ใน database ด้วยตัวเอง

func (t *TodoHandler) NewTask(c *gin.Context) {

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

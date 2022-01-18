package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"
)

type Todo struct {
	Id        int    `json:"id"`
	Frst_name string `json:"frst_name,omitempty"`
	Last_name string `json:"last_name,omitempty"`
	Address   string `json:"address,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Email     string `json:"email,omitempty validate: email,required"`
}

func CreateTodoTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createError := db.CreateTable(&Todo{}, opts)
	if createError != nil {
		log.Printf("Error while creating todo table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Todo table created")
	return nil
}

var dbConnect *pg.DB

func InitiateDB(db *pg.DB) {
	dbConnect = db
}

func GetAllTodos(c *gin.Context) {
	var todos []Todo
	err := dbConnect.Model(&todos).Select()

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Todos",
		"data":    todos,
	})
	return
}
func CreateTodo(c *gin.Context) {
	var todo = Todo{
		Id:        7,
		Frst_name: "krishna",
		Last_name: "gump",
		Address:   "hyd",
		Email:     "hh@gmail.com",
		Mobile:    "9393009018",
	}

	c.BindJSON(&todo)
	id := todo.Id
	frst_name := todo.Frst_name
	last_name := todo.Last_name
	address := todo.Address
	mobile := todo.Mobile
	email := todo.Email

	insertError := dbConnect.Insert(&Todo{
		Id:        id,
		Frst_name: frst_name,
		Last_name: last_name,
		Address:   address,
		Mobile:    mobile,
		Email:     email,
	})
	if insertError != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Todo created Successfully",
	})
	return
}
func GetSingleTodo(c *gin.Context) {

	c.Param("todoid")
	todo := &Todo{Id: 6}
	err := dbConnect.Select(todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data":    todo,
	})
	return
}
func EditTodo(c *gin.Context) {
	c.Param("todoId")
	var todo = Todo{
		Id:        6,
		Last_name: "gumpalli"}
	c.BindJSON(&todo)
	last_name := todo.Last_name

	_, err := dbConnect.Model(&Todo{}).Set("last_name = ?", last_name).Where("id = ?", todo.Id).Update()
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Todo Edited Successfully",
	})
	return
}
func DeleteTodo(c *gin.Context) {
	c.Param("todoId")
	todo := &Todo{Id: 6}

	err := dbConnect.Delete(todo)
	if err != nil {
		log.Printf("Error while deleting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully",
	})
	return
}

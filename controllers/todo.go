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
type request1 struct {
	Rid int `json:"Rid"`
}
type CSVResponse struct {
	CSVData string `json:"csv_data,omitempty"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
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
		//c.JSON(http.StatusInternalServerError, gin.H{
		//	"status":  http.StatusInternalServerError,
		//	"message": "Something went wrong",
		//})
		res := &CSVResponse{}
		res.CSVData = err.Error()
		res.Success = false
		c.JSON(http.StatusInternalServerError, res)
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
	todo := new(Todo)
	c.BindJSON(todo)
	insertError := dbConnect.Insert(todo)
	if insertError != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", insertError)
		res := &CSVResponse{}
		res.CSVData = insertError.Error()
		res.Success = false
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Todo created Successfully",
	})
	return
}
func GetSingleTodo(c *gin.Context) {

	c.Param("todoId")
	req := new(request1)
	c.BindJSON(&req)
	rid := req.Rid
	todo := &Todo{Id: rid}
	err := dbConnect.Select(todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		res := &CSVResponse{}
		res.CSVData = err.Error()
		res.Success = false
		c.JSON(http.StatusInternalServerError, res)

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
	var todo Todo

	c.BindJSON(&todo)
	last_name := todo.Last_name

	_, err := dbConnect.Model(&Todo{}).Set("last_name = ?", last_name).Where("id = ?", todo.Id).Update()
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		res := &CSVResponse{}
		res.CSVData = err.Error()
		res.Success = false
		c.JSON(http.StatusInternalServerError, res)
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
	req := new(request1)
	c.BindJSON(&req)
	rid := req.Rid
	todo := &Todo{Id: rid}

	err := dbConnect.Delete(todo)
	if err != nil {
		log.Printf("Error while deleting a single todo, Reason: %v\n", err)
		res := &CSVResponse{}
		res.CSVData = err.Error()
		res.Success = false
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully",
	})
	return
}

package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

var todoList = []todo{
	{ID: 0, Description: "sample description"},
}

func binarySearch(arr *[]todo, target int64) (*todo, error) {

	var low float64 = 0
	var high float64 = float64(len((*arr)))

	for low <= high {
		mid := int64(math.Floor((high + low) / 2))
		if (*arr)[mid].ID == target {
			return &(*arr)[mid], nil
		} else if (*arr)[mid].ID < target {
			low = float64(mid + 1)
		} else if (*arr)[mid].ID > target {
			high = float64(mid - 1)
		}
	}

	return nil, errors.New("Error not able to find target")

}

func getAllTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todoList)
}

func addTodo(c *gin.Context) {
	var newTodo todo

	if err := c.BindJSON(&newTodo); err != nil {
		fmt.Println("ERRROR")
		return
	}

	for _, item := range todoList {
		if item.ID == newTodo.ID {
			c.IndentedJSON(http.StatusAlreadyReported, gin.H{"message": "todo exists"})
			return
		}
	}

	todoList = append(todoList, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodoByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 36, 64)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "invalid id"})
		return
	}

	foundTodo, err := binarySearch(&todoList, id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, (*foundTodo))
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 36, 64)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "invalid id"})
		return
	}

	var deleteInstance *todo = nil
	newArr := []todo{}

	for _, item := range todoList {
		if item.ID == id {
			deleteInstance = &item
			continue
		}
		newArr = append(newArr, item)
	}

	if deleteInstance == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "invalid id"})
		return
	}

	todoList = newArr

	c.IndentedJSON(http.StatusOK, (*deleteInstance))

}

func updateTodo(c *gin.Context) {
	var updateTodo todo

	if err := c.BindJSON(&updateTodo); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	for ind, item := range todoList {
		if item.ID == updateTodo.ID {
			todoList[ind].Description = updateTodo.Description
			c.IndentedJSON(http.StatusOK, todoList[ind])
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
}

func main() {
	router := gin.Default()
	router.GET("/todos", getAllTodos)
	router.GET("/todo/:id", getTodoByID)
	router.POST("/todo", addTodo)
	router.DELETE("/todo/:id", deleteTodo)
	router.PUT("/todo/:id", updateTodo)

	router.Run("localhost:8080")
}

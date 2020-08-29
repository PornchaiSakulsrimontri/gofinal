package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var todos = map[int]*Todo{
	1: &Todo{ID: 1, Title: "pay phone bills", Status: "active"},
}

func getTodoByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	todo := findByID(id)
	if todo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{})
		return
	} else {
		c.JSON(http.StatusOK, todo)
	}
}

func getTodosHandler(c *gin.Context) {
	items := findAll()
	c.JSON(http.StatusOK, items)
}

func createTodosHandler(c *gin.Context) {
	newTodo := Todo{}
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insert(&newTodo)

	c.JSON(http.StatusOK, newTodo)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/todos", getTodosHandler)
	r.GET("/todos/:id", getTodoByIdHandler)
	r.POST("/todos", createTodosHandler)

	return r
}

func main() {
	r := setupRouter()
	r.Run()
}

// ===================================  DB func
func createTable(db *sql.DB) {
	createTb := `
		CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT,
		status TEXT
		);
		`
	_, err := db.Exec(createTb)
	if err != nil {
		log.Println("can't create table", err)
	}
	log.Println("create table success")
}

func findByID(searchId int) Todo {
	db := getConnection()
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1")
	if err != nil {
		log.Println("can't preprare query one row statement", err)
		return Todo{}
	}

	row := stmt.QueryRow(searchId)
	var id int
	var title, status string

	err = row.Scan(&id, &title, &status)
	if err != nil {
		log.Println("can't prepare query todo by id statment", err)
		return Todo{}
	}

	log.Println("one row", id, title, status)
	return Todo{ID: id, Title: title, Status: status}
}

func delete(deleteId int) {
	db := getConnection()
	defer db.Close()

	row := db.QueryRow("DELETE FROM todos WHERE id =$1 ", deleteId)

	err := row.Scan(&newTodo.ID)
	if err != nil {
		log.Println("can't scan id", err)
		return
	}
	log.Println("insert todo success id : ", newTodo.ID)
}

func findAll() []Todo {
	db := getConnection()
	defer db.Close()

	items := []Todo{}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos")
	if err != nil {
		log.Println("can't prepare query all todos statment", err)
		return items
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Println("can't query all todos", err)
		return items
	}

	for rows.Next() {
		var id int
		var title, status string
		err := rows.Scan(&id, &title, &status)
		if err != nil {
			log.Println("can't Scan row into variable", err)
			return items
		}
		item := Todo{ID: id, Title: title, Status: status}
		items = append(items, item)

	}
	log.Println("query all todos success")
	return items
}

func insert(newTodo *Todo) {
	db := getConnection()
	defer db.Close()

	row := db.QueryRow("INSERT INTO todos (title, status) values ($1, $2) RETURNING id", newTodo.Title, newTodo.Status)

	err := row.Scan(&newTodo.ID)
	if err != nil {
		log.Println("can't scan id", err)
		return
	}
	log.Println("insert todo success id : ", newTodo.ID)
}

func getConnection() *sql.DB {
	// log -> uber-go zap
	db, err := sql.Open("postgres", "postgres://cbebltqs:4JW8rDYtGPJJ8SKsy096Yti_6xQ00_Va@rosie.db.elephantsql.com:5432/cbebltqs")
	// os.Getenv("DATABASE_URL")
	// set DATABASE_URL=postgres://cbebltqs:4JW8rDYtGPJJ8SKsy096Yti_6xQ00_Va@rosie.db.elephantsql.com:5432/cbebltqs
	if err != nil {
		log.Println("Connect to database error", err)
	}
	return db
}

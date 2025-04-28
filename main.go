package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var rnd *renderer.Render
var db *gorm.DB

const (
	port = ":9000"
)

type (
	todoModel struct {
		ID        uint      `gorm:"primaryKey" json:"id"`
		Title     string    `gorm:"not null" json:"title"`
		Completed bool      `gorm:"default:false" json:"completed"`
		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	}
	todo struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

func init() {
	rnd = renderer.New()

	dsn := "host=localhost user=postgres password=root123 dbname=todolistdb port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	checkErr(err)

	// Auto migrate the todoModel schema
	err = db.AutoMigrate(&todoModel{})
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Welcome to the Todo API",
	})
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	var todos []todoModel
	if err := db.Find(&todos).Error; err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err.Error(),
		})
		return
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todos,
	})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title is required",
		})
		return
	}

	tm := todoModel{
		Title:     t.Title,
		Completed: false,
	}

	if err := db.Create(&tm).Error; err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to create todo",
			"error":   err.Error(),
		})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo_id": tm.ID,
	})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is required",
		})
		return
	}

	if err := db.Delete(&todoModel{}, id).Error; err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to delete todo",
			"error":   err.Error(),
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo deleted successfully",
	})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is required",
		})
		return
	}

	var t todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title is required",
		})
		return
	}

	var todoItem todoModel
	if err := db.First(&todoItem, id).Error; err != nil {
		rnd.JSON(w, http.StatusNotFound, renderer.M{
			"message": "Todo not found",
		})
		return
	}

	todoItem.Title = t.Title
	todoItem.Completed = t.Completed

	if err := db.Save(&todoItem).Error; err != nil {
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{
			"message": "Failed to update todo",
			"error":   err.Error(),
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
	})
}

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", port, err)
		}
	}()
	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-app/internal/task"
	"todo-app/internal/user"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var db *gorm.DB
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=todo_db sslmode=disable password=199978")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&task.Task{})
	go task.CheckExpiredTasks(db)
	router := mux.NewRouter()
	router = setupRouter(db, router)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRouter(db *gorm.DB, router *mux.Router) *mux.Router {
	router.HandleFunc("/api/users/register/", func(w http.ResponseWriter, r *http.Request) {
		user.Register(db, w, r)
	}).Methods("POST")
	router.HandleFunc("/api/users/login/", func(w http.ResponseWriter, r *http.Request) {
		user.Login(db, w, r)
	}).Methods("POST")
	router.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		user.GetUsers(db, w, r)
	}).Methods("GET")

	router.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		task.GetTasks(db, w, r)
	}).Methods("GET")
	router.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		task.CreateTask(db, w, r)
	}).Methods("POST")
	router.HandleFunc("/api/tasks/{id}/", func(w http.ResponseWriter, r *http.Request) {
		task.UpdateTask(db, w, r)
	}).Methods("PUT")
	router.HandleFunc("/api/tasks/{id}/", func(w http.ResponseWriter, r *http.Request) {
		task.DeleteTask(db, w, r)
	}).Methods("DELETE")
	return router
}

package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Task struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Deadline  time.Time `json:"deadline"`
	Expired   bool      `json:"expired"`
}

func GetTasks(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "List of tasks")
}

func CreateTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Create(&task).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db.Save(&task)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
	fmt.Fprintf(w, "Task created successfully")

}

func UpdateTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task Task
	if err := db.First(&task, params["id"]).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Save(&task).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
	fmt.Fprintf(w, "Task updated successfully")
}

func DeleteTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task Task
	if err := db.First(&task, params["id"]).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := db.Delete(&task).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "Task deleted successfully")
}

func CheckExpiredTasks(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var tasks []Task
			db.Find(&tasks)
			if len(tasks) == 0 {
				return
			}
			for _, task := range tasks {
				if time.Now().After(task.Deadline) {
					task.Expired = true
					db.Save(&task)
				}
			}
			fmt.Println("Expiration checked!")
		}
	}
}

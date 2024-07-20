package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if CheckIfUsernameExists(db, user.Username) {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}
	if err := user.HashPassword(); err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	fmt.Fprintf(w, "User registered successfully")
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}
	if len(u.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func CheckIfUsernameExists(db *gorm.DB, username string) bool {
	var count int
	db.Model(&User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

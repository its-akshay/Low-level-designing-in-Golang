package basics

//============Struct============

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User is our "class" — lowercase fields = private within package
type User struct {
	id        string
	email     string
	password  string // hashed — never exposed raw
	role      string
	createdAt time.Time
}

// Constructor — the Go "new" pattern
func NewUser(email, hashedPassword, role string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	if hashedPassword == "" {
		return nil, errors.New("password hash is required")
	}

	if role == "" {
		role = "user"
	}

	return &User{
		id:        uuid.New().String(),
		email:     email,
		password:  hashedPassword,
		role:      role,
		createdAt: time.Now(),
	}, nil
}

// Controlled getters — expose only what's needed
func (u *User) ID() string {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Role() string {
	return u.role
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// Business logic method — behavior lives on the struct
func (u *User) IsAdmin() bool {
	return u.role == "admin"
}

// No setter for password — use a dedicated method with validation
func (u *User) ChangePassword(newHash string) error {
	if newHash == "" {
		return errors.New("password hash cannot be empty")
	}

	u.password = newHash
	return nil
}

// Additional business logic
func (u *User) PromoteToAdmin() {
	u.role = "admin"
}

// Safe display method — never expose password
func (u *User) Display() {
	fmt.Println("User Details")
	fmt.Println("------------")
	fmt.Println("ID:", u.id)
	fmt.Println("Email:", u.email)
	fmt.Println("Role:", u.role)
	fmt.Println("Created At:", u.createdAt.Format(time.RFC1123))
}

// Example usage
func Example() {
	// Create user
	user, err := NewUser(
		"akshay@example.com",
		"hashed_password_123",
		"user",
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Access through getters
	fmt.Println("User ID:", user.ID())
	fmt.Println("User Email:", user.Email())
	fmt.Println("User Role:", user.Role())

	// Business logic
	fmt.Println("Is Admin?", user.IsAdmin())

	// Change password safely
	err = user.ChangePassword("new_hashed_password")

	if err != nil {
		fmt.Println("Password change failed:", err)
	}

	// Promote user
	user.PromoteToAdmin()

	fmt.Println("Is Admin After Promotion?", user.IsAdmin())

	// Safe display
	user.Display()
}

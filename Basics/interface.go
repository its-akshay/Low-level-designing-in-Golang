// =======================================================
// USER REPOSITORY INTERFACE
// =======================================================

// Define a contract (blueprint)
//
// Any struct that implements ALL these methods
// automatically becomes a UserRepository.
//
// This is the core idea of Dependency Inversion:
// service depends on BEHAVIOR, not implementation.

package basics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

type UserRepository interface {

	// Find a user using ID
	//
	// ctx -> helps with timeout/cancel request
	// id  -> user id to search
	//
	// returns:
	// *User -> found user
	// error -> if something failed
	FindByID(ctx context.Context, id string) (*User, error)

	// Find user using email
	FindByEmail(ctx context.Context, email string) (*User, error)

	// Save user into database/storage
	Save(ctx context.Context, user *User) error

	// Delete user using ID
	Delete(ctx context.Context, id string) error
}

// =======================================================
// POSTGRES IMPLEMENTATION
// =======================================================

// Concrete implementation using PostgreSQL database
//
// This struct "HAS A" database connection.
type PostgresUserRepository struct {
	db *sql.DB
}

// Method implementation for PostgreSQL repository
//
// Because this struct implements all interface methods,
// Go automatically treats it as UserRepository.
func (r *PostgresUserRepository) FindByID(
	ctx context.Context,
	id string,
) (*User, error) {

	// Execute SQL query
	//
	// SELECT id, email, role
	// FROM users
	// WHERE id = $1
	//
	// $1 gets replaced by "id"
	//
	// QueryRowContext returns ONE row.
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, email, role FROM users WHERE id=$1",
		id,
	)

	// Create empty user object
	//
	// We'll fill this using Scan().
	u := &User{}

	// Copy database columns into struct fields
	//
	// Example:
	// DB row:
	// (1, "akshay@gmail.com", "admin")
	//
	// becomes:
	// u.id    = 1
	// u.email = "akshay@gmail.com"
	// u.role  = "admin"
	if err := row.Scan(&u.id, &u.email, &u.role); err != nil {

		// If no row found
		//
		// sql.ErrNoRows is built-in DB error.
		if errors.Is(err, sql.ErrNoRows) {

			// Return custom business error
			return nil, ErrUserNotFound
		}

		// Wrap unknown DB error
		//
		// %w preserves original error.
		return nil, fmt.Errorf(
			"query user by id: %w",
			err,
		)
	}

	// Success -> return filled user
	return u, nil
}

// Other interface methods
func (r *PostgresUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	/* ... */
}

func (r *PostgresUserRepository) Save(
	ctx context.Context,
	user *User,
) error {
	/* ... */
}

func (r *PostgresUserRepository) Delete(
	ctx context.Context,
	id string,
) error {
	/* ... */
}

// =======================================================
// IN-MEMORY IMPLEMENTATION
// =======================================================

// Another implementation of SAME interface
//
// Instead of PostgreSQL,
// data is stored inside memory using map.
//
// Mostly used in:
// - testing
// - mock data
// - fast prototyping
type InMemoryUserRepository struct {

	// RWMutex protects map from concurrent access
	//
	// Multiple readers allowed
	// One writer allowed
	mu sync.RWMutex

	// map[userID] => *User
	users map[string]*User
}

// In-memory implementation of FindByID
func (r *InMemoryUserRepository) FindByID(
	ctx context.Context,
	id string,
) (*User, error) {

	// Read lock
	//
	// Prevents race conditions.
	r.mu.RLock()

	// Automatically unlock at function end
	defer r.mu.RUnlock()

	// Check if user exists in map
	if u, ok := r.users[id]; ok {

		// User found
		return u, nil
	}

	// User not found
	return nil, ErrUserNotFound
}

// =======================================================
// SERVICE LAYER
// =======================================================

// UserService contains business logic
//
// IMPORTANT:
//
// It DOES NOT care whether data comes from:
// - PostgreSQL
// - MySQL
// - Redis
// - Memory
//
// It only cares about the interface.
//
// This is loose coupling.
type UserService struct {

	// Depends on abstraction/interface
	//
	// NOT:
	// repo *PostgresUserRepository
	//
	// GOOD because now implementation can change easily.
	repo UserRepository
}

// Constructor function
//
// Dependency Injection happens here.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// =======================================================
// HOW IT WORKS IN REAL LIFE
// =======================================================

// PRODUCTION:
//
// postgresRepo := &PostgresUserRepository{db: db}
// service := NewUserService(postgresRepo)
//
// Service now uses PostgreSQL internally.

// TESTING:
//
// memoryRepo := &InMemoryUserRepository{
//     users: map[string]*User{},
// }
//
// service := NewUserService(memoryRepo)
//
// Same service code.
// Different repository.
// No code changes needed.

// =======================================================
// MAIN BENEFITS
// =======================================================

// 1. Loose Coupling
// Service is independent of database type.

// 2. Easy Testing
// Use in-memory repo instead of real DB.

// 3. Easy Replacement
// Can switch PostgreSQL -> MongoDB later.

// 4. Clean Architecture
// Business logic separated from storage logic.

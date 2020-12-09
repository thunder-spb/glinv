package pgsql

import (
	"database/sql"
	"strings"

	"glinv/pkg/models"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// UserModel type which wraps a sql.DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// Insert method to add a new record to the users table.
func (m *UserModel) Insert(username, email, password string) error {
	stmt := `INSERT INTO users (user_name, email, hashed_password)
	VALUES($1, $2, $3)`

	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(stmt, username, email, string(hashedPassword))
	if err != nil {
		// TODO
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && strings.Contains(pqErr.Code.Name(), "unique_violation") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

// Auth method to verify whether a user exists with the provided email
// address and password. This will return the relevant user ID if they do.
func (m *UserModel) Auth(email, password string) (int, error) {
	stmt := `SELECT id, hashed_password
	FROM users WHERE email = $1`

	// Retrieve the id and hashed password associated with the given email.
	// If no matching email exists, we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	// Check hashed password and plain-text password provided match.
	// If not, return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

// Get method to fetch details for a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := `SELECT id, user_name, user_role, email, created 
	FROM users WHERE id = $1`

	user := &models.User{}

	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.UserName, &user.UserRole, &user.Email, &user.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

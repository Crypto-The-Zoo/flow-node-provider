package queries

import (
	"InceptionAnimals/app/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserQueries struct{ *sqlx.DB }

// GetUser method for getting one user by given ID.
func (q *UserQueries) GetUser(id uuid.UUID) (models.User, error) {
	// Define user variable
	user := models.User{}

	// Define query string
	query := `SELECT * FROM users WHERE id = $1`

	// Send query to database
	err := q.Get(&user, query, id)
	if err != nil {
		// Return empty object and error
		return user, err
	}

	// Return query result
	return user, nil
}

// CreateUser method for creating one user with email and username
func (q *UserQueries) CreateUser(user *models.User) error {
	// Define query string
	query := `INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	// Send query to database
	_, err := q.Exec(query, user.ID, user.CreatedAt, user.UpdatedAt, user.Username, user.Email, nil, false, nil)
	if err != nil {
		return err
	}

	// return nothing if success
	return nil
}

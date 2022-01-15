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

func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	// Define user variable
	user := models.User{}

	// Define query string
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`

	// Send query to database
	err := q.Get(&user, query, email)
	if err != nil {
		// Return empty object and error
		return user, err
	}

	// Return query result
	return user, nil
}

func (q *UserQueries) GetUserByUsername(username string) (models.User, error) {
	// Define user variable
	user := models.User{}

	// Define query string
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1`

	// Send query to database
	err := q.Get(&user, query, username)
	if err != nil {
		// Return empty object and error
		return user, err
	}

	// Return query result
	return user, nil
}

func (q *UserQueries) GetUserByFlowAddress(flowAddress string) (models.User, error) {
	// Define user variable
	user := models.User{}

	// Define query string
	query := `SELECT * FROM users WHERE flow_address = $1 LIMIT 1`

	// Send query to database
	err := q.Get(&user, query, flowAddress)
	if err != nil {
		// Return empty object and error
		return user, err
	}

	// Return query result
	return user, nil
}

func (q *UserQueries) GetUserPublicByEmail(email string) (models.UserPublic, error) {
	// Define user variable
	userPublic := models.UserPublic{}

	// Define query string
	query := `SELECT id, username, email, flow_address FROM users WHERE email = $1 AND is_active=true LIMIT 1`

	// Send query to database
	err := q.Get(&userPublic, query, email)
	if err != nil {
		// Return empty object and error
		return userPublic, err
	}

	// Return query result
	return userPublic, nil
}

// CreateUser method for creating one user with email and username
func (q *UserQueries) CreateUser(user *models.User) error {
	// Define query string
	query := `INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	// Send query to database
	_, err := q.Exec(query, user.ID, user.CreatedAt, user.UpdatedAt, user.Username, user.Email, user.FlowAddress, false, user.LoginObj)
	if err != nil {
		return err
	}

	// return nothing if success
	return nil
}

// CreateUser method for creating one user with email and username
func (q *UserQueries) CreateLoginCode(user *models.User) error {
	// Define query string
	query := `UPDATE users SET login_obj = $1 WHERE id = $2`

	// Send query to database
	_, err := q.Exec(query, user.LoginObj, user.ID)
	if err != nil {
		return err
	}

	// return nothing if success
	return nil
}

func (q *UserQueries) GetLoginCode(email string) (models.LoginObj, error) {
	loginObj := models.LoginObj{}

	query := `SELECT login_obj FROM users where email = $1 LIMIT 1`

	err := q.Get(&loginObj, query, email)
	if err != nil {
		return loginObj, err
	}

	return loginObj, err
}

func (q *UserQueries) DeleteLoginCode(email string) error {
	loginObj := models.LoginObj{}
	query := `UPDATE users SET login_obj=$1 where email = $2`

	_, err := q.Exec(query, loginObj, email)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) ActivateUser(email string) error {
	query := `UPDATE users SET is_active=true where email = $1`

	_, err := q.Exec(query, email)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) DeleteInactiveUser(user *models.User) error {
	query := `DELETE from users where email = $1 OR username= $2 AND is_active=false`

	_, err := q.Exec(query, user.Email, user.Username)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQueries) AddFlowAddressToUser(email string, flowAddress string) error {
	query := `UPDATE users SET flow_address = $1 WHERE email = $2`

	_, err := q.Exec(query, flowAddress, email)
	if err != nil {
		return err
	}

	return nil
}

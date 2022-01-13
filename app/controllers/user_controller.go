package controllers

import (
	"InceptionAnimals/app/models"
	"InceptionAnimals/platform/database"
	"time"

	"InceptionAnimals/pkg/utils"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
)

func GetUser(ctx *fiber.Ctx) {
	// Catch user id from URL
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get user by id
	// TODO: get user by jwt token
	user, err := db.GetUser(id)
	if err != nil {
		// Return, if user not found
		ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "user with the given ID is not found",
			"user":  nil,
		})
	}

	// Return status 200 OK
	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
	})

	// user := ctx.Locals("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// id := claims["sub"].(string)
	// ctx.Send(fmt.Sprintf("Sanity Check: id %s", id))
}

func CreateUser(ctx *fiber.Ctx) {

	// Create new User struct
	user := &models.User{}

	// Check, if received JSON data is valid
	if err := ctx.BodyParser(user); err != nil {
		// Return status 400 and error message
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// Create database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// Create a new validator for a User model
	validate := utils.NewValidator()

	// Set initialized default data for user
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Validate user fields
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
		return
	}

	// Create user in database
	if err := db.CreateUser(user); err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// Send status 201 with user object
	ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
	})
}

// func Login(ctx *fiber.Ctx) {
// 	type request struct {
// 		Email string `json:"email"`
// 		Code  string `json:"code"`
// 	}

// 	var body request
// 	err := ctx.BodyParser(&body)
// 	if err != nil {
// 		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "failed_to_parse_json",
// 		})
// 		return
// 	}

// 	if body.Email != "alliu930410@gmail.com" || body.Code != "CodeFromDb" {
// 		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Bad Credentials",
// 		})
// 		return
// 	}

// 	token := jwt.New(jwt.SigningMethodHS256)
// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["sub"] = "1"
// 	claims["exp"] = time.Now().Add(time.Hour * 24 * 7) // valid for a week

// 	s, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
// 	if err != nil {
// 		ctx.SendStatus(fiber.StatusInternalServerError)
// 		return
// 	}

// 	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"token": s,
// 		"user": struct {
// 			Id    int    `json:"id"`
// 			Email string `json:"email"`
// 		}{
// 			Id:    1,
// 			Email: "alliu930410@gmail.com",
// 		},
// 	})
// }

func GetLoginCode(ctx *fiber.Ctx) {
	now := time.Now()

	// Define request struct
	type request struct {
		Email string `json:"email"`
	}

	// Parse body from JSON request
	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// TODO: validate email field
	// if body.Email == nil {
	// 	ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   "invalid_email",
	// 	})
	// 	return
	// }

	// Create database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// Create User Login Object
	// duration := 10
	loginObj := models.LoginObj{
		Code:      "12345678",
		CreatedAt: now,
		// ExpiresAt: now.Add(time.Minute * duration),
		ExpiresAt: now,
	}

	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// CreateLoginCode in database
	user.LoginObj = loginObj
	if err := db.CreateLoginCode(&user); err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// Return LoginCode
	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"loginObj": loginObj,
	})
}

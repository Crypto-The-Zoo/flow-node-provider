package controllers

import (
	"InceptionAnimals/app/models"
	"InceptionAnimals/platform/database"
	"fmt"
	"math/rand"
	"time"

	"InceptionAnimals/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUser(ctx *fiber.Ctx) error {
	// Catch user id from URL
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get user by id
	user, err := db.GetUser(id)
	if err != nil {
		// Return, if user not found
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "user with the given ID is not found",
			"user":  nil,
		})
	}

	// Return status 200 OK
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
	})
}

func Register(ctx *fiber.Ctx) error {

	// Create new User struct
	user := &models.User{}

	// Check, if received JSON data is valid
	if err := ctx.BodyParser(user); err != nil {
		// Return status 400 and error message
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Check if email is taken in database
	existingEmailUser, err := db.GetUserByEmail(user.Email)
	if err == nil {
		if existingEmailUser.IsActive {
			// Return, if some fields are not valid.
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "email_taken",
			})
		}
	}

	// Check if username is taken in database
	existingUsernameUser, err := db.GetUserByUsername(user.Username)
	if err == nil {
		if existingUsernameUser.IsActive {
			// Return, if some fields are not valid.
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "username_taken",
			})
		}
	}

	// Delete inactive user records if there is any
	if existingEmailUser.LoginObj.Code != "" {
		if err := db.DeleteInactiveUser(user); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}
	if existingUsernameUser.LoginObj.Code != "" {
		if err := db.DeleteInactiveUser(user); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}

	// Create User Login Object, valid for 30 minutes
	now := time.Now()
	duration := 30
	code := fmt.Sprintf("%d%d%d%d%d%d%d%d", rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9))
	user.LoginObj = models.LoginObj{
		Code:      code,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Minute * time.Duration(duration)),
	}

	// Create user in database
	if err := db.CreateUser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Send status 201 with user object
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
	})
}

func Login(ctx *fiber.Ctx) error {
	type request struct {
		Email string `json:"email" validate:"required,min=3,max=32"`
		Code  string `json:"code" validate:"required,min=3,max=32"`
	}

	var body request
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed_to_parse_json",
		})
	}

	// Create a new validator for a request model
	validate := utils.NewValidator()
	// Validate request fields
	if err := validate.Struct(body); err != nil {
		// Return, if some fields are not valid.
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get user object loginObj
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Validate login code
	if user.LoginObj.Code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "code_not_initialized",
		})
	}
	if body.Code != user.LoginObj.Code {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid_code",
		})
	}
	if now := time.Now().Unix(); now > user.LoginObj.ExpiresAt.Unix() {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "code_expired",
		})
	}

	// Invalidate login code
	if err := db.DeleteLoginCode(body.Email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	user.LoginObj = models.LoginObj{}

	// Delete user login code in database
	if err := db.DeleteLoginCode(body.Email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Activate user if active field is false
	if !user.IsActive {
		if err := db.ActivateUser(body.Email); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}

	if err := db.DeleteLoginCode(body.Email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Generate new access token for user
	userPublic, err := db.GetUserPublicByEmail(body.Email)
	if err != nil {
		// Return status 500 and token generation error.
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	token, err := utils.GenerateNewAccessToken(&userPublic)
	if err != nil {
		// Return status 500 and token generation error.
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return access token
	return ctx.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"token": token,
	})
}

func AddFlowAddress(ctx *fiber.Ctx) error {
	now := time.Now().Unix()

	type request struct {
		FlowAddress string `json:"flowAddress" validate:"required,len=18"`
	}

	// Parse body from JSON request
	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a request model
	validate := utils.NewValidator()
	// Validate request fields
	if err := validate.Struct(body); err != nil {
		// Return, if some fields are not valid.
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Get claims from JWT
	claims, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		// Return status 500 and JWT parse error.
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Check if user has already has wallet registered
	emailUser, err := db.GetUserByEmail(claims.Email)
	if err == nil {
		if emailUser.FlowAddress != "" && emailUser.FlowAddress == body.FlowAddress {
			// silent return 200
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"error": false,
				"msg":   "wallet_already_registered",
			})
		}
		if emailUser.FlowAddress != "" && emailUser.FlowAddress != body.FlowAddress {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "different_wallet_registered",
			})
		}
	}

	// Check if wallet has already been registered by user
	walletUser, err := db.GetUserByFlowAddress(body.FlowAddress)
	if err == nil {
		if walletUser.Email == claims.Email {
			// silent return 200
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"error": false,
				"msg":   "wallet_already_registered",
			})
		}
		if walletUser.Email != claims.Email {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "wallet_already_registered_by_another_user",
			})
		}
	}

	// Register wallet with user in database
	if err := db.AddFlowAddressToUser(claims.Email, body.FlowAddress); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
	})
}

func GetLoginCode(ctx *fiber.Ctx) error {
	now := time.Now()

	// Define request struct
	type request struct {
		Email string `json:"email" validate:"required,email,min=6,max=32"`
	}

	// Parse body from JSON request
	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a request model
	validate := utils.NewValidator()
	// Validate request fields
	if err := validate.Struct(body); err != nil {
		// Return, if some fields are not valid.
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create database connection
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// check last loginObj
	lastLoginObj, err := db.GetLoginCode(body.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	fmt.Println(lastLoginObj)
	if lastLoginObj.Code != "" && lastLoginObj.CreatedAt.Add(time.Minute).Unix() > now.Unix() {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "retry_too_frequent",
		})
	}

	// Create User Login Object, valid for 30 minutes
	duration := 30
	code := fmt.Sprintf("%d%d%d%d%d%d%d%d", rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9))
	loginObj := models.LoginObj{
		Code:      code,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Minute * time.Duration(duration)),
	}

	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// CreateLoginCode in database
	user.LoginObj = loginObj
	if err := db.CreateLoginCode(&user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return LoginCode
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"loginObj": loginObj,
	})
}

package auth

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidUsername    = errors.New("username must be 3-50 characters and alphanumeric")
	ErrInvalidEmail       = errors.New("invalid email format")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,50}$`)

// Register creates a new user account
func Register(username, email, password string) (*User, string, error) {
	// Validate input
	if err := validateUsername(username); err != nil {
		return nil, "", err
	}
	if err := validateEmail(email); err != nil {
		return nil, "", err
	}
	if err := validatePassword(password); err != nil {
		return nil, "", err
	}

	// Check if user already exists
	if existingUser, _ := GetUserByUsername(username); existingUser != nil {
		return nil, "", errors.New("username already taken")
	}
	if existingUser, _ := GetUserByEmail(email); existingUser != nil {
		return nil, "", errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, "", errors.New("failed to hash password")
	}

	// Create user
	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := CreateUser(user); err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// Login authenticates a user and returns a token
func Login(username, password string) (*User, string, error) {
	// Get user from database
	user, err := GetUserByUsername(username)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	// Verify password
	if err := verifyPassword(user.PasswordHash, password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// ValidateUser checks if a user exists and returns their information
func ValidateUser(userID string) (*User, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ValidateToken validates a token and returns the associated user
func ValidateToken(tokenString string) (*User, error) {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Helper functions

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

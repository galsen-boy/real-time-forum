package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/structure"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles the registration endpoint
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Prevents the endpoint being called by other URL paths
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Prevent all request types other than POST
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Stores the unmarshalled register data
	var newUser structure.User

	// Decodes the request body into the user struct
	// Returns a bad request if there's an error
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "400 bad request.", http.StatusBadRequest)
		return
	}
	// Validate the email format
	if !isValidEmail(newUser.Email) {
		http.Error(w, "400 bad request: Invalid email address.", http.StatusBadRequest)
		return
	}
	//checks if age is on valid format
	age, err := strconv.Atoi(newUser.DOB)
	if err != nil || age < 0 {
		http.Error(w, "400 bad request: Invalid date of birth.", http.StatusBadRequest)
		return
	}

	// Check if the email or username already exists
	emailExists, err := database.UserExists(config.Path, newUser.Email)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	usernameExists, err := database.UserExists(config.Path, newUser.Username)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	if emailExists && usernameExists {
		http.Error(w, "409 conflict: Email and username already exist.", http.StatusConflict)
		return
	} else if emailExists {
		http.Error(w, "409 conflict: The email you entered is already taken.", http.StatusConflict)
		return
	} else if usernameExists {
		http.Error(w, "409 conflict: The username you entered is already taken.", http.StatusConflict)
		return
	}

	// Generate the password hash for the user
	passwordHash, err := GenerateHash(newUser.Password)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	newUser.Password = passwordHash

	// Attempts to add the new user to the database
	err = database.NewUser(config.Path, newUser)
	if err != nil {
		http.Error(w, "500 internal server error: Failed to register user.", http.StatusInternalServerError)
		return
	}

	// Sends a message back if successfully registered
	var msg = structure.Resp{Msg: "Successful registration"}

	resp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "500 internal server error: Failed to marshal response. ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Generates a hash from a given password
func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)

	return string(hash), err
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	// Email validation pattern using regular expression
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/structure"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	// Prevents the endpoint being called by other URL paths
	if r.URL.Path != "/session" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Prevents all request types other than POST
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		// Handle the absence of session cookie here
		// For example, you can log the error or perform other actions
		fmt.Println("Session cookie not found or expired")
		// Set a dummy value for the cookie to avoid "401 unauthorized" error
		cookie = &http.Cookie{Name: "session", Value: "dummy"}
	}

	foundVal := cookie.Value

	curr, err := database.CurrentUser(config.Path, foundVal)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	cid := strconv.Itoa(curr.Id)

	// Create the response struct
	resp := structure.Resp{
		Msg: cid + "|" + curr.Username,
	}

	// Encode the response as JSON and write it to the response writer
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
}

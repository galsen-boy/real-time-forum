package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/structure"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//Prevents the endpoint being called by other url paths
	if r.URL.Path != "/logout" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	//Prevents all request types other than POST
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	//Opens the database
	db, err := database.OpenDB(config.Path)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	//Checks for session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	//Removes session from the database
	_, err = db.Exec(database.RemoveCookie, cookie.Value)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	//Removes cookie from the browser
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	//Sends a message back if successfully logged out
	var msg = structure.Resp{Msg: "Goodbye"}

	resp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

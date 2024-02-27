package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	//Prevents the endpoint being called by other url paths
	if r.URL.Path != "/user" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	//Prevents all request types other than GET
	if r.Method != "GET" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	//Check whether an id is passed in the url (/user?id=)
	//If no, get all users. If yes, get user with matching id
	id := r.URL.Query().Get("id")
	if id == "" {
		users, err := database.FindAllUsers(config.Path)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Marshals the array of user structs to a json object
		resp, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Writes the json object to the frontend
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else {
		user, err := database.FindUserByParam(config.Path, "id", id)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Marshals the user struct to a json object
		resp, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Writes the json object to the frontend
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

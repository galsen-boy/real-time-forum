package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/structure"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/chat" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	user_id := r.URL.Query().Get("user_id")
	if user_id == "" {
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	uid, err := strconv.Atoi(user_id)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	users, err := database.FindUserChats(config.Path, uid)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	var chatUsers []int

	for _, u := range users {
		if u.User_one == uid {
			chatUsers = append(chatUsers, u.User_two)
		} else {
			chatUsers = append(chatUsers, u.User_one)
		}
	}

	var msg = structure.OnlineUsers{
		UserIds:  chatUsers,
		Msg_type: "",
	}

	resp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

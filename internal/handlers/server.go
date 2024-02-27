package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"real-time-forum/internal/chat"
	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
)

// Sets up the router with endpoints and starts the server
func StartServer() {
	database.InitDB(config.Path)

	mux := http.NewServeMux()
	hub := chat.NewHub()
	go hub.Run()

	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("./frontend"))))

	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/session", SessionHandler)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/logout", LogoutHandler)
	mux.HandleFunc("/register", RegisterHandler)
	mux.HandleFunc("/user", UserHandler)
	mux.HandleFunc("/post", PostHandler)
	mux.HandleFunc("/message", MessageHandler)
	mux.HandleFunc("/comment", CommentHandler)
	mux.HandleFunc("/like", LikeHandler)
	mux.HandleFunc("/chat", ChatHandler)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})

	fmt.Println("Server running on port 8000....")
	openBrowser("http://localhost:8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}

// Opens the browser to the specified url
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}

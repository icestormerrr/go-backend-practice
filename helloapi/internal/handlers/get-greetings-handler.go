package handlers

import "net/http"

func GetGreetingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/entrypoint", habrbestbot.entrypoint)
	http.ListenAndServe(":8080", nil)
}

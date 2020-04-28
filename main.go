package main

import (
	"net/http"

	"github.com/baor/habr-best-bot/habrbestbot"
)

func main() {
	http.HandleFunc("/entrypoint", habrbestbot.Entrypoint)
	http.ListenAndServe(":8080", nil)
}

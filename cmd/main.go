package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var startPage = template.Must(template.ParseFiles("templates/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := startPage.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	port := ""

	err := godotenv.Load()
	if err != nil {
		log.Println(color.RedString("Error: Can't loading .env file"))
	} else {
		port = os.Getenv("PORT")
	}

	if port == "" {
		port = "4999"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)

	fmt.Println(color.GreenString("Starting server at:"))
	fmt.Printf(color.GreenString("http://0.0.0.0:%s", port))

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("could not start server:\n\t %v", err)
	}
}

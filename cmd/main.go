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
	fmt.Println(color.CyanString("Start server"))

	port := ""

	err := godotenv.Load()
	if err != nil {
		log.Println(color.RedString("Error loading .env file"))
	} else {
		port = os.Getenv("PORT")
	}

	if port == "" {
		port = "4999"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}

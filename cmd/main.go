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

type GeneralInfo struct {
	SNP      string `json:"SNP"`
	Birthday string `json:"birthday"`
	About    string `json:"about"`
	Position string `json:"position"`
}

type Educations []struct {
	Type         string `json:"type"`
	Organisation string `json:"organisation"`
	Grad         string `json:"grad"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Department   string `json:"department"`
}

type Contacts struct {
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	GithubLink string `json:"github_link"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Relocation string `json:"relocation"`
}

type Languages []struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type WorkExpiriences struct {
	Seniority  string `json:"seniority"`
	Workplaces []struct {
		Organisation string `json:"organisation"`
		Position     string `json:"position"`
		Start        string `json:"start"`
		End          string `json:"end"`
		Projects     []struct {
			Description string `json:"description"`
			Stack       string `json:"stack"`
			Comment     string `json:"comment"`
		} `json:"projects"`
	} `json:"workplaces"`
}

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

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	fmt.Println(color.GreenString("Starting server at:"))
	fmt.Printf(color.GreenString("http://0.0.0.0:%s", port))

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("could not start server:\n\t %v", err)
	}
}

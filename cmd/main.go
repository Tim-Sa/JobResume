package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type JsonLoad interface{}

func jsonLoad(filePath string, data JsonLoad) {
	file, _ := ioutil.ReadFile(filePath)
	_ = json.Unmarshal([]byte(file), data)
}

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

type Projects []struct {
	Description string `json:"description"`
	Stack       string `json:"stack"`
	Comment     string `json:"comment"`
}

type Workplaces []struct {
	Organisation string   `json:"organisation"`
	Position     string   `json:"position"`
	Start        string   `json:"start"`
	End          string   `json:"end"`
	Projects     Projects `json:"projects"`
}

type WorkExpiriences struct {
	Seniority  string     `json:"seniority"`
	Workplaces Workplaces `json:"workplaces"`
}

type Skills []struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type PageContent struct {
	GeneralInfo     GeneralInfo
	Educations      Educations
	Contacts        Contacts
	Languages       Languages
	WorkExpiriences WorkExpiriences
	Skills          Skills
}

func getContent() PageContent {
	generalInfo := &GeneralInfo{}
	jsonLoad("assets/content/general_info.json", generalInfo)

	educations := &Educations{}
	jsonLoad("assets/content/educations.json", educations)

	contacts := &Contacts{}
	jsonLoad("assets/content/contacts.json", contacts)

	languages := &Languages{}
	jsonLoad("assets/content/languages.json", languages)

	workExpiriences := &WorkExpiriences{}
	jsonLoad("assets/content/work_exp.json", workExpiriences)

	skills := &Skills{}
	jsonLoad("assets/content/skills.json", skills)

	pageContent := PageContent{
		GeneralInfo:     *generalInfo,
		Educations:      *educations,
		Contacts:        *contacts,
		Languages:       *languages,
		WorkExpiriences: *workExpiriences,
		Skills:          *skills,
	}

	return pageContent
}

var startPage = template.Must(template.ParseFiles("templates/index.html"))
var pageContent = getContent()

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := startPage.Execute(w, pageContent)
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

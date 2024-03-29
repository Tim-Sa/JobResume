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
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/joho/godotenv"
)

func getConfig(path string) map[string]any {
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles(path)
	if err != nil {
		panic(err)
	}
	configData := config.Data()
	log.Printf(color.CyanString("Configuration: \n %#v\n"), configData)

	return configData
}

type JsonLoad interface{}

func jsonLoad(filePath string, data JsonLoad) {
	file, _ := ioutil.ReadFile(filePath)
	_ = json.Unmarshal([]byte(file), data)
}

type GeneralInfo struct {
	SNP      string `json:"SNP"`
	Position string `json:"position"`
}

type Summary []struct {
	Paragraph string `json:"paragraph"`
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
	Telegram   string `json:"telegram"`
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
	Summary         Summary
	Educations      Educations
	Contacts        Contacts
	Languages       Languages
	WorkExpiriences WorkExpiriences
	Skills          Skills
	PortraitPath    string
}

func getContent() PageContent {
	contentConfig := getConfig("content_config.yaml")
	jsonPath := fmt.Sprintf("%v", contentConfig["content_path"])
	portraitPath := fmt.Sprintf("%v", contentConfig["portrait_path"])

	generalInfo := &GeneralInfo{}
	jsonLoad(jsonPath+"general_info.json", generalInfo)

	educations := &Educations{}
	jsonLoad(jsonPath+"educations.json", educations)

	contacts := &Contacts{}
	jsonLoad(jsonPath+"contacts.json", contacts)

	languages := &Languages{}
	jsonLoad(jsonPath+"languages.json", languages)

	workExpiriences := &WorkExpiriences{}
	jsonLoad(jsonPath+"work_exp.json", workExpiriences)

	skills := &Skills{}
	jsonLoad(jsonPath+"skills.json", skills)

	summary := &Summary{}
	jsonLoad(jsonPath+"summary.json", summary)

	pageContent := PageContent{
		GeneralInfo:     *generalInfo,
		Educations:      *educations,
		Contacts:        *contacts,
		Languages:       *languages,
		WorkExpiriences: *workExpiriences,
		Skills:          *skills,
		Summary:         *summary,
		PortraitPath:    portraitPath,
	}

	return pageContent
}

var startPage = template.Must(template.ParseFiles("templates/index.html"))
var pdfPage = template.Must(template.ParseFiles("templates/to_pdf.html"))
var pageContent = getContent()

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("\nStart Page GET request from: %s\n", r.RemoteAddr)

	err := startPage.Execute(w, pageContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("\nPdf Page GET request from: %s\n", r.RemoteAddr)

	err := pdfPage.Execute(w, pageContent)
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
	mux.HandleFunc("/pdf", pdfHandler)

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	fmt.Println(color.GreenString("Starting server at:"))
	fmt.Printf(color.GreenString("http://0.0.0.0:%s", port))

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("could not start server:\n\t %v", err)
	}
}

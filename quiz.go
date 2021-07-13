package main

import (
	"encoding/json"
	f "fmt"
	"html/template"
	"log"
	h "net/http"
	"os"

	"github.com/gorilla/mux"
)

type RadioButton struct {
	Id    string
	Name  string
	Value string
	Text  string
}

type Input struct {
	Question     string
	RadioButtons []RadioButton
}
type PageInput struct {
	PageTitle string
	Inputs    []Input
}

type Answer struct {
	AnswerGiven string `json:"Answer:"`
}
type Answers []Answer

func main() {
	rtr := mux.NewRouter().StrictSlash(true)
	rtr.HandleFunc("/", homePage).Methods("GET")
	rtr.HandleFunc("/questions", quizPage).Methods("GET")
	rtr.HandleFunc("/answers", answerPage).Methods("POST")
	rtr.HandleFunc("check", checkAnswers).Methods("GET")
	log.Fatal(h.ListenAndServe(":8080", rtr))
}

func homePage(w h.ResponseWriter, r *h.Request) {
	f.Fprintf(w, "<h1>Welcome to the homepage!</h1>")
	f.Println("Endpoint hit: homePage")
}

func quizPage(w h.ResponseWriter, r *h.Request) {
	f.Println("Endpoint hit: quizPage")
	t, err := template.ParseFiles("html_quiz.html")
	if err != nil {
		f.Println(err)
	}
	RadioButtonsQ1 := []RadioButton{
		RadioButton{Id: "Two", Name: "Answerq1", Value: "2", Text: "2"},
		RadioButton{Id: "Three", Name: "Answerq1", Value: "3", Text: "3"},
		RadioButton{Id: "Four", Name: "Answerq1", Value: "4", Text: "4"},
	}

	DataQ1 := Input{
		Question:     "What is 1+1?",
		RadioButtons: RadioButtonsQ1,
	}

	RadioButtonsQ2 := []RadioButton{
		RadioButton{Id: "Twelve", Name: "Answerq2", Value: "12", Text: "12"},
		RadioButton{Id: "Fifteen", Name: "Answerq2", Value: "15", Text: "15"},
		RadioButton{Id: "Fourteen", Name: "Answerq2", Value: "14", Text: "14"},
	}

	DataQ2 := Input{
		Question:     "What is 5+10?",
		RadioButtons: RadioButtonsQ2,
	}

	PageData := PageInput{
		PageTitle: "Quiz questions:",
		Inputs: []Input{
			DataQ1,
			DataQ2,
		},
	}

	err = t.Execute(w, PageData)
	if err != nil {
		f.Println(err)
	}
}

func answerPage(w h.ResponseWriter, r *h.Request) {
	f.Println("Endpoint hit: answerPage")

	file, err := os.OpenFile("answers.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		f.Println(err)
	}
	defer file.Close()

	r.ParseForm()
	MyAnswers := Answers{
		Answer{r.FormValue("Answerq1")},
		Answer{r.FormValue("Answerq2")},
	}

	jfile, err := json.Marshal(MyAnswers)
	if err != nil {
		f.Println(err)
	}

	file.Write(jfile)
	file.Close()

	//This will output JSON with answers on the page .../answers
	//r.ParseForm()
	//MyAnswers := Answers{
	//Answer{r.FormValue("Answerq1")},
	//Answer{r.FormValue("Answerq2")},
	//}
	//json.NewEncoder(w).Encode(MyAnswers)

	//You can simply output answers in the terminal:
	//fmt.Println("Your answer to question one is", r.FormValue("Answerq1"))
	//fmt.Println("Your answer to questions two is", r.FormValue("Answerq2"))
}

func checkAnswers(w h.ResponseWriter, r *h.Request) {

}

//Here we create a REST API
//We output a message to the homepage
//We create a database with a table that has all our questions then we populate the database
//We display questions to the .../pages page that we get from our database but we also use html template
//We then save the answers to our database
//We then compare the answers to the correct ones

package main

import (
	"database/sql"
	f "fmt"
	"html/template"
	"log"
	h "net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "postgres"
	NAME     = "myquiz"
)

type RadioButton struct {
	Id         string
	Name       string
	Value      string
	Text       string
	QuestionID string
}

type Input struct {
	Question     string
	RadioButtons []RadioButton
}

type PageInputs struct {
	Inputs []Input
}

var count int

var numQuestions = make([]Input, count)

var PageData = PageInputs{
	Inputs: numQuestions,
}

var Radio RadioButton
var db *sql.DB
var err error

var wrongQuestions string
var given_answer string
var correct_answer string
var Outputs []string

func Database() {

	//dbinfo is all the information we need to connect to our postgres database
	dbinfo := f.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASSWORD, NAME)
	dbPoint := &db
	//sql.Open opens a connection to ourdatabase but does not create a connection, it simply validates the arguments provided
	*dbPoint, err = sql.Open("postgres", dbinfo)
	if err != nil {
		f.Println("Open err: ", err)
	}
	//db.Ping actually creates a connection to our database
	err = db.Ping()
	if err != nil {
		f.Println("Ping err: ", err)
	}

	//create a table radiobuttons, answers, and questions
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS radiobuttons (id varchar(1024), name varchar(1024), value varchar(1024), text varchar(1024), question_id varchar(1024));
	CREATE TABLE IF NOT EXISTS questions (id TEXT, content TEXT);
	CREATE TABLE IF NOT EXISTS answers (id varchar(1024), question_id varchar(1024), correct_answer varchar(1024), given_answer varchar(1024))`)
	if err != nil {
		f.Println("Create table error:", err)
	}
	//insert records into radiobuttons
	_, err = db.Query(`INSERT INTO radiobuttons (id, name, value, text, question_id) VALUES ('Two', 'Answerq1', '2', '2', '1'),
	('Three', 'Answerq1', '3', '3', '1'), ('Four', 'Answerq1', '4', '4', '1'), ('Twelve', 'Answerq2', '12', '12', '2'), 
	('Fifteen', 'Answerq2', '15', '15', '2'), ('Fourteen', 'Answerq2', '14', '14', '2')`)
	if err != nil {
		f.Println("Insert into table:", err)
	}

	_, err = db.Query("INSERT INTO questions (id, content) VALUES ('1', 'What is 1+1?'), ('2', 'What is 10 + 5?')")
	if err != nil {
		f.Println("Insert into table:", err)
	}

	_, err = db.Query("INSERT INTO answers (id, question_id, correct_answer, given_answer) VALUES ('1', '1', '2', NULL),('2', '2', '15', NULL)")
	if err != nil {
		f.Println("Insert into table:", err)
	}

	//here we will learn how many questions we have in out database

	err := db.QueryRow("SELECT COUNT(*) FROM answers").Scan(&count)
	if err != nil {
		f.Println(err)
	}

	pointNumQuestions := &numQuestions
	*pointNumQuestions = make([]Input, count)

	pointPageData := &PageData
	*pointPageData = PageInputs{
		Inputs: numQuestions,
	}

	//populate struct from the table radiobuttons to use later in html
	rows, err := db.Query("SELECT * FROM radiobuttons WHERE question_id='1'")
	if err != nil {
		f.Println("Select from radiobuttons:", err)
	}

	for rows.Next() {
		err = rows.Scan(&Radio.Id, &Radio.Name, &Radio.Value, &Radio.Text, &Radio.QuestionID)
		if err != nil {
			f.Println("Scan radiobuttons:", err)
		}
		(PageData.Inputs[0]).RadioButtons = append((PageData.Inputs[0]).RadioButtons, Radio)

	}

	row := db.QueryRow("SELECT content FROM questions WHERE id='1' ")
	err = row.Scan(&PageData.Inputs[0].Question)
	if err != nil {
		f.Println("Select question:", err)
	}

	rows, err = db.Query("SELECT * FROM radiobuttons WHERE question_id='2'")
	if err != nil {
		f.Println("Select from radiobuttons:", err)
	}

	for rows.Next() {
		err = rows.Scan(&Radio.Id, &Radio.Name, &Radio.Value, &Radio.Text, &Radio.QuestionID)
		if err != nil {
			f.Println("Scan radiobuttons:", err)
		}
		(PageData.Inputs[1]).RadioButtons = append((PageData.Inputs[1]).RadioButtons, Radio)

	}

	row = db.QueryRow("SELECT content FROM questions WHERE id='2' ")
	err = row.Scan(&PageData.Inputs[1].Question)
	if err != nil {
		f.Println("Select question:", err)
	}

}

func main() {

	Database()

	rtr := mux.NewRouter().StrictSlash(true)
	rtr.HandleFunc("/", homePage).Methods("GET")
	rtr.HandleFunc("/questions", quizPage).Methods("GET")
	rtr.HandleFunc("/answers", answerPage).Methods("POST")
	rtr.HandleFunc("/check", checkAnswers).Methods("GET")
	log.Fatal(h.ListenAndServe(":8080", rtr))
}

func homePage(w h.ResponseWriter, r *h.Request) {
	//note here we output "Welcome to the homepage" as a html heading
	f.Fprintf(w, "<h1>Welcome to the homepage!</h1>")
	f.Println("Endpoint hit: homePage")
}

func quizPage(w h.ResponseWriter, r *h.Request) {

	f.Println("Endpoint hit: quizPage")
	t, err := template.ParseFiles("html_quiz_notitle.html")
	if err != nil {
		f.Println("Parsing error:", err)
	}

	err = t.Execute(w, PageData)
	if err != nil {
		f.Println("Execute error:", err)
	}

}

func answerPage(w h.ResponseWriter, r *h.Request) {
	f.Println("Endpoint hit: answerPage")

	r.ParseForm()

	AnswergivenQ1 := r.FormValue("Answerq1")
	AnswerGivenQ2 := r.FormValue("Answerq2")

	_, err := db.Query("UPDATE answers SET given_answer = $1 WHERE question_id='1' ", AnswergivenQ1)
	if err != nil {
		f.Println(err)
	}

	_, err = db.Query("UPDATE answers SET given_answer = $1 WHERE question_id='2' ", AnswerGivenQ2)
	if err != nil {
		f.Println(err)
	}

}

func checkAnswers(w h.ResponseWriter, r *h.Request) {
	//row := db.QueryRow("SELECT question_id FROM answers WHERE correct_answer != given_answer AND id='1' ")
	//if err != nil {
	//f.Println(err)
	//}
	//err = row.Scan(&Output1)
	//if err != nil {
	//row = db.QueryRow("SELECT question_id FROM answers WHERE correct_answer != given_answer AND id='2' ")
	//row.Scan(&Output1)
	//}
	//f.Fprintf(w, "The question you got wrong is "+Output1)

	rows, err := db.Query("SELECT question_id, given_answer, correct_answer FROM answers WHERE correct_answer != given_answer")
	if err != nil {
		f.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&wrongQuestions, &given_answer, &correct_answer)
		if err != nil {
			f.Println(err)
		}
		f.Fprintf(w, "You got this question wrong: "+wrongQuestions+", the answer you gave is "+given_answer+", the correct answer is "+correct_answer+"\n")
	}
	//err = rows.Scan(&Outputs[i-1])
	//if err != nil {
	//f.Println(err)
	//}

	//f.Println(Outputs)

	//if err != nil {
	//row = db.QueryRow("SELECT question_id FROM answers WHERE correct_answer != given_answer AND id='2' ")
	//row.Scan(&Output1)
	//}
	//f.Fprintf(w, "The question you got wrong is "+Output1)

}

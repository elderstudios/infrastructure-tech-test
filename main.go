package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type server struct{}

var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
os.Getenv("DB_HOST"),os.Getenv("PSQL_USER"), os.Getenv("PSQL_PASS"), os.Getenv("PSQL_DBNAME"), os.Getenv("PSQL_PORT"))



func runMigrations()  {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("FailedToConnectToDatabase")
	}

	db.AutoMigrate(&Member{})
	log.Infof("SchemaUpToDate")
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		get(w)
	case "POST":
		post(w,r)
	default:
		w.WriteHeader(404)
	}
}

func post(w http.ResponseWriter, r *http.Request)  {

	decoder := json.NewDecoder(r.Body)
	var person MemberDto
	err := decoder.Decode(&person)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("FailedToDeserialiseDto")
		w.WriteHeader(500)
		return
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("FailedToConnectToDatabase")
		w.WriteHeader(500)
		return
	}

	db.Create(&Member{ID: uuid.New(), Age: person.Age, Name: person.Name, CreatedAt: time.Now()})

	json.NewEncoder(w).Encode(person)
}

func get(w http.ResponseWriter)  {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("FailedToConnectToDatabase")
		w.WriteHeader(500)
		return
	}

	var persons []Member
	db.Find(&persons)

	json.NewEncoder(w).Encode(persons)
}

func main() {
	runMigrations()
	s := &server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

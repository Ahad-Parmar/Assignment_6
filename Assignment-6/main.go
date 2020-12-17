package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

type Truck struct {
	gorm.Model
	DriverName    string
	CleanerName   string
	TruckNo       int
}

var db *gorm.DB

var err error

func main() {
	router := mux.NewRouter()
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=password")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Truck{})

	router.HandleFunc("/", GetTrucks).Methods("GET")
	router.HandleFunc("/truck/{id}", GetTruck).Methods("GET")
	router.HandleFunc("/truck/{id}", DeleteTruck).Methods("DELETE")
	router.HandleFunc("/addtruck", AddTruck).Methods("POST")
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func GetTrucks(w http.ResponseWriter, r *http.Request) {
	var trucks []truck
	db.Find(&trucks)
	json.NewEncoder(w).Encode(&trucks)
}

func GetTruck(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var truck Truck
	db.First(&truck, params["id"])
	json.NewEncoder(w).Encode(&truck)
}
func DeleteTruck(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var truck Truck
	db.First(&truck, params["id"])
	db.Delete(&truck)
	var trucks []Truck
	db.Find(&trucks)
	json.NewEncoder(w).Encode(&trucks)
}

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var truck Truck
	json.NewDecoder(r.Body).Decode(&truck)
	db.Create(&truck)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(truck)
}

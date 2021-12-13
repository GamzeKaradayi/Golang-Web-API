package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gamzekaradayi/book/models"
	"github.com/gorilla/mux"
)

var books []models.Book
var apiRoot string

func main() {

	apiRoot := "/api"

	books = []models.Book{
		{ID: 1, Name: "Serenad", Writer: "Zülfü Livaneli", Category: "Aşk", InStock: true},
		{ID: 2, Name: "Beyoğlu Rapsodisi", Writer: "Ahmet Ümit", Category: "Polisiye", InStock: true},
		{ID: 3, Name: "Körlük", Writer: "Jose Saramago", Category: "Bilim", InStock: false},
		{ID: 4, Name: "Fareler ve İnsanlar", Writer: "John Steinbeck", Category: "Trajedi", InStock: true},
	}

	myRouter := mux.NewRouter()

	//GET http://localhost:8080/api/books
	myRouter.HandleFunc(apiRoot+"/books", getBooks)

	//GET http://localhost:8080/api/book/id
	myRouter.HandleFunc(apiRoot+"/book/{id:[0-9]}", getSingleBook).Methods("GET")

	//POST http://localhost:8080/api/book
	myRouter.HandleFunc(apiRoot+"/book", createNewBook).Methods("POST")

	//PUT http://localhost:8080/api/book/id
	myRouter.HandleFunc(apiRoot+"/book/{id}", updateBook).Methods("PUT")

	//DELETE http://localhost:8080/api/book/id
	myRouter.HandleFunc(apiRoot+"/book/{id}", deleteBook).Methods("DELETE")

	http.Handle("/", myRouter)
	http.ListenAndServe(":8080", nil)
}

//bring the list of books
func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func getSingleBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := strconv.Atoi(vars["id"])

	//return the book encoded as JSON
	for _, book := range books {
		if book.ID == key {
			json.NewEncoder(w).Encode(book)
		}
	}

}

func createNewBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	// get the body of our POST request
	err := json.NewDecoder(r.Body).Decode(&book)
	checkError(err)

	book.ID = rand.Intn(10000)
	books = append(books, book)

	//api response message
	apiResult := models.Api{"Successfully added!", false}
	output, err := json.Marshal(apiResult)
	checkError(err)
	fmt.Fprintf(w, string(output))
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	vars := mux.Vars(r)
	updatedBookId, _ := strconv.Atoi(vars["id"])

	for index, item := range books {
		if item.ID == updatedBookId {
			books = append(books[:index], books[index+1:]...)
			err := json.NewDecoder(r.Body).Decode(&book)
			checkError(err)

			book.ID = updatedBookId
			books = append(books, book)
			json.NewEncoder(w).Encode(books)

			//api response message
			apiResult := models.Api{"Updated!", false}
			output, err := json.Marshal(apiResult)
			checkError(err)
			fmt.Fprintf(w, string(output))
		}
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for index, book := range books {
		if book.ID == id {
			//book array updated after deletion
			books = append(books[:index], books[index+1:]...)

			//api response message
			apiResult := models.Api{"Deleted!", false}
			output, err := json.Marshal(apiResult)
			checkError(err)
			fmt.Fprintf(w, string(output))
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal Error : ", err.Error())
		os.Exit(1)
	}
}

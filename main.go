package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gamzekaradayi/book/models"
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

	http.HandleFunc(apiRoot+"/books", getBooks)
	http.HandleFunc(apiRoot+"/books/", getBookByCategoryName)
	http.HandleFunc(apiRoot+"/books/add", addBook)
	http.HandleFunc(apiRoot+"/books/delete/", deleteBook)
	http.HandleFunc(apiRoot+"/books/update/", updateBook)

	http.ListenAndServe(":8080", nil)
}

//bring the list of books
func getBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		output, err := json.Marshal(books)
		checkError(err)
		fmt.Fprintf(w, string(output))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//bring the list of books by category name
func getBookByCategoryName(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var result []models.Book
		url := r.URL.Path
		parts := strings.Split(url, "/")

		for _, p := range books {
			if strings.EqualFold(p.Category, parts[3]) {
				result = append(result, p)
			}
		}

		if len(result) == 0 {
			http.Error(w, http.StatusText(404), 404)
		} else {
			output, err := json.Marshal(result)
			checkError(err)
			fmt.Fprintf(w, string(output))
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//add a book
func addBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		checkError(err)

		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		writer := r.FormValue("writer")
		category := r.FormValue("category")
		instock, _ := strconv.ParseBool(r.FormValue("instock"))

		bookItem := models.Book{ID: id, Name: name, Writer: writer, Category: category, InStock: instock}
		books = append(books, bookItem)

		apiResult := models.Api{"Successfully added!", false}
		output, err := json.Marshal(apiResult)
		checkError(err)
		fmt.Fprintf(w, string(output))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		deletedBookId := r.URL.Path[len(r.URL.Path)-1:]
		deleted, _ := strconv.Atoi(deletedBookId)
		for index, book := range books {
			if book.ID == deleted {
				books = append(books[:index], books[index+1:]...)
				apiResult := models.Api{"Deleted!", false}
				output, err := json.Marshal(apiResult)
				checkError(err)
				fmt.Fprintf(w, string(output))
			}
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var updatedBook models.Book
		reqBody, err := ioutil.ReadAll(r.Body)
		checkError(err)

		//get id value
		updatedBookId := r.URL.Path[len(r.URL.Path)-1:]
		updated, _ := strconv.Atoi(updatedBookId)

		json.Unmarshal(reqBody, &updatedBook)

		for i, book := range books {
			if book.ID == updated {
				book.Name = updatedBook.Name
				book.Writer = updatedBook.Writer
				book.Category = updatedBook.Category
				book.InStock = updatedBook.InStock
				books = append(books[i:], book)

				apiResult := models.Api{"Updated!", false}
				output, err := json.Marshal(apiResult)
				checkError(err)
				fmt.Fprintf(w, string(output))
			}
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal Error : ", err.Error())
		os.Exit(1)
	}
}

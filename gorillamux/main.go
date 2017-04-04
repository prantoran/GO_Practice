package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//HomeHandler blah
func HomeHandler(w http.ResponseWriter, r *http.Request) {

}

//ProductHandler blah
func ProductHandler(w http.ResponseWriter, r *http.Request) {

}

//ArticlesHandler blah
func ArticlesHandler(w http.ResponseWriter, r *http.Request) {

}

//ArticlesCategoryHandler blah
func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/products", ProductHandler)
	r.HandleFunc("/articles", ArticlesHandler)

	//paths can have variables, if expression pattern is not defined, the matched variable
	//will be anything until the next slash

	r.HandleFunc("/products/{key}", ProductHandler)
	r.HandleFunc("/articles/{category}/", ArticlesHandler)
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticlesHandler)
	r.HandleFunc("/articles/{category}/{sort:(?:asc|desc|new)}", ArticlesCategoryHandler)

	http.Handle("/", r)

	req := http.Request{}
	req.Method = "GET"

	vars := mux.Vars(&req)
	fmt.Printf("vars of req: %v \n", vars)
	fmt.Printf("Method of req: %v \n", req.Method)
	//category := vars["category"]

	//Subrouting
	r2 := mux.NewRouter()
	s := r2.Host("databd.slack.com").Subrouter()
	s.HandleFunc("/products/", ProductHandler)
	s.HandleFunc("/products/{key}", ProductHandler)
	s.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticlesHandler)

	//routes can be named
	r3 := mux.NewRouter()
	r3.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticlesHandler).Name("article")

	//to build a URL, get the route and call the URL () method
	//passing a sequence of key/value pairs for the route variables
	url, err := r3.Get("article").URL("category", "technology", "id", "42")
	fmt.Printf("pass\n")
	if err != nil {
		panic(err)
	}
	fmt.Printf("url.EscapedPath(): %v \n", url.EscapedPath())

}

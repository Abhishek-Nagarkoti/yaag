package main

import (
	"fmt"
	"github.com/akshaykumar12527/yaag/middleware"
	"github.com/akshaykumar12527/yaag/yaag"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, time.Now().String())
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = ioutil.ReadAll(r.Body)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("test", "tesasasdasd")
	fmt.Fprintf(w, time.Now().String())
}

func main() {
	yaag.Init(&yaag.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.HandleFunc(handler))
	r.HandleFunc("/testing", middleware.HandleFunc(postHandler)).Methods("POST")
	http.ListenAndServe(":8080", r)
}

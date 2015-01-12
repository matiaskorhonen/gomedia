package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var (
	bucketName string
	baseURL    string
)

func init() {
	flag.StringVar(&bucketName, "b", "", "Bucket Name")
	flag.StringVar(&baseURL, "u", "", "Base URL")
}

func rootHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Nothing to see here!")
}

func tweetbot(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(r.ContentLength)

	// message := r.Form["message"]
	// source := r.Form["source"]
	file, header, err := r.FormFile("media")
	defer file.Close()

	if err != nil {
		panic(err.Error())
	}

	timeStamp := time.Now().Unix()
	random := rand.Intn(999999)
	filename := fmt.Sprintf("%x-%x-%s", timeStamp, random, header.Filename)

	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}

	// Open Bucket
	s := s3.New(auth, aws.EUWest)
	bucket := s.Bucket(bucketName)

	path := fmt.Sprintf("tweetbot/%s", filename)

	contentLength, err := strconv.Atoi(header.Header.Get("Content-Length"))
	if err != nil {
		panic(err.Error())
	}

	buffer := make([]byte, contentLength)
	cBytes, err := file.Read(buffer)

	s3Headers := map[string][]string{"Content-Type": {header.Header.Get("Content-Type")}, "Cache-Control": {"public, max-age=315360000"}}

	err = bucket.PutHeader(path, buffer[0:cBytes], s3Headers, s3.PublicRead)
	if err != nil {
		panic(err.Error())
	}

	url := fmt.Sprintf("%s/%s", baseURL, path)

	fmt.Printf("\nFile %s (%s) uploaded successfully.\n", header.Filename, header.Header)

	responseMap := map[string]string{"url": url}
	jsonResponse, _ := json.Marshal(responseMap)
	fmt.Println(string(jsonResponse))

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonResponse))
}

func main() {
	goji.Get("/", rootHandler)
	re := regexp.MustCompile("/tweetbot")
	goji.Post(re, tweetbot)
	goji.Serve()
}

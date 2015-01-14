package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
	"github.com/goji/httpauth"
	"github.com/tv42/base58"
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
	multiReader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Open Bucket
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}
	s := s3.New(auth, aws.EUWest)
	bucket := s.Bucket(bucketName)

	// We only want the first part, the media
	part, err := multiReader.NextPart()
	if err == io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	originalFilename := part.FileName()
	fileExt := filepath.Ext(originalFilename)
	unixTime := time.Now().UTC().Unix()
	b58buf := base58.EncodeBig(nil, big.NewInt(unixTime))

	filename := fmt.Sprintf("%s%s", b58buf, fileExt)
	path := "tweetbot/" + filename

	contentType := part.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(fileExt)

		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	contentLength, err := strconv.ParseInt(part.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bucket.PutReader(path, part, contentLength, contentType, s3.PublicRead, s3.Options{CacheControl: "public, max-age=315360000"})
	if err != nil {
		panic(err.Error())
	}

	// fmt.Printf("\nFile %s (%s) uploaded successfully.\n", originalFilename, path)

	url := fmt.Sprintf("%s/%s", baseURL, path)

	responseMap := map[string]string{"url": url}
	jsonResponse, _ := json.Marshal(responseMap)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonResponse))
}

func main() {
	goji.Use(httpauth.SimpleBasicAuth("x", "password"))

	goji.Get("/", rootHandler)

	re := regexp.MustCompile("/tweetbot")
	goji.Post(re, tweetbot)

	goji.Serve()
}

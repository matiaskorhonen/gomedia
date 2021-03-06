package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/matiaskorhonen/gomedia/Godeps/_workspace/src/github.com/crowdmob/goamz/aws"
	"github.com/matiaskorhonen/gomedia/Godeps/_workspace/src/github.com/crowdmob/goamz/s3"
	"github.com/matiaskorhonen/gomedia/Godeps/_workspace/src/github.com/tobi/airbrake-go"
	"github.com/matiaskorhonen/gomedia/Godeps/_workspace/src/github.com/tv42/base58"
	"github.com/matiaskorhonen/gomedia/Godeps/_workspace/src/github.com/zenazn/goji"
	"github.com/matiaskorhonen/gomedia/Godeps/_workspace/src/github.com/zenazn/goji/web"
)

var (
	bucketName  string
	baseURL     string
	awsRegion   aws.Region
	username    string
	password    string
	useAirbrake bool
)

func init() {
	bucketName = os.Getenv("BUCKET_NAME")
	baseURL = os.Getenv("BASE_URL")
	username = os.Getenv("HTTP_USER")
	password = os.Getenv("HTTP_PASSWORD")

	if os.Getenv("AWS_REGION") == "" {
		awsRegion = aws.GetRegion("us-east-1")
	} else {
		awsRegion = aws.GetRegion(os.Getenv("AWS_REGION"))
		if awsRegion.Name == "" {
			fmt.Printf("Unknown AWS region: " + os.Getenv("AWS_REGION") + "\n")
			os.Exit(1)
		}
	}

	// Configure Airbrake/Errbit
	airbrakeAPIKey := os.Getenv("AIRBRAKE_API_KEY")
	airbrakeEndpoint := os.Getenv("AIRBRAKE_ENDPOINT")
	airbrakeEnvironment := os.Getenv("AIRBRAKE_ENVIRONMENT")
	useAirbrake = false
	if airbrakeAPIKey != "" {
		airbrake.ApiKey = airbrakeAPIKey
		useAirbrake = true
		if airbrakeEndpoint != "" {
			airbrake.Endpoint = airbrakeEndpoint
		}
		if airbrakeEnvironment != "" {
			airbrake.Environment = airbrakeEnvironment
		}
	}
}

func reportIfAirbake(err error, r *http.Request) {
	if useAirbrake {
		if r != nil {
			airbrake.Error(err, r)
		} else {
			airbrake.Notify(err)
		}
	}
}

func openBucket() (*s3.Bucket, error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return nil, err
	}
	s := s3.New(auth, awsRegion)

	if bucketName == "" {
		return nil, errors.New("BUCKET_NAME not set")
	}

	bucket := s.Bucket(bucketName)
	return bucket, nil
}

func readerToS3(ioReader io.Reader, basePath string, originalFilename string, generateNewFileName bool, contentType string, contentLength int64) (string, error) {
	bucket, err := openBucket()
	if err != nil {
		return "", err
	}

	fileExt := filepath.Ext(originalFilename)

	var filename string

	if generateNewFileName {
		unixTime := time.Now().UTC().Unix()
		b58buf := base58.EncodeBig(nil, big.NewInt(unixTime))
		filename = fmt.Sprintf("%s%s", b58buf, fileExt)
	} else {
		filename = originalFilename
	}

	path := basePath + filename

	if contentType == "" || contentType == "application/octet-stream" {
		contentType = mime.TypeByExtension(fileExt)

		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	err = bucket.PutReader(path, ioReader, contentLength, contentType, s3.PublicRead, s3.Options{CacheControl: "public, max-age=315360000"})
	if err != nil {
		return "", err
	}

	var url string

	if baseURL == "" {
		url = bucket.URL(path)
	} else {
		url = fmt.Sprintf("%s/%s", baseURL, path)
	}

	return url, nil
}

func uploadPartToS3(part *multipart.Part, basePath string) (string, error) {
	originalFilename := part.FileName()

	contentType := part.Header.Get("Content-Type")

	contentLength, err := strconv.ParseInt(part.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return "", err
	}

	url, err := readerToS3(part, basePath, originalFilename, true, contentType, contentLength)

	return url, err
}

func rootHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Nothing to see here!")
}

func tweetbot(c web.C, w http.ResponseWriter, r *http.Request) {
	multiReader, err := r.MultipartReader()
	if err != nil {
		reportIfAirbake(err, r)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var part *multipart.Part
	err = nil

	// Find the media part
	for err == nil {
		part, err = multiReader.NextPart()

		if err != nil {
			reportIfAirbake(err, r)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if part.FormName() == "media" {
			break
		}
	}

	// Ensure that the Content-Length is set
	_, err = strconv.ParseInt(part.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		reportIfAirbake(err, r)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url, err := uploadPartToS3(part, "tweetbot/")
	if err != nil {
		reportIfAirbake(err, r)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseMap := map[string]string{"url": url}
	jsonResponse, _ := json.Marshal(responseMap)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonResponse))
}

func webDavUpload(c web.C, w http.ResponseWriter, r *http.Request) {
	// Ensure that the Content-Length is set
	if r.ContentLength < 1 {
		http.Error(w, "Content-Length must be set", http.StatusBadRequest)
		return
	}

	originalFilename := c.URLParams["name"]

	contentType := r.Header.Get("Content-Type")

	basePath := ""

	url, err := readerToS3(r.Body, basePath, originalFilename, false, contentType, r.ContentLength)

	if err != nil {
		reportIfAirbake(err, r)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusCreated)
}

func webDavDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Deleting files is not supported", http.StatusNotImplemented)
}

func propfindIntercept(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PROPFIND" {
			xml := "<?xml version=\"1.0\" ?>\n" +
				"<D:multistatus xmlns:D=\"DAV:\">\n" +
				"<D:response>\n" +
				"<D:href>http://www.contoso.com/public/container/</D:href>\n" +
				"<D:propstat>\n" +
				"<D:status>HTTP/1.1 200 OK</D:status>\n" +
				"</D:propstat>\n" +
				"</D:response>\n" +
				"</D:multistatus>\n"

			w.WriteHeader(207)
			w.Header().Set("Content-Type", "text/xml")
			w.Write([]byte(xml))
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func main() {
	if username != "" && password != "" {
		authOpts := AuthOptions{
			Realm:    "Restricted",
			User:     username,
			Password: password,
		}

		goji.Use(BasicAuth(authOpts))
	}

	goji.Use(propfindIntercept)

	goji.Get("/", rootHandler)

	goji.Put("/:name", webDavUpload)
	goji.Delete("/:name", webDavDelete)

	re := regexp.MustCompile(`\A/tweetbot/?\z`)
	goji.Post(re, tweetbot)

	goji.Serve()
}

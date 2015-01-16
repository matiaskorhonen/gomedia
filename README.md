# GoMedia

Tweetbot custom media endpoint for uploads to S3

## Install/Build/Deploy

### Cross compile for Heroku

```sh
go get github.com/laher/goxc
goxc -os="linux" -arch="amd64" -d=./build -tasks-="downloads-page,deb,deb-dev,go-test,go-vet"
```

### Deploy with Heroku Anvil

Install the `heroku-anvil` plugin

```sh
heroku plugins:install https://github.com/ddollar/heroku-anvil
```

Cross compile the app as described above, then build and release the slug using heroku-anvil.

```sh
# Extract the binary to the slug directory
tar -xf ./build/snapshot/gomedia_linux_amd64.tar.gz -C ./slug --strip-components=1

# Compile slug and release
heroku build ./slug -b https://github.com/ryandotsmith/null-buildpack.git -r gomedia
```

(replace **gomedia** with the name of your Heroku app)

See below for the required environment variables on Heroku.

### Everything together

Just run `make deploy APP=myherokuappname`

### Environment Variables / Configuration

* `AWS_ACCESS_KEY_ID` — (required) self-explanatory
* `AWS_SECRET_ACCESS_KEY` — (required) self-explanatory
* `AWS_REGION` — (optional) the AWS region. Defaults to `us-east-1`
* `BASE_URL` — (required) the base URL for generated URLs (no trailing slash)
    * Use the bucket URL if you don't have a CNAME
    * If you're using a custom CNAME (or a CDN) set it to the hostname (beginning with http or https)
* `BUCKET_NAME` — (required) the name of your S3 bucket
* `HTTP_PASSWORD` — (recommended) protect the upload end points with basic auth
* `HTTP_USER` — (recommended) protect the upload end points with basic auth

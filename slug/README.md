# GoMedia

Tweetbot custom media endpoint for uploads to S3

## Install/Build/Deploy

### Cross compile for Heroku

```sh
go get github.com/laher/goxc
goxc -os="linux" -arch="amd64" -d=./build
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
heroku build ./slug -r gomedia
```

(replace **gomedia** with the name of your Heroku app)

### Everything together

```sh
goxc -os="linux" -arch="amd64" -d=./build &&
tar -xf ./build/snapshot/gomedia_linux_amd64.tar.gz -C ./slug --strip-components=1 &&
heroku build ./slug -r gomedia
```

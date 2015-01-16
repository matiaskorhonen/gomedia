build:
	goxc -os="linux" -arch="amd64" -d=./build -tasks-="downloads-page,deb,deb-dev,go-test,go-vet"

extractbin:
	tar -xf ./build/snapshot/gomedia_linux_amd64.tar.gz -C ./slug --strip-components=1

deploy: build extractbin
	heroku build ./slug -r ${APP}

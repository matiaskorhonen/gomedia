# GoMedia

Tweetbot custom media endpoint for uploads to S3

## Install/Build/Deploy

### Configure the buildpack

If you're creating a new heroku app, set the buildpack on creation:

```sh
heroku create -b https://github.com/kr/heroku-buildpack-go.git --region=eu
```

For an existing Heroku app, set the `BUILDPACK_URL`:

```sh
heroku config:set BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git
```

### Set environment variables

GoMedia uses environment variables for configuration (see all available options below).

Before pushing to Heroku, you should at the very least set the required options, e.g.

```sh
heroku config:set AWS_ACCESS_KEY_ID="…" AWS_SECRET_ACCESS_KEY="…"  \
AWS_REGION="eu-west-1" BUCKET_NAME="mybucket"
```

It is highly recommended that you also set a username and password (for the obvious reasons).

### Push to Heroku

```sh
git push heroku master
```

### Environment Variables / Configuration

* `AWS_ACCESS_KEY_ID` — (required) self-explanatory
* `AWS_SECRET_ACCESS_KEY` — (required) self-explanatory
* `AWS_REGION` — (optional) the AWS region. Defaults to `us-east-1`
* `BUCKET_NAME` — (required) the name of your S3 bucket
* `BASE_URL` — (optional) the base URL for generated URLs (no trailing slash)
    * Defaults to the bucket url (e.g. `https://s3-eu-west-1.amazonaws.com/BUCKET_NAME/`)
    * If you're using a custom CNAME (or a CDN) set it to the hostname (beginning with http or https)
* `HTTP_PASSWORD` — (recommended) protect the upload end points with basic auth
* `HTTP_USER` — (recommended) protect the upload end points with basic auth
* `AIRBRAKE_API_KEY` — (optional) set this to your Airbrake API key if you want airbrake support
* `AIRBRAKE_ENDPOINT` — (optional) configure this if you use Errbit (or some other custom endpoint)
    * e.g. `https://myerribt.herokuapp.com/notifier_api/v2/notices`
* `AIRBRAKE_ENVIRONMENT` — (optional) set the environment for Airbrake, defaults to `development`

## Use with Tweetbot

GoMedia works as a [Custom Media Upload][custom] endpoint for Tweetbot (both iOS and OS X).

All you need to do is set the custom endpoint URL to:

```
https://MYAPPNAME.herokuapp.com/tweetbot?username=USERNAME&password=PASSWORD
```

Replace `USERNAME` and `PASSWORD` with your own credentials and `MYAPPNAME` with your Heroku app name…

![Tweetbot for Mac configuration][tweetbot-mac]

![Tweetbot for iOS configuration][tweetbot-ios]

[custom]: http://tapbots.net/tweetbot/custom_media/
[tweetbot-mac]: http://shots.matiaskorhonen.fi/tweetbot-mac-gomedia_nibjww.png
[tweetbot-ios]: http://shots.matiaskorhonen.fi/tweetbot-ios-gomedia_nibjwd.png

## Use with Monosnap

For monosnap you need to configure a WebDAV endpoint with the app URL and credentials. Do **not** set the port!

![GoMedia Monosnap configuration][monoconf]

You need to configure the Base URL in Monosnap regardless of whether you're using a custom base URL or not.

### Recommended

With Monosnap it is recommended to add `%R` to the filename template because the filename cannot be altered server side (thus there's a small chance that the filename will conflict with an existing file in the bucket).

![Advanced Monosnap configuration][advanced]

[monoconf]: http://shots.matiaskorhonen.fi/monosnap-gomedia-configuration_nibjb9.png
[advanced]: http://shots.matiaskorhonen.fi/Monosnap_advanced.png

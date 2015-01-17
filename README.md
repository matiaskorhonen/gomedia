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

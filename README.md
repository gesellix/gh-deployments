# gh-deployments

A little cli on top of go-github.

Official docs:
- https://docs.github.com/en/free-pro-team@latest/rest/guides/delivering-deployments
- https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#create-a-deployment

## Create Deployment

````shell script
$ docker run --rm -it \
  --env GITHUB_TOKEN \
  --env GITHUB_OWNER=gesellix \
  --env GITHUB_REPO=gh-deployments \
  deploy create \
  --ref=7d9c662978d50faf7a0ba489fcc94e543f49da61 \
  --description=ein\ test \
  --environment=test
deployment_id=279713949
````

## Update Deployment Status

````shell script
$ docker run --rm -it \
  --env GITHUB_TOKEN \
  --env GITHUB_OWNER=gesellix \
  --env GITHUB_REPO=gh-deployments \
  deploy status \
  --deployment-id=279713949 \
  --log_url=https://www.gesellix.net \
  --state=success \
  --description=ein\ test \
  --environment=test
deployment_status_id=413547422
````

# Drone Lambda Plugin
## Disclaimer
Has been forked from the original [plugin](https://github.com/omerxx/drone-lambda-plugin), as it currently only supports us-east-1 region. Following features in addition:
* AWS Lambda region is configurable with function_region within drone plugin configuration (see example below)
* Actually fails when something goes wrong (exit code 1, which will let the build fail) 

### The plugin utilizes AWS go-sdk to update an existing function's code; build your code, zip it with dependencies and upload it to S3. Then trigger the plugin for deploy.

## Build:
Build the binary:
```
go build main.go
```

## Docker:
Build the container:
```
docker build --rm=true -t danielhill/drone-lambda-plugin .
```

## Usage:

#### Execute from the working directory; 
This will update `my-function` with a zip file under `S3://some-bucket/lambda-dir/lambda-project-1.zip`:
```bash
docker run --rm \
  -e PLUGIN_FUNCTION_NAME=my-function \
  -e PLUGIN_S3_BUCKET=some-bucket \
  -e PLUGIN_FILE_NAME=lambda-directory/lambda-project-1.zip \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  --privileged \
  plugins/docker --dry-run
```

#### Example:

```yaml
pipeline:
  deploy-lambda:
    image: danielhill/drone-lambda-plugin
    pull: true
    function_name: my-function
    s3_bucket: some-bucket
    file_name: lambda-dir/lambda-project-${DRONE_BUILD_NUMBER}.zip
```

#### Example of a complete Lambda project's pipeline:

```yaml
pipeline:
  build:
    image: python:2.7-alpine
    commands:
      - apk update && apk add zip
      - pip install -r requirements.txt -t .
      - zip -r -9 lambda-project-${DRONE_BUILD_NUMBER}.zip *

  s3-publish:
    image: plugins/s3
    acl: private
    region: us-east-1
    bucket: some-bucket
    target: lambda-dir
    source: lambda-project-${DRONE_BUILD_NUMBER}.zip

  deploy-lambda:
    image: danielhill/drone-lambda-plugin
    pull: true
    function_name: my-function
    function_region: eu-central-1
    s3_bucket: some-bucket
    file_name: lambda-dir/revenue-report-${DRONE_BUILD_NUMBER}.zip

  notify-slack-releases:
    image: plugins/slack
    channel: product-releases
    webhook: https://hooks.slack.com/services/ABCD/XYZ
    username: Drone-CI
```

# Serverless Multicloud

Example architecture demonstrating the use of a thin boundary layer to abstract the serverless
(or other) runtime environment in an API application.

## Make Targets

* **aws-deploy** - Compile and deploy to AWS with API Gateway and Lambda
* **gcp-func-deploy** - Compile and deploy using GCP Functions
* **run-local** - Run local http server process
* **docker-run** - Build docker image and run locally 

**Note:** `aws-deploy` requires `CODE_BUCKET` environment variable to be set 
to an S3 bucket name where code will be stored. In addition, the [AWS CLI](https://aws.amazon.com/cli/) environment
will need to be installed and configured using `aws configure` to access the account and bucket.

**Note:** `gcp-func-deploy` requires [gcloud CLI](https://cloud.google.com/sdk/gcloud/) tool to be installed and configured.
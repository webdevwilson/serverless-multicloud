HTTP_SERVER := dist/http
HTTP_SERVER_SRC := ./src/boundary/http/main.go

LAMBDA_HANDLER := dist/lambda_handler
LAMBDA_HANDLER_SRC := ./src/boundary/lambda/main.go

SAM_STACK_NAME := serverless-multicloud-site
SAM_TEMPLATE := aws-lambda.yaml
SAM_PKG_TEMPLATE := dist/template.pkg.yaml

GCP_ZIP := dist/gcp_func.zip
GCP_SRC := dist/gcp
GCP_BUCKET := serverless-multicloud-src

clean:
	rm -rf dist

dist:
	mkdir dist

$(SAM_PKG_TEMPLATE): dist
	aws cloudformation package \
		--template-file $(SAM_TEMPLATE) \
		--s3-bucket $(S3_BUCKET) \
		--output-template-file $(SAM_PKG_TEMPLATE)

# sam aws lambda targets
$(LAMBDA_HANDLER): dist
	GOOS=linux GOARCH=amd64 go build -o $(LAMBDA_HANDLER) $(LAMBDA_HANDLER_SRC)

aws-deploy: all
	aws cloudformation deploy \
		--template-file $(SAM_PKG_TEMPLATE) \
		--stack-name $(SAM_STACK_NAME) \
		--capabilities CAPABILITY_IAM

all: $(LAMBDA_HANDLER) $(SAM_PKG_TEMPLATE)

serve:
	go build -o $(LAMBDA_HANDLER) $(LAMBDA_HANDLER_SRC)
	sam-cli local

# Docker targets
docker-build:
	docker build -t serverless-multicloud .

docker-run: docker-build
	docker run -e MSG=local-docker -it --rm --publish 8080:8080 --name serverless-multicloud serverless-multicloud

# local http server
$(HTTP_SERVER): dist
	go build -o $(HTTP_SERVER) $(HTTP_SERVER_SRC)

run-local:
	MSG=local-http go run src/server/main.go

# gcp function targets
$(GCP_SRC): dist
	mkdir -p $(GCP_SRC)
	cp -r ./src $(GCP_SRC)
	cp ./src/boundary/gcpfunc/function.go $(GCP_SRC)/
	cp go.mod $(GCP_SRC)

$(GCP_ZIP): dist
	zip -r $(GCP_ZIP) ./src
	cp ./src/boundary/http/main.go dist/function.go
	zip -rj $(GCP_ZIP) ./dist/function.go

gcp-upload: $(GCP_ZIP)
	gsutil cp $(GCP_ZIP) gs://$(GS_BUCKET)

gcp-func-deploy: $(GCP_SRC)
	gcloud functions deploy serverless-multicloud \
		--memory=128MB \
		--entry-point=Handler \
		--runtime=go111 \
		--source=$(GCP_SRC) \
		--set-env-vars=MSG=gcp-func \
		--trigger-http


.PHONY: all aws-deploy clean serve docker-build docker-run deploy-gcp-func run-local


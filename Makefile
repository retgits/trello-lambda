.PHONY: deps clean build deploy test-lambda

deps:
	go get -u ./...

clean: 
	rm -rf ./bin
	
build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/trello-lambda *.go

test-lambda: clean build
	sam local invoke Trello -e ./test/event.json

deploy: clean build
	sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket retgits-apps
	sam deploy --template-file packaged.yaml --stack-name trello-lambda --capabilities CAPABILITY_IAM
# trello-lambda - A serverless app to create Trello cards

This serverless function is designed to create a new Trello Cards each time it gets invoked.

## Layout
```bash
.                    
├── test            
│   └── event.json      <-- Sample event to test using SAM local
├── .gitignore          <-- Ignoring the things you do not want in git
├── function.go         <-- Test main function code
├── LICENSE             <-- The license file
├── main.go             <-- The Flogo Lambda trigger code
├── Makefile            <-- Makefile to build and deploy
├── README.md           <-- This file
└── template.yaml       <-- SAM Template
```

## Installing
There are a few ways to install this project

### Get the sources
You can get the sources for this project by simply running
```bash
$ go get -u github.com/retgits/github-trello/...
```

### Deploy
Deploy the Lambda app by running
```bash
$ make deploy
```

## Parameters
### AWS Systems Manager parameters
The code will automatically retrieve the below list of parameters from the AWS Systems Manager Parameter store:

* **/trello/appkey**: Your Trello App token
* **/trello/apptoken**: Your Trello App key
* **/trello/list**: The ID of the list you want to send the card to

_Details on how to get the `appkey` and `apptoken` for Trello can be found in the [Trello API documentation](https://trello.readme.io/docs/get-started)_

### Deployment parameters
In the `template.yaml` there are certain deployment parameters:

* **region**: The AWS region in which the code is deployed

## Make targets
trello-lambda has a _Makefile_ that can be used for most of the operations

```
usage: make [target]
```

* **deps**: Gets all dependencies for this app
* **clean** : Removes the dist directory
* **build**: Builds an executable to be deployed to AWS Lambda
* **test-lambda**: Clean, builds and tests the code by using the AWS SAM CLI
* **deploy**: Cleans, builds and deploys the code to AWS Lambda

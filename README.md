# Trello app for Lambda

This serverless function is designed to create a new Trello Cards each time it gets invoked.

## Layout
```bash
.
├── .travis.yml                 <-- Travis-CI build file
├── event.json                  <-- Sample event to test using SAM local
├── README.md                   <-- This file
├── src                         <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda trigger code
│   └── function.go             <-- Lambda function code
└── template.yaml               <-- SAM Template
```

## Build and Deploy
Building and deploying this function is done through Travis-CI using [lambda-builder](https://github.com/retgits/lambda-builder)

## AWS Systems Manager
Within the AWS Systems Manager Parameter store there are three parameters that are used in this app:

* /trello/appkey
* /trello/apptoken
* /trello/list

_Details on how to get the `appkey` and `apptoken` for Trello can be found in the [Trello API documentation](https://trello.readme.io/docs/get-started)_

## TODO
- [ ] Remove `shim_support.go` and `shim.go` as soon as the shim generation support gets into the master branch of [flogo-lib](https://github.com/TIBCOSoftware/flogo-lib)


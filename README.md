# Trello app for Lambda

This Serverless function is designed to create a new Trello Cards each time it gets invoked.

## Layout
```bash
.
├── build.sh                    <-- Make to automate build
├── event.json                  <-- Sample event to test using SAM local
├── README.md                   <-- This file
├── src                         <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
└── template.yaml               <-- SAM Template
```

## build.sh
The `build.sh` file has seven commands to make working with this app easier than it already is

* deps: go get and update all the dependencies
* clean: removes the ./bin folder
* test: uses SAM local and the event in `event.json` to test the implementation
* build: creates the executable
* getparams: updates the SAM template with the values from the AWS Systems Manager Parameter Store
* delparams: removes the values of the environment variables in the SAM template
* deploy: deploy the function to AWS Lambda

## Prerequisites
While executing the the build script there are a few programs that are used:

* [jq](https://stedolan.github.io/jq/)
* [yq](https://github.com/mikefarah/yq)
* [aws cli](https://github.com/aws/aws-cli)
* [sam cli](https://github.com/awslabs/aws-sam-cli)

## AWS Systems Manager
Within the AWS Systems Manager Parameter store there are three parameters that are used in this app:

* /trello/appkey
* /trello/apptoken
* /trello/defaultlist

_Details on how to get the `appkey` and `apptoken` for Trello can be found in the [Trello API documentation](https://trello.readme.io/docs/get-started)_

## TODO
- [ ] Update the `deps` target in build.sh to make use of dep or simply have a smarter approach than list all dependencies

## License
The MIT License (MIT)

Copyright (c) 2018 retgits

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
# trello-lambda

[![Go Report Card](https://goreportcard.com/badge/github.com/retgits/trello-lambda?style=flat-square)](https://goreportcard.com/report/github.com/retgits/trello-lambda)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/retgits/trello-lambda)
![GitHub](https://img.shields.io/github/license/retgits/trello-lambda?style=flat-square)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/retgits/trello-lambda?sort=semver&style=flat-square)

> Trello-Lambda is an opinionated stack of Lambda functions to help automate Trello.

I use [Trello](https://trello.com) all the time and instead of manually copy/pasting and clearing out cards, I built a set of serverless apps to help me with that.

## Available functions

* [Weekly Archiver](./weekly-archiver): Archives all cards from the Done list
* [Five Minute Journal](./five-minute-journal): Creates a new [Five Minute Journal](https://www.intelligentchange.com/collections/all/products/the-five-minute-journal) card every day

## Prerequisites

* [Go (at least Go 1.12)](https://golang.org/dl/)
* [AWS](https://portal.aws.amazon.com/) account, with access to Lambda, SSM and KMS

### AWS Systems Manager Parameter Store

The app relies on [AWS Systems Manager Parameter Store](https://aws.amazon.com/systems-manager/features/) (SSM) to store encrypted variables on how to connect to Trello. The variables it relies on are:

* `/<stage>/global/kmskey`: The ID of the KMS key used to decrypt environment variables.
* `/<stage>/trello/apikey`: The Trello API Key.
* `/<stage>/trello/apptoken`: The Trello app token.
* `/<stage>/trello/lists/main-done`: The ID of the Trello list from which to archive cards.
* `/<stage>/trello/lists/main-today`: The ID of the Trello list to create a new Five Minute Journal card in.

With the _`<stage>`_ variable, you can have different stack referencing different sets of API keys. Details on how to get the `appkey` and `apptoken` for Trello can be found in the [Trello API documentation](https://trello.readme.io/docs/get-started).

### AWS Key Management Service

The parameters are encrypted using [Amazon KMS](https://aws.amazon.com/kms/) and retrieved from the Parameter Store on deployment. This way the encrypted variables are given to the Lambda function and the function needs to take care of decrypting them at runtime. To be able to decrypt the variables at runtime, the Lambda function will need permission to access the KMS Key with the KeyID specified in `/<stage>/global/kmskey`

To create the encrypted variables, run the below command for all of the variables

```bash
aws ssm put-parameter                       \
   --type String                            \
   --name "<your variable>"                 \
   --value $(aws kms encrypt                \
              --output text                 \
              --query CiphertextBlob        \
              --key-id <YOUR_KMS_KEY_ID>    \
              --plaintext "PLAIN TEXT HERE")
```

## Usage

There are several `make` targets available to help build and deploy the function

| Target | Description                                       |
|--------|---------------------------------------------------|
| build  | Build the executable for Lambda                   |
| clean  | Remove all generated files                        |
| deploy | Deploy the app to AWS Lambda                      |
| deps   | Get the Go modules from the GOPROXY               |
| help   | Displays the help for each target (this message). |
| local  | Run SAM to test the Lambda function using Docker  |
| test   | Run all unit tests and print coverage             |

## Contributing

[Pull requests](https://github.com/retgits/trello-lambda/pulls) are welcome. For major changes, please open [an issue](https://github.com/retgits/trello-lambda/issues) first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

See the [LICENSE](./LICENSE) file in the repository

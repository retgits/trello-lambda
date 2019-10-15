# trello-lambda

Trello-Lambda is an opinionated stack of Lambda functions to help automate Trello.

## Prerequisites

### Environment Variables

The app relies on [AWS Systems Manager Parameter Store](https://aws.amazon.com/systems-manager/features/) (SSM) to store encrypted variables on how to connect to Trello. The variables it relies on are:

* `/<stage>/global/kmskey`: The ID of the KMS key used to decrypt environment variables.
* `/<stage>/trello/apikey`: The Trello API Key.
* `/<stage>/trello/apptoken`: The Trello app token.
* `/<stage>/trello/lists/main-done`: The ID of the Trello list from which to archive cards.

With the _`<stage>`_ variable, you can have different stack referencing different sets of API keys. Details on how to get the `appkey` and `apptoken` for Trello can be found in the [Trello API documentation](https://trello.readme.io/docs/get-started).

These parameters are encrypted using [Amazon KMS](https://aws.amazon.com/kms/) and retrieved from the Parameter Store on deployment. This way the encrypted variables are given to the Lambda function and the function needs to take care of decrypting them at runtime. To be able to decrypt the variables at runtime, the Lambda function will need permission to access the KMS Key with the KeyID specified in `/<stage>/global/kmskey`

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

## Build and Deploy

There are several `Make` targets available to help build and deploy the function

| Target | Description                                       |
|--------|---------------------------------------------------|
| build  | Build the executable for Lambda                   |
| clean  | Remove all generated files                        |
| deploy | Deploy the app to AWS Lambda                      |
| deps   | Get the Go modules from the GOPROXY               |
| help   | Displays the help for each target (this message). |
| local  | Run SAM to test the Lambda function using Docker  |
| test   | Run all unit tests and print coverage             |

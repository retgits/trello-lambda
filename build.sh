#!/bin/bash

FUNC=Trello

# Make sure all the dependencies are available
# TODO: Update this to make use of dep
deps() { 
    echo Getting all dependencies...
    go get -u github.com/adlio/trello
    go get -u github.com/aws/aws-lambda-go/...
}

# Remove the bin folder
clean() {
    echo Cleaning up older versions...
    rm -rf bin
}

# Create the executable
build() {
    echo Executing build...
    GOOS=linux GOARCH=amd64 go build -o bin/${FUNC,,} src/*.go
}

# Use SAM local to test the code
test() {
    sam local invoke "${FUNC}" -e event.json
}

# Get the variables from the AWS Systems Manager Parameter Store
getparams() {
    echo Getting parameters from AWS Systems Manager Parameter Store...
    for row in $(yq r template.yaml Resources.${FUNC}.Properties.Environment.Variables); do 
        if [[ $row = *":"* ]]; then 
            param=${row::-1};
            sVar=$(aws ssm get-parameter --name /${FUNC,,}/${param} --with-decryption | jq '.Parameter.Value')
            $(yq w -i template.yaml Resources.${FUNC}.Properties.Environment.Variables.${param} ${sVar})
        fi; 
    done
}

# Replace all parameters in the template with xxx
delparams() {
    echo Removing all parameter values from the template...
    for row in $(yq r template.yaml Resources.${FUNC}.Properties.Environment.Variables); do
        if [[ $row = *":"* ]]; then
            param=${row::-1};
            $(yq w -i template.yaml Resources.${FUNC}.Properties.Environment.Variables.${param} xxx);
        fi;
    done
}

getversion() {
    echo Getting the commit version...
    sVar=$(git log -n 1 --pretty=format:"%H")
    [ ${#sVar} -ge 5 ] && sVar=$sVar || sVar=no-commits
    $(yq w -i template.yaml Resources.${FUNC}.Properties.Tags.commit ${sVar})
}

# Deploy the function to AWS Lambda
deploy() {
    # Make the necessary preparation
    clean
    build
    getparams
    getversion

    # Create a new S3 bucket
    today=`date +%Y%m%d`
    # Get the difference between 53 (which is the max bucket size, 63, minus the length of today and 2 hyphens)
    num=`expr 53 - ${#FUNC,,}`
    suffix=`cat /dev/urandom | tr -dc 'a-z0-9' | fold -w $num | head -n 1`
    bucket=`aws s3 mb s3://${FUNC,,}-$today-$suffix`
    bucket="${bucket:13}"

    # Package it up!
    sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket $bucket

    # Create CF templates and deploy
    sam deploy --template-file packaged.yaml --stack-name ${FUNC,,} --capabilities CAPABILITY_IAM

    # Clean up...
    delparams
    rm packaged.yaml
}

case "$1" in
    "deps")
        deps
        ;;
    "clean")
        clean
        ;;
    "test")
        test
        ;;
    "build")
        build
        ;;
    "getparams")
        getparams
        ;;
    "delparams")
        delparams
        ;;
    "deploy")
        deploy
        ;;
    *)
        echo "The target {$1} want to execute doesn't exist"
        echo 
        echo "Usage"
        echo "./build deps      : go get and update all the dependencies"
        echo "./build clean     : removes the ./bin folder"
        echo "./build test      : uses SAM local and the event in event.json to test "
        echo "                    the implementation"
        echo "./build build     : creates the executable"
        echo "./build getparams : updates the SAM template with the values from the AWS"
        echo "                    Systems Manager Parameter Store"
        echo "./build delparams : removes the values of the environment variables in the "
        echo "                    SAM template"
        echo "./build deploy    : deploy the function to AWS Lambda"
        exit 2
        ;; 
esac

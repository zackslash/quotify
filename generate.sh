#!/bin/bash
# File: generate.sh
# Description: create packages for lambda deployment

echo -e "\033[4mEnsuring go-bindata...\033[0m"
(cd vendor/github.com/zackslash/go-bindata && go install .)

echo -e "\033[4mCollecting resources...\033[0m"
go-bindata -pkg quotify -o resources.go resources/

echo -e "\033[4mCompile for target platform (linux)\033[0m"
(cd delivery && GOOS=linux go build -o delivery)
(cd generation && GOOS=linux go build -o generation)

echo -e "\033[4mPackage for upload\033[0m"
(cd delivery && zip deployment.zip ./delivery)
(cd generation && zip deployment.zip ./generation)

#!/bin/bash 

#build the project for linux
env GOOS=linux GOARCH=amd64 go build -o books-api main.go

#build the container
docker build -t askewseth/books-api .


#run the container if the run flag is given
if [[ "$1" == "run" ]]; then
	docker run -p 5555:5555 --name books-api -d askewseth/books-api
fi


# Using the latest golang version 
FROM golang:latest

# Set the application directory 
WORKDIR /app

# Set the GOPATH 
ENV GOPATH $HOME/gopath

# Installation of libraries 
RUN go get github.com/Syfaro/telegram-bot-api

# Copy the code from current folder 
# to /app inside the container 
COPY . /app 

# Make port 80 available for links 
# and/or publish 
EXPOSE 80 

# Command to run the service
CMD go run main.go

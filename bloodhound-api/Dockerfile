# syntax=docker/dockerfile:1

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.17-alpine

# create a working directory inside the image
WORKDIR /app

ADD . .

# copy Go modules and dependencies to image
COPY go.mod ./
COPY go.sum ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY *.go ./

# copy the .env file
COPY .env ./

# compile application
RUN go build -o /godocker
 
# tells Docker that the container listens on specified network ports at runtime
EXPOSE 3000

# command to be used to execute when the image is used to start a container
CMD [ "/godocker" ]
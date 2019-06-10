FROM golang:1.12.4
 
RUN mkdir -p /app
 
WORKDIR /app
 
ADD . /app
 
RUN go build ./main.go
 
CMD ["./app"]

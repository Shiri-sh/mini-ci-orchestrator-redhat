FROM golang:1.26-alpine AS builder

#defining the working directory inside the container
WORKDIR /app

#copying the go.mod and go.sum and downloading the dependencies
COPY go.mod go.sum ./
RUN go mod download

#copying the source code into the container
COPY . .

#building the application for linux with CGO disabled to create a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

#final image
FROM alpine:latest

WORKDIR /root/

# copying the binary from the builder stage to the final image
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
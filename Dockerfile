# Build stage
FROM golang:1.24-alpine AS builder
 
WORKDIR /go/src/github.com/rohanchauhan02/internal-transfer
ENV GO111MODULE=on
ENV GODEBUG="madvdontneed=1"
 
# We want to populate the module cache based on the go.mod and go.sum files.
COPY go.mod .
COPY go.sum .
 
RUN go mod download
 
# Build stage
FROM builder AS server_builder
 
WORKDIR /go/src/github.com/rohanchauhan02/internal-transfer
 
COPY . .
 
# Build the application
RUN go build -o engine app/main.go
 
# Final stage
FROM gcr.io/distroless/base:nonroot
 
WORKDIR /internal-transfer/app
EXPOSE 11001
 
# Copy the built binary and configuration files
COPY --from=server_builder /go/src/github.com/rohanchauhan02/internal-transfer/engine .
COPY --from=server_builder /go/src/github.com/rohanchauhan02/internal-transfer/configs ./configs/
 
CMD ["/internal-transfer/app/engine"]

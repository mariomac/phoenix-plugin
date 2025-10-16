FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /code-reviewer main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates git

WORKDIR /

COPY --from=builder /code-reviewer /code-reviewer

ENTRYPOINT ["/code-reviewer"]

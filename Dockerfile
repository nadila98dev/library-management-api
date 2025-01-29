# version: '3.8'

# services:
#   goapp:
#   build: ,
#   container_name: goapp
#   ports:
#     - "8080:8080"
#     volumes:
#       - .:/app
  

  # Build Stage
  FROM golang:alpine as builder

  # Set the working directory inside the container
  WORKDIR /build
  
  # Copy Go modules and dependencies
  COPY go.mod go.sum ./
  RUN go mod download
  
  # Copy the application source code
  COPY . .
  
  # Build the application statically to run on scratch
  RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
  
  # Final Stage
  FROM scratch
  
  # Set working directory in the scratch image
  WORKDIR /app
  
  # Copy the compiled binary from the builder stage
  COPY --from=builder /build/main .
  
  # Expose the application port
  EXPOSE 8080
  
  # Command to run the application
  CMD ["./main"]
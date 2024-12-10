# Step 1: Use the official Golang image as the base image for building the app
FROM golang:1.23.4-alpine AS builder

# Step 2: Set the Current Working Directory inside the container
WORKDIR /app

# Step 3: Copy the Go modules manifests (go.mod and go.sum) into the container
COPY go.mod go.sum ./

# Step 4: Download all the dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Step 5: Copy the rest of the application's source code into the container
COPY . .

# Step 6: Build the Go app (replace 'main.go' with your app's entry point)
# We explicitly specify the output binary file to avoid the error
RUN go build -o myapp .

# Step 7: Create a smaller image for running the app
FROM alpine:latest

# Step 8: Install dependencies required to run the app (if needed)
RUN apk --no-cache add ca-certificates

# Step 9: Copy the pre-built binary from the 'builder' image
COPY --from=builder /app/myapp /usr/local/bin/myapp

# Step 10: Expose the port the app will run on (optional)
EXPOSE 8080

# Step 11: Command to run the binary
CMD ["myapp"]

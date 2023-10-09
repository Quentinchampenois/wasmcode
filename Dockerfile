# Use an official Go runtime as a parent image
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the current directory contents into the container
COPY . .

# Fetch dependencies
RUN go get -u github.com/skip2/go-qrcode

# Build the Go WebAssembly binary
RUN GOOS=js GOARCH=wasm go build -o main.wasm main.go

FROM nginx:latest

# Remove default Nginx files
RUN rm -rf /usr/share/nginx/html/*

# Copy files from builder stage and additional files
COPY --from=builder /go/src/app/main.wasm /usr/share/nginx/html/
COPY --from=builder /go/src/app/index.html /usr/share/nginx/html/
COPY --from=builder /usr/local/go/misc/wasm/wasm_exec.js /usr/share/nginx/html/

# Copy Nginx configuration file
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 8080
EXPOSE 8080

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]

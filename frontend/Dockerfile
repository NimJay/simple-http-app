FROM golang:1.19.2 as builder
WORKDIR /app
# Define the Go module.
RUN go mod init frontend
COPY *.go ./
# Build the Go module.
RUN CGO_ENABLED=0 GOOS=linux go build -o /frontend

# We use build two images — one for compiling the code, and another image for running the code.
# Learn more about multi-stage builds at https://docs.docker.com/build/building/multi-stage/
FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=builder /frontend /frontend
ENV PORT 8080
USER nonroot:nonroot
CMD ["/frontend"]

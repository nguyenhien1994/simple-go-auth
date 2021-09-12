# build stage
FROM golang as builder
# Add dependencies
WORKDIR /go/src/app
ADD . /go/src/app
# Build app
# RUN go mod download
RUN go build -o /go/bin/app ./main.go

# final stage
FROM alpine:latest
ARG PORT=8080

COPY --from=builder /go/bin/app /
EXPOSE $PORT
CMD  ["/app"]

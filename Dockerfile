FROM golang:1.15
WORKDIR /go/src/github.com/uhlig-it/shorts
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shorts .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /go/src/github.com/uhlig-it/shorts/shorts .
COPY deployment/files/shorts.yml .
CMD ["./shorts", "--urls", "shorts.yml"]

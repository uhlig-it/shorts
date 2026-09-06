FROM golang:1.27 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o shorts .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /src/shorts .
COPY deployment/files/shorts.yml .
CMD ["./shorts", "--urls", "shorts.yml"]

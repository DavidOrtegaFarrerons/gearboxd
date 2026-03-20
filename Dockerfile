FROM golang:1.26-alpine AS go
RUN apk add --no-cache tzdata
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/api

FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=Europe/Madrid

WORKDIR /app
COPY --from=go /app/app .
#COPY --from=go /app/migrations ./migrations

EXPOSE 4000
CMD ["./app"]
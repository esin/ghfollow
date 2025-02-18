FROM golang:1.24-alpine

WORKDIR /src

COPY githubfollower.go .
COPY go.mod .

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o githubfollower

FROM alpine:3.21

COPY --from=0 /src/githubfollower /app/githubfollower

ENTRYPOINT [ "/app/githubfollower" ]

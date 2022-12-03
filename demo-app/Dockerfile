FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /demo-app

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /demo-app /demo-app

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/demo-app"]

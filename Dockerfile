FROM golang:1.23 AS development
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
# want: hot-reload dev environment
CMD ["sh", "-c", "while :; do sleep 2024; done"]


FROM golang:1.23 AS build
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12 AS production
COPY --from=builder /go/bin/app /app
CMD [ "/app" ]

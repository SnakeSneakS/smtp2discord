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
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12 AS production
COPY --from=build /go/bin/app /go/bin/app
CMD [ "/go/bin/app" ]

#docker build --target production -t smtp2discord .
#docker run --rm smtp2discord


FROM golang:1.24 AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download
# RUN go vet -v
# RUN go test -v

RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o /go/bin/app

FROM scratch
WORKDIR /app

COPY --from=build /go/bin/app ./bin
COPY devices/ ./devices/

ENV STORAGE_PATH="/data/"
EXPOSE 51826/tcp
EXPOSE 5353/udp

CMD ["/app/bin"]
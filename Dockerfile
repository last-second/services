FROM golang:1.17

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY ./pkg ./pkg
COPY ./cmd ./cmd
COPY main.go ./
RUN go build -o app ./main.go

ENTRYPOINT ["./app"]
CMD ["serve"]

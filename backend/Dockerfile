FROM golang:1.19.5-bullseye

WORKDIR /backend
COPY . .
RUN go build .
CMD ["./walkietalkie"]
FROM --platform=amd64 golang:1.20-alpine as base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/* cmd/
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gotube



FROM --platform=arm64 golang:1.20-alpine

WORKDIR /app
COPY --from=base gotube .
CMD [ "./gotube" ]
FROM golang:1.25.7-alpine3.23 AS build

WORKDIR /app

COPY cmd/internal-tool/go.mod cmd/internal-tool/go.sum ./
RUN go mod download

COPY cmd/internal-tool/ .
RUN CGO_ENABLED=0 go build -o internal-tool .

FROM alpine:3.23.0
RUN apk add --no-cache ca-certificates

COPY --from=build /app/internal-tool /app/internal-tool

ENTRYPOINT ["/app/internal-tool"]

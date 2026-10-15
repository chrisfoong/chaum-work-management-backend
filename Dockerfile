# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/server ./cmd/server

FROM alpine:3.23 AS runtime
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 chaum
WORKDIR /app
COPY --from=build --chown=chaum:chaum /out/server ./server
USER chaum
ENV PORT=8080 GIN_MODE=release TZ=Asia/Bangkok
EXPOSE 8080
CMD ["./server"]

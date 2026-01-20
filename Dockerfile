FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o go-p2p

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/go-p2p .
EXPOSE 8081
# Using wget for healthcheck because it's standard in Alpine
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://localhost:8081/health || exit 1
CMD ["./go-p2p", "-ci"]
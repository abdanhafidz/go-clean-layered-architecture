# =====================
# BUILD STAGE
# =====================
FROM golang:1.25.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# =====================
# RUNTIME STAGE
# =====================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/app /app/app

# Non-root user
USER nonroot:nonroot

EXPOSE 8080

CMD ["/app/app"]

FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/main .

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["/app/main"]

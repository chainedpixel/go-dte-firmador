FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o signserver ./cmd/signserver

FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/signserver .

# Copy config files
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/configs/locales /app/configs/locales

# Make uploads directory and set permissions
RUN mkdir -p /app/uploads
COPY ./uploads/ /app/uploads/
RUN rm -f /app/uploads/*.md
RUN chmod -R 755 /app/uploads

EXPOSE 8113

ENTRYPOINT ["/app/signserver"]
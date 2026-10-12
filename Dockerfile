FROM golang:1.26.1-bookworm AS builder

WORKDIR /src

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN echo "=== BUILDER TEMPLATES ===" && \
    find /src/templates -type f && \
    echo "=== BUILDER STATIC ===" && \
    find /src/static -type f

RUN go build -tags fts5 -o /out/forum .

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /app/data

COPY --from=builder /out/forum ./forum
COPY --from=builder /src/templates ./templates/
COPY --from=builder /src/static ./static/

EXPOSE 8080

CMD ["./forum"]
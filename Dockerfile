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

RUN mkdir -p /out && go build -tags fts5 -o /out/forum .

FROM debian:bookworm-slim

WORKDIR /app

RUN mkdir -p /app/data

COPY --from=builder /out/forum ./forum
COPY --from=builder /src/templates ./templates/
COPY --from=builder /src/static ./static/

EXPOSE 8080

CMD ["./forum"]
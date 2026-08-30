FROM node:24-bookworm-slim AS frontend

RUN corepack enable
WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

FROM golang:1.26-bookworm AS backend

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=frontend /app/frontend/out/. ./internal/frontend/dist/
RUN CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /nilchan ./cmd

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /data
WORKDIR /app
COPY --from=backend /nilchan ./nilchan
COPY internal/config/config.yml ./internal/config/config.yml
ENV STORAGE_PATH=/data/storage.db
EXPOSE 8080
CMD ["./nilchan"]

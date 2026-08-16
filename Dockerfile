FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 \
    go build -trimpath -ldflags="-s -w" -o /out/batch-rename ./cmd/batch-rename

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
COPY --from=build /out/batch-rename /usr/local/bin/batch-rename

ENTRYPOINT ["batch-rename"]

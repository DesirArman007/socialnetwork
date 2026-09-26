# syntax=docker/dockerfile:1

FROM golang:1.25.4-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/social-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/social-api /social-api

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/social-api"]

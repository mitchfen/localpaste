# Build Stage
FROM golang:1.26-alpine AS build
WORKDIR /source

COPY go.mod ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o /app/localpaste ./cmd/localpaste

# Runtime Stage
FROM scratch
WORKDIR /app

COPY --from=build /app/localpaste .
COPY --from=build /source/assets/ assets/

EXPOSE 8080

ENTRYPOINT ["./localpaste"]

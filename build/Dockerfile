# Build Stage
FROM golang:1.26-alpine AS build
WORKDIR /source

COPY src/go.mod ./
RUN go mod download

COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/localpaste .

# Runtime Stage
FROM scratch
WORKDIR /app

COPY --from=build /app/localpaste .
COPY --from=build /source/assets/ assets/

EXPOSE 8080

ENTRYPOINT ["./localpaste"]

FROM golang:latest AS builder

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY ./restServer/go.* ./

RUN go mod download

COPY ./restServer/ ./
RUN swag init --parseDependency --parseInternal --parseDepth 1

RUN CGO_ENABLED=0 GOOS=linux go build -o /build

FROM alpine:latest AS runner

COPY --from=builder build /bin/build

EXPOSE 8080

ENTRYPOINT [ "/bin/build" ]
FROM golang:1.15-alpine AS build
COPY . /app
ENV CGO_ENABLED=0
WORKDIR /app
RUN go build -o bank-service-api ./cmd/api

FROM alpine:3
COPY --from=build /app/bank-service-api /app/bank-service-api
CMD ["/app/bank-service-api"]
EXPOSE 9999

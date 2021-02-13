FROM alpine:3

COPY bank-service-api /app/

CMD ["/app/bank-service-api"]

EXPOSE 9999

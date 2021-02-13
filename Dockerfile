FROM alpine:3

COPY bank /app/

CMD ["/app/bank"]

EXPOSE 9999

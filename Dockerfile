FROM alpine:3.5

COPY books-api .

EXPOSE 5555

CMD ["./books-api"]

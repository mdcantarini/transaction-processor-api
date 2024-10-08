FROM golang:1.22.8-alpine

WORKDIR /app
COPY . .
RUN go build -o api .

EXPOSE 8000

CMD ["./api"]
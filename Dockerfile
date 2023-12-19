 #Use the official Golang image as a builder stage
FROM golang:latest as builder
LABEL maintainer="Iqra Shams <iqra.shams339@gmail.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
FROM postgres:latest as postgres
ENV POSTGRES_PASSWORD=${DB_PASSWORD} \
    POSTGRES_USER=${DB_USER} \
    POSTGRES_DB=${DB_NAME}
FROM builder
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 1304
CMD ["./main"]



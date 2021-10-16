# Build stage
from golang:1.17-alpine as builder
workdir /app
copy . .
run go build -o banco main.go

# run stage
from alpine:3.14
workdir /app
copy --from=builder /app/banco .
copy app.env /app
copy db/migrations /app/db/migrations
copy wait-for.sh /app/wait-for.sh
expose 8080
cmd ["/app/banco"]
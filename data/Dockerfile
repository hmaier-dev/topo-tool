FROM golang:1.20
LABEL authors="hmaier@ipb-halle.de"
WORKDIR /app
COPY . .

RUN go mod download github.com/go-sql-driver/mysql
RUN go build -o scraper.go
CMD ["./scraper.go"]

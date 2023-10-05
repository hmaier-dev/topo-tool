FROM golang:1.20
LABEL authors="hmaier@ipb-halle.de"
WORKDIR /app
COPY . .

RUN go mod download github.com/go-sql-driver/mysql
RUN go build -o main.go
CMD ["go", "run", "main.go"]

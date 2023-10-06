FROM golang:1.20
LABEL authors="hmaier@ipb-halle.de"
WORKDIR /app
COPY . .
# RUN wget https://unpkg.com/htmx.org/dist/htmx.min.js > ./js/htmx.min.js
RUN mkdir -p /app/js
ADD https://unpkg.com/htmx.org/dist/htmx.min.js /app/js/htmx.min.js

RUN go mod download github.com/go-sql-driver/mysql
RUN go build -o server.go

EXPOSE 8181
#CMD ["go", "run", "server.go"]
CMD ["./server.go"]

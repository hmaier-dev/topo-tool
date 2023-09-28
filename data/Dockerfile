FROM golang:1.20
LABEL authors="hmaier@ipb-halle.de"
WORKDIR /app
COPY . .
# RUN wget https://unpkg.com/htmx.org/dist/htmx.min.js > ./js/htmx.min.js

CMD ["go", "run", "clearpass.go"]

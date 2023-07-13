FROM python:3.11-alpine
WORKDIR /app
#gcc musl-dev linux-headers
RUN apk add --no-cache mariadb-dev build-base
COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt
COPY . .

CMD ["python","./main.py"]
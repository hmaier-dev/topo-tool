FROM python:3.11
LABEL authors="hmaier"

WORKDIR /app
COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt
COPY . .

ENTRYPOINT ["python", "-u", "main.py"]
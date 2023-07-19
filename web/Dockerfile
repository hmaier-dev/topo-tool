FROM python:3.11
LABEL authors="hmaier@ipb-halle.de"
WORKDIR /app

# Turns off buffering for easier container logging
ENV PYTHONUNBUFFERED=1
ENV FLASK_APP=main.py

COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt
COPY . .

EXPOSE 8080
#CMD ["stdbuf", "-oL", "python", "main.py"]
#CMD ["python", "-m","flask","run"]
CMD ["python", "main.py"]

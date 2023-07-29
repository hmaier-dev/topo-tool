FROM python:3.11
WORKDIR /app

# Turns off buffering for easier container logging
ENV PYTHONUNBUFFERED=1

COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt
COPY . .

EXPOSE 8181

#CMD ["python","-u","main.py"]
# CMD ["stdbuf", "-oL", "python", "main.py", "--scanner"]
CMD ["stdbuf", "-oL", "python", "main.py", "--server"]

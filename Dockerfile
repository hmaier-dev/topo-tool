FROM python
WORKDIR .
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt
CMD ["python","main.py"]
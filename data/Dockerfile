FROM python
WORKDIR .
COPY requirements.txt requirements.txt
COPY cred.py cred.py
COPY api.py api.py
COPY discover.py discover.py
COPY filter.py filter.py
COPY db.py db.py
RUN pip install -r requirements.txt
CMD ["python","main.py"]
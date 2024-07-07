FROM python:3.12-slim

ENV PIP_DEFAULT_TIMEOUT=100 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    PIP_NO_CACHE_DIR=1

RUN apt-get update && \
    apt-get install --no-install-recommends -y build-essential git python3-dev python3-pip && \
    rm -rf /var/lib/apt/lists/*

RUN pip install mfrc522 requests RPi.GPIO

WORKDIR /app
COPY raspi.py ./
COPY rfid_reader.py ./

CMD ["python3", "rfid_reader.py"]

FROM python:3.9-slim

WORKDIR /app

ENV PIP_DEFAULT_TIMEOUT=100 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    PIP_NO_CACHE_DIR=1

RUN apt-get update && \
    apt-get install --no-install-recommends -y build-essential git python3-dev python3-pip python3-smbus && \
    rm -rf /var/lib/apt/lists/*

RUN pip install spidev mfrc522

COPY rfid_reader.py ./

CMD ["python3", "rfid_reader.py"]

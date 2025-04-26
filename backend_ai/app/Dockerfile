FROM python:3.12-slim

WORKDIR /app

RUN apt update -y && \
    apt install -y curl && \
    rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN apt-get update && apt-get install ffmpeg libsm6 libxext6  -y
RUN python -m pip install --upgrade pip
RUN pip install -r requirements.txt

COPY . .

ENV PYTHONPATH="/app"

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]
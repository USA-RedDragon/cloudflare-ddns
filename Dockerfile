FROM python:3.7-alpine

WORKDIR /app

RUN pip install requests

COPY cloudflare-ddns.py .

CMD [ "python", "/app/cloudflare-ddns.py" ]

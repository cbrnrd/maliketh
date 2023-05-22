FROM nginx:alpine

WORKDIR /app

ENV PROXY_HOST=proxy.example.com
ENV C2_HOST=c2.example.com

COPY nginx.conf /etc/nginx/nginx.conf
COPY fullchain.pem  /app/fullchain.pem
COPY privkey.pem /app/privkey.pem

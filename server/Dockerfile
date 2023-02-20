FROM python:3.11
ADD . /app
WORKDIR /app
ENV FLASK_APP=maliketh.app:app
RUN pip install -r requirements.txt
CMD flask run
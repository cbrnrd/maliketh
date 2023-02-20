FROM python:3.11
ADD . /app
WORKDIR /app
ENV FLASK_APP=maliketh.app:app
ENV FLASK_ENV=development
RUN pip install -r requirements.txt
CMD flask run -h 0.0.0.0

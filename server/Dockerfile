FROM python:3.11-slim as base

ENV PYTHONUNBUFFERED=TRUE
RUN pip install poetry

COPY ./pyproject.toml ./poetry.lock /app/

WORKDIR /app
RUN apt update && \
    apt install build-essential g++-mingw-w64 -y && \
    poetry config virtualenvs.create false && \
    poetry install --no-dev && \
    apt purge --autoremove build-essential -y

COPY . /app

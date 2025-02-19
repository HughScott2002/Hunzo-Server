FROM python:3.11-slim

WORKDIR /app
COPY . .

# Install Poetry
RUN pip install poetry
RUN poetry config virtualenvs.create false && poetry install --no-dev

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
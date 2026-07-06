FROM python:3.13-slim

WORKDIR /app

COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv

COPY pyproject.toml uv.lock ./
RUN uv sync --frozen

COPY . .

EXPOSE 8000

CMD ["sh", "-c", "[ -f db.sqlite3 ] || cat seed_data/db.sqlite3.gz.part-* | gunzip -c > db.sqlite3; uv run python manage.py migrate && uv run python manage.py runserver 0.0.0.0:8000"]

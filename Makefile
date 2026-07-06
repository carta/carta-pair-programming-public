.PHONY: run stop run-backend restore-db test test-v typecheck generate-migrations migrate

run:
	docker compose up --build

stop:
	docker compose down

# Reassemble + decompress the committed, split, gzipped seed database.
# Idempotent: does nothing if db.sqlite3 already exists.
restore-db:
	cd backend && [ -f db.sqlite3 ] || cat seed_data/db.sqlite3.gz.part-* | gunzip -c > db.sqlite3

run-backend: restore-db
	cd backend && uv sync && uv run python manage.py migrate && uv run python manage.py runserver 8000

test:
	cd backend && uv sync && uv run python -m pytest

test-v:
	cd backend && uv run python -m pytest -v

typecheck:
	cd backend && uv run python -m mypy .

generate-migrations:
	cd backend && uv run python manage.py makemigrations

migrate:
	cd backend && uv run python manage.py migrate

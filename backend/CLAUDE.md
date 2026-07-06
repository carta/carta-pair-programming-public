# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Banking/transaction management app built for pair programming interviews. Django 5.2 backend with Django Ninja REST API, SQLite database, paired with a Vite React frontend.

## Commands

All commands run from the repo root via `make` (they `cd backend` internally):

```bash
make run-backend           # uv sync + migrate + runserver on :8000
make test                  # uv sync + pytest
make test-v                # pytest verbose
make test-cov              # pytest with coverage
make typecheck             # mypy (strict mode)
make generate-migrations   # makemigrations
make migrate               # apply migrations
```

Run a single test from the backend directory:
```bash
cd backend && uv run pytest -k "test_name"
```

The database ships **pre-seeded**. It is committed as a gzipped, split artifact under
`backend/seed_data/db.sqlite3.gz.part-*` and reassembled into `backend/db.sqlite3` at setup
(the Docker entrypoint and `make restore-db` do this automatically; `db.sqlite3` itself is
gitignored). To restore it manually:
```bash
make restore-db
# equivalently: cd backend && cat seed_data/db.sqlite3.gz.part-* | gunzip -c > db.sqlite3
```
There is intentionally **no data-generation command in this repo** — the seed generator is
interviewer-only and lives in `carta/interviews` at `onsite_coding/seed_data/`.

## Architecture

**Django project**: `bank/` — settings, root URL config, WSGI/ASGI.

**Single app**: `accounts/` — all domain logic.

**API layer** (`accounts/api.py`): Django Ninja Router with three endpoints, all under `/api/`:
- `GET /users` — list all users
- `GET /users/{user_id}/accounts` — user's accounts with computed balance
- `GET /accounts/{account_id}/transactions?search=` — transactions with optional text search on memo/merchant/category

The router is mounted in `bank/urls.py` via `NinjaAPI().add_router("/", accounts_router)`.

**Models** (`accounts/models.py`):
- `Account` → FK to Django `User`, types: checking/savings/credit
- `Transaction` → FK to `Account`, types: debit/credit, optional FK to `Category` and `Merchant`
- `Category` and `Merchant` — simple name-only lookup tables

**Schemas**: Ninja `Schema` classes in `api.py` define typed API responses (`UserOut`, `AccountOut`, `TransactionOut`, etc.).

## Testing

Pytest with `pytest-django`. Tests live in `accounts/tests.py`.

Test factories (`accounts/factories.py`) use Factory Boy for all models. Tests use `@pytest.mark.django_db` decorator and Django's test `Client`.

## Type Checking

Mypy runs in strict mode with `mypy_django_plugin`. Test/factory modules have relaxed type checking rules (see `pyproject.toml` overrides).

## Dependencies

Managed with `uv`. Run `uv sync` to install. Lock file is `uv.lock`.

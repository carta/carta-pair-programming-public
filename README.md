# Pair Programming Interview

A banking application used for Carta's pair programming interviews. Candidates work with an interviewer to extend and improve the app during the session.

## Overview

The app lets users view bank accounts (checking, savings, credit), see transaction histories, and search transactions by merchant, category, or memo.

**Tech stack:** Django + Django Ninja API backend (SQLite), React + TypeScript + Vite frontend, styled with Tailwind CSS and shadcn/ui.

### Data Models

- **Account** — belongs to a User. Has a name, type (checking/savings/credit), and optional notes. Balance is computed from transactions.
- **Transaction** — belongs to an Account. Tracks date, memo, amount, and type (debit/credit). Optionally linked to a Category and Merchant.
- **Category** — a unique name (e.g. "Groceries").
- **Merchant** — a unique name (e.g. "Whole Foods").

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/users` | List all users |
| GET | `/api/users/{id}/accounts` | List accounts for a user (with computed balances) |
| GET | `/api/accounts/{id}/transactions` | List transactions for an account (supports `?search=`) |

## Prerequisites
Please ensure prior to your interview that you are able to run the application via the instructions below.

- [Docker](https://www.docker.com/)

## Running

```bash
docker compose up --build
```

Visit [http://localhost:5173](http://localhost:5173).

> On first start, the backend restores its pre-seeded database from the compressed files in
> `backend/seed_data/` (a few seconds) before the API becomes available.

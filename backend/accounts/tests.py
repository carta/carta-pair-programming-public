from datetime import date
from decimal import Decimal

import pytest
from django.test import Client

from accounts.factories import (
    AccountFactory,
    CategoryFactory,
    MerchantFactory,
    TransactionFactory,
    UserFactory,
)


@pytest.mark.django_db
def test_list_users(client: Client) -> None:
    users = UserFactory.create_batch(3)

    resp = client.get("/api/users")

    assert resp.status_code == 200
    data = resp.json()
    assert len(data) == 3
    returned_usernames = {u["username"] for u in data}
    assert returned_usernames == {u.username for u in users}


@pytest.mark.django_db
def test_list_users_includes_transaction_count(client: Client) -> None:
    busy_user = UserFactory()
    busy_account = AccountFactory(user=busy_user)
    TransactionFactory.create_batch(3, account=busy_account)

    quiet_user = UserFactory()
    AccountFactory(user=quiet_user)

    resp = client.get("/api/users")

    assert resp.status_code == 200
    counts = {u["username"]: u["transaction_count"] for u in resp.json()}
    assert counts[busy_user.username] == 3
    assert counts[quiet_user.username] == 0


@pytest.mark.django_db
def test_list_accounts_with_balance(client: Client) -> None:
    user = UserFactory()
    account = AccountFactory(user=user, name="Checking", account_type="checking")
    TransactionFactory(account=account, amount=Decimal("500.00"), type="credit")
    TransactionFactory(account=account, amount=Decimal("150.00"), type="debit")

    resp = client.get(f"/api/users/{user.id}/accounts")

    assert resp.status_code == 200
    data = resp.json()
    assert len(data) == 1
    assert data[0]["name"] == "Checking"
    assert Decimal(data[0]["balance"]) == Decimal("350.00")


@pytest.mark.django_db
def test_list_transactions_search(client: Client) -> None:
    account = AccountFactory()
    category = CategoryFactory(name="Groceries")
    merchant = MerchantFactory(name="Whole Foods")
    TransactionFactory(
        account=account,
        category=category,
        merchant=merchant,
        date=date(2025, 1, 15),
    )
    TransactionFactory(
        account=account,
        date=date(2025, 1, 16),
    )

    # search by merchant name
    resp = client.get(f"/api/accounts/{account.id}/transactions?search=whole")

    assert resp.status_code == 200
    data = resp.json()
    assert len(data) == 1
    assert data[0]["merchant"]["name"] == "Whole Foods"

from datetime import date
from decimal import Decimal

from django.contrib.auth.models import User
from django.db.models import Count
from django.http import HttpRequest
from django.shortcuts import get_object_or_404
from ninja import Query, Router, Schema

from .models import Account, Transaction

router = Router()


class UserOut(Schema):
    id: int
    username: str
    email: str
    transaction_count: int


class AccountOut(Schema):
    id: int
    name: str
    account_type: str
    notes: str | None
    balance: Decimal


class CategoryOut(Schema):
    id: int
    name: str


class MerchantOut(Schema):
    id: int
    name: str


class TransactionAccountOut(Schema):
    id: int
    name: str
    account_type: str


class TransactionOut(Schema):
    id: int
    date: date
    amount: Decimal
    category: CategoryOut | None
    merchant: MerchantOut | None
    type: str
    account: TransactionAccountOut


@router.get("/users", response=list[UserOut])
def list_users(request: HttpRequest) -> list[UserOut]:
    users = User.objects.annotate(transaction_count=Count("accounts__transactions"))
    return [
        UserOut(
            id=u.id,
            username=u.username,
            email=u.email,
            transaction_count=u.transaction_count,
        )
        for u in users
    ]


def _compute_balance(transactions: list[Transaction]) -> Decimal:
    return sum(
        (t.amount if t.type == Transaction.TransactionType.CREDIT else -t.amount for t in transactions),
        Decimal(0),
    )


@router.get("/users/{user_id}/accounts", response=list[AccountOut])
def list_accounts(request: HttpRequest, user_id: int) -> list[AccountOut]:
    user = get_object_or_404(User, id=user_id)
    accounts = user.accounts.all()
    return [
        AccountOut(
            id=a.id,
            name=a.name,
            account_type=a.account_type,
            notes=a.notes,
            balance=_compute_balance(list(a.transactions.all())),
        )
        for a in accounts
    ]


@router.get("/accounts/{account_id}/transactions", response=list[TransactionOut])
def list_transactions(
    request: HttpRequest, account_id: int, search: str = Query("")  # type: ignore[type-arg]
) -> list[TransactionOut]:
    account = get_object_or_404(Account, id=account_id)
    transactions = list(account.transactions.order_by("date").all())

    result = [
        TransactionOut(
            id=t.id,
            date=t.date,
            amount=t.amount,
            category=CategoryOut(id=t.category.id, name=t.category.name) if t.category else None,
            merchant=MerchantOut(id=t.merchant.id, name=t.merchant.name) if t.merchant else None,
            type=t.type,
            account=TransactionAccountOut(id=account.id, name=account.name, account_type=account.account_type),
        )
        for t in transactions
    ]

    if search:
        term = search.lower()
        return [
            t
            for t in result
            if (t.merchant and term in t.merchant.name.lower())
            or (t.category and term in t.category.name.lower())
        ]

    return result

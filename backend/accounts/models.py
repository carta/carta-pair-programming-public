from django.db import models
from django.contrib.auth.models import User


class Account(models.Model):
    class AccountType(models.TextChoices):
        CHECKING = "checking"
        SAVINGS = "savings"
        CREDIT = "credit"

    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name="accounts")
    name = models.CharField(max_length=255)
    account_type = models.CharField(max_length=20, choices=AccountType.choices)
    notes = models.TextField(null=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)

    def __str__(self) -> str:
        return f"{self.name} ({self.account_type})"


class Category(models.Model):
    name = models.CharField(max_length=100, unique=True)

    class Meta:
        verbose_name_plural = "categories"

    def __str__(self) -> str:
        return self.name


class Merchant(models.Model):
    name = models.CharField(max_length=255, unique=True)

    def __str__(self) -> str:
        return self.name


class Transaction(models.Model):
    class TransactionType(models.TextChoices):
        DEBIT = "debit"
        CREDIT = "credit"

    account = models.ForeignKey(Account, on_delete=models.CASCADE, related_name="transactions")
    date = models.DateField()
    amount = models.DecimalField(max_digits=12, decimal_places=2)
    category = models.ForeignKey(
        Category, on_delete=models.SET_NULL, null=True, blank=True, related_name="transactions"
    )
    merchant = models.ForeignKey(
        Merchant, on_delete=models.SET_NULL, null=True, blank=True, related_name="transactions"
    )
    type = models.CharField(max_length=10, choices=TransactionType.choices)
    created_at = models.DateTimeField(auto_now_add=True)

    def __str__(self) -> str:
        return f"{self.date} {self.merchant} {self.amount}"

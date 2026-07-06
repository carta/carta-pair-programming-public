import factory
from django.contrib.auth.models import User

from .models import Account, Category, Merchant, Transaction


class UserFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = User

    username = factory.LazyAttribute(lambda obj: f"{obj.first_name}.{obj.last_name}".lower())
    email = factory.LazyAttribute(lambda obj: f"{obj.username}@example.com")
    first_name = factory.Faker("first_name")
    last_name = factory.Faker("last_name")
    password = factory.django.Password("testpass123")


class AccountFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Account

    user = factory.SubFactory(UserFactory)
    name = factory.Iterator(["Checking", "Savings", "Credit Card"])
    account_type = factory.Iterator(Account.AccountType.values)


class CategoryFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Category

    name = factory.Sequence(lambda n: f"Category {n}")


class MerchantFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Merchant

    name = factory.Sequence(lambda n: f"Merchant {n}")


class TransactionFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Transaction

    account = factory.SubFactory(AccountFactory)
    date = factory.Faker("date_object")
    amount = factory.Faker("pydecimal", left_digits=4, right_digits=2, positive=True)
    category = factory.SubFactory(CategoryFactory)
    merchant = factory.SubFactory(MerchantFactory)
    type = factory.Faker("random_element", elements=Transaction.TransactionType.values)

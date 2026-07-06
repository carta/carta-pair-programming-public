from django.contrib import admin
from django.urls import path
from ninja import NinjaAPI

from accounts.api import router as accounts_router

api = NinjaAPI()
api.add_router("/", accounts_router)

urlpatterns = [
    path("admin/", admin.site.urls),
    path("api/", api.urls),
]

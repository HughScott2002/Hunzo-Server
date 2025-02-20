# app/api/main.py
from fastapi import FastAPI # type: ignore
from .routes import router

app = FastAPI(title="FastAPI Project")
app.include_router(router)
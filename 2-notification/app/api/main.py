# app/api/main.py
from fastapi import FastAPI
from .routes import router

app = FastAPI(title="FastAPI Project")
app.include_router(router)
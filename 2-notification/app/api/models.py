# app/api/models.py
from pydantic import BaseModel

class Item(BaseModel):
    id: int
    name: str
    description: str | None = None
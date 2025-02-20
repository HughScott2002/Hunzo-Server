# app/api/routes.py
from fastapi import APIRouter # type: ignore
from .models import Item

router = APIRouter()


@router.get("/health")
async def health():
    return {
        "OKAY": "OKAY"
    }
    
@router.get("/api/notifications")
async def notifications():
    return {
        "OKAY": "OKAY"
    }

@router.get("/items/{item_id}")
async def read_item(item_id: int):
    return {"item_id": item_id}

@router.post("/items/")
async def create_item(item: Item):
    return item
# app/api/routes.py
from datetime import datetime
import uuid
from fastapi import APIRouter # type: ignore
from .models import Item, NotificationResponse

router = APIRouter(prefix="/api/notifications")


NOTIFICATIONS = [
    {
        "notification_id": str(uuid.uuid4()),
        "account_id": "23423423423443",
        "is_read": False,
        "was_dismissed": False,
        "label": "Tenner Stafford",
        "content": "You have sent $200.00 to Tenner Stafford",
        "date": datetime.utcnow(),
        "icon": "https://github.com/shadcn.png",
        "type": None,
        "priority": "normal",
        "category": "transaction",
        "action_url": None
    },
     {
        "notification_id": str(uuid.uuid4()),
        "account_id": "23423423423443",
        "is_read": False,
        "was_dismissed": False,
        "label": "Tenner Stafford",
        "content": "You have sent $200.00 to Tenner Stafford",
        "date": datetime.utcnow(),
        "icon": "https://github.com/shadcn.png",
        "type": "action",
        "priority": "normal",
        "category": "transaction",
        "action_url": None
    }
    
]

@router.get("/health")
async def health():
     return {
        "status": "healthy",
        "service": "notification-service",
        "kafka_consumer": "running"
    }
    
   
@router.get("", response_model=NotificationResponse)
async def get_notifications(page: int = 1, page_size: int = 10):
    start_idx = (page - 1) * page_size
    end_idx = start_idx + page_size
    paginated_notifications = NOTIFICATIONS[start_idx:end_idx]
    unread_count = sum(1 for n in NOTIFICATIONS if not n["is_read"])
    
    return {
        "notifications": paginated_notifications,
        "total": len(NOTIFICATIONS),
        "page": page,
        "page_size": page_size,
        "unread_count": unread_count
    }

@router.get("/{notification_id}")
async def get_notification(notification_id: str):
    notification = next(
        (n for n in NOTIFICATIONS if n["notification_id"] == notification_id),
        None
    )
    if notification:
        return notification
    return {"error": "Notification not found"}

@router.get("/items/{item_id}")
async def read_item(item_id: int):
    return {"item_id": item_id}

@router.post("/items/")
async def create_item(item: Item):
    return item



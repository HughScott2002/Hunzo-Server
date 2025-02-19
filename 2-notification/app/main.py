from fastapi import FastAPI
from app.core.config import settings
from app.consumer.kafka_consumer import start_kafka_consumer
import asyncio
from pydantic import BaseModel, Field
from typing import List, Optional
from datetime import datetime
import uuid

app = FastAPI()

class Notification(BaseModel):
    notification_id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    account_id: str
    is_read: bool = False
    was_dismissed: bool = False
    label: str
    content: str
    date: datetime = Field(default_factory=datetime.utcnow)
    type: Optional[str] = None
    icon: Optional[str] = None
    priority: Optional[str] = "normal"
    category: Optional[str] = None
    action_url: Optional[str] = None

class NotificationResponse(BaseModel):
    notifications: List[Notification]
    total: int
    page: int
    page_size: int
    unread_count: int

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
    }
]
# TODO: FIX THE KAFKA IMPLEMENTATION FOR NOTIFICATIONS

# @app.on_event("startup")
# async def startup_event():
#     loop = asyncio.get_event_loop()
#     loop.run_in_executor(None, start_kafka_consumer)

@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "service": "notification-service",
        "kafka_consumer": "running"
    }

@app.get("/health/detailed")
async def detailed_health_check():
    return {
        "status": "healthy",
        "components": {
            "api": "up",
            "kafka_consumer": "running",
            "version": settings.VERSION if hasattr(settings, 'VERSION') else "1.0.0"
        }
    }

@app.get("/notifications", response_model=NotificationResponse)
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

@app.get("/notifications/{account-id}/last-updated")
async def notication_lastupdated():
    return {
        "total": len(NOTIFICATIONS),
        "last-updated": datetime.utcnow(),
    }


@app.get("/notifications/{notification_id}")
async def get_notification(notification_id: str):
    notification = next(
        (n for n in NOTIFICATIONS if n["notification_id"] == notification_id),
        None
    )
    if notification:
        return notification
    return {"error": "Notification not found"}
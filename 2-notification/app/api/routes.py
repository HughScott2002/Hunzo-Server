# app/api/routes.py
from datetime import datetime, timedelta
from typing import Optional
import uuid
from fastapi import APIRouter, Query# type: ignore
from .models import Item, NotificationResponse
from fastapi.responses import JSONResponse

router = APIRouter(prefix="/api/notifications")

# Generate a list of 25 notifications to ensure we have multiple pages
NOTIFICATIONS = []

# Create notification types and templates
notification_templates = [
    {
        "label": "Tenner Stafford",
        "content": "You have sent $200.00 to Tenner Stafford",
        "type": None,
        "icon": "https://github.com/shadcn.png",
        "category": "transaction"
    },
    {
        "label": "Kafla Winser",
        "content": "You have received a payment request from Kafla Winser for $800.00",
        "type": "action",
        "icon": "https://github.com/shadcn.png",
        "category": "request"
    },
    {
        "label": "Security Alert",
        "content": "We've detected a login from a new device. Please verify if this was you.",
        "type": "action",
        "icon": None,
        "category": "security"
    },
    {
        "label": "System Notification",
        "content": "Your account verification is complete. You now have full access to all features.",
        "type": None,
        "icon": None,
        "category": "system"
    },
    {
        "label": "Promotion",
        "content": "Transfer funds this week and get a $10 bonus on your next transaction!",
        "type": None,
        "icon": None,
        "category": "promotion"
    }
]

# Generate 25 notifications with different timestamps and some read/unread
for i in range(25):
    template = notification_templates[i % len(notification_templates)]
    # Create notifications with different dates (newer to older)
    notification_date = datetime.utcnow() - timedelta(hours=i*3)
    
    # Make some notifications read (every third one)
    is_read = (i % 3 == 0)
    
    notification = {
        "notification_id": str(uuid.uuid4()),
        "account_id": "23423423423443",  # Same account for all in this example
        "is_read": is_read,
        "was_dismissed": False,
        "label": template["label"],
        "content": template["content"],
        "date": notification_date,
        "icon": template["icon"],
        "type": template["type"],
        "priority": "high" if "Security" in template["label"] else "normal",
        "category": template["category"],
        "action_url": None
    }
    NOTIFICATIONS.append(notification)

@router.get("/health")
async def health():
    return {
        "status": "healthy",
        "service": "notification-service",
        "kafka_consumer": "running"
    }

@router.options("")
async def options_notifications():
    return JSONResponse(
        content={},
        headers={
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "GET, POST, OPTIONS, PUT",
            "Access-Control-Allow-Headers": "Content-Type, Authorization",
        },
    )

@router.get("", response_model=NotificationResponse)
async def get_notifications(
    account_id: Optional[str] = Query(None),
    page: int = Query(1),
    page_size: int = Query(10)
):
    # If account ID is provided, filter by it
    # For simplicity, if no account ID is provided, return all notifications
    if account_id:
        filtered_notifications = [n for n in NOTIFICATIONS if n["account_id"] == account_id]
    else:
        filtered_notifications = NOTIFICATIONS
        
    # Sort notifications by date (newest first)
    filtered_notifications = sorted(
        filtered_notifications, 
        key=lambda x: x["date"] if isinstance(x["date"], datetime) else datetime.fromisoformat(str(x["date"])), 
        reverse=True
    )
    
    # Paginate the filtered notifications
    start_idx = (page - 1) * page_size
    end_idx = start_idx + page_size
    paginated_notifications = filtered_notifications[start_idx:end_idx]
    
    # Count unread notifications
    unread_count = sum(1 for n in filtered_notifications if not n["is_read"])
    
    # Create response
    response = {
        "notifications": paginated_notifications,
        "total": len(filtered_notifications),
        "page": page,
        "page_size": page_size,
        "unread_count": unread_count
    }
    
    return response

# @router.get("/{notification_id}")
# async def get_notification(notification_id: str):
#     notification = next(
#         (n for n in NOTIFICATIONS if n["notification_id"] == notification_id),
#         None
#     )
#     if notification:
#         return notification
#     return {"error": "Notification not found"}

# Add endpoints for marking notifications as read
@router.put("/{notification_id}/read")
async def mark_as_read(notification_id: str):
    notification = next(
        (n for n in NOTIFICATIONS if n["notification_id"] == notification_id),
        None
    )
    if notification:
        notification["is_read"] = True
        return {"success": True}
    return {"error": "Notification not found"}

@router.put("/read-all")
async def mark_all_as_read(account_id: str = Query(...)):
    read_count = 0
    for notification in NOTIFICATIONS:
        if notification["account_id"] == account_id and not notification["is_read"]:
            notification["is_read"] = True
            read_count += 1
    return {"success": True, "read_count": read_count}

@router.get("/all-for-testing")
async def just_show_all():
    # Sort notifications by date (newest first)
    sorted_notifications = sorted(
        NOTIFICATIONS, 
        key=lambda x: x["date"] if isinstance(x["date"], datetime) else datetime.fromisoformat(str(x["date"])), 
        reverse=True
    )
    
    # Count unread notifications
    unread_count = sum(1 for n in NOTIFICATIONS if not n["is_read"])
    
    return {
        "notifications": sorted_notifications,  # Return all notifications
        "total": len(NOTIFICATIONS),
        "page": 1,
        "page_size": len(NOTIFICATIONS),  # Set page size to total count
        "unread_count": unread_count
    }
    
@router.get("/items/{item_id}")
async def read_item(item_id: int):
    return {"item_id": item_id}

@router.post("/items/")
async def create_item(item: Item):
    return item
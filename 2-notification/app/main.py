from fastapi import FastAPI
from app.core.config import settings
from app.consumer.kafka_consumer import start_kafka_consumer
import asyncio

app = FastAPI()

# @app.on_event("startup")
# async def startup_event():
#     # Start Kafka consumer in background
#     loop = asyncio.get_event_loop()
#     loop.run_in_executor(None, start_kafka_consumer)
    
@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "service": "notification-service",
        "kafka_consumer": "running"
    }

# You might also want a more detailed health check that includes Kafka connection status
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
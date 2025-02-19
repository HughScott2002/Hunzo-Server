from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # Kafka settings
    kafka_bootstrap_servers: str = "broker:9092"  # Use container name from docker-compose
    kafka_topic: str = "notifications"
    kafka_group_id: str = "notification_service"
    
    # Email settings
    smtp_server: str = "mailhog"
    smtp_port: int = 1025
    email_from: str = "test@example.com"
    smtp_password: str = ""

    class Config:
        env_prefix = "TESTING_NOTIFICATION_SERVICE_"

settings = Settings()

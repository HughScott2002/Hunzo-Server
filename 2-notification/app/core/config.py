from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    kafka_bootstrap_servers: str = "localhost:9092"
    sendgrid_api_key: str
    email_from: str = "noreply@example.com"

    class Config:
        env_file = ".env"

settings = Settings()
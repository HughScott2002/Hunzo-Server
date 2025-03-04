# app/api/main.py
from fastapi import FastAPI # type: ignore
from fastapi.middleware.cors import CORSMiddleware #type: ignore
from .routes import router

app = FastAPI(title="Notification Service")

# app.add_middleware(
#     CORSMiddleware,
#     allow_origins=["http://localhost:3000/"],  # For development. Restrict this in production
#     allow_credentials=True,
#     allow_methods=["*"],
#     allow_headers=["*"],
# )
# Add CORS middleware to allow cross-origin requests
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # In production, specify your frontend domain
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.include_router(router)
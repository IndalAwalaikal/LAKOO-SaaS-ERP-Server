import os
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from routers import insight, predict

from dotenv import load_dotenv
load_dotenv()

app = FastAPI(
    title="Lakoo AI Service",
    description="Python Data Science Backend for Sales Analytics and Predictions",
    version="1.0.0"
)

# Allow CORS for main Go API or Web Frontend
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(insight.router)
app.include_router(predict.router)

@app.get("/health")
def health_check():
    return {"status": "ok", "service": "lakoo-ai-engine"}

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", 8000))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)

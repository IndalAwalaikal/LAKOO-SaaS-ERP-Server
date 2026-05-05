from fastapi import APIRouter
from pydantic import BaseModel
from typing import List, Dict, Any
from services.prediction_service import forecast_demand

router = APIRouter(
    prefix="/predict",
    tags=["Stock & Demand Prediction"]
)

class DailySale(BaseModel):
    day_index: int
    sold_qty: float

class PredictionRequest(BaseModel):
    product_id: str
    historical_data: List[DailySale]
    days_to_predict: int = 7

@router.post("/demand")
def predict_demand(data: PredictionRequest) -> Dict[str, Any]:
    """
    Uses a simple Linear Regression model via Scikit-Learn to forecast inventory demand 
    for the next N days based on historical daily sales. Handles logic in the Prediction Service.
    """
    historical_data = [d.dict() for d in data.historical_data]
    return forecast_demand(data.product_id, historical_data, data.days_to_predict)

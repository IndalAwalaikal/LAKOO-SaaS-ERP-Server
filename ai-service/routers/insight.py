from fastapi import APIRouter
from pydantic import BaseModel
from typing import List, Dict, Any
from services.analytics_service import process_sales_insights

router = APIRouter(
    prefix="/insight",
    tags=["Sales Insights"]
)

class SaleRecord(BaseModel):
    date: str
    product_id: str
    amount: float
    qty: float

class AnalyticsRequest(BaseModel):
    tenant_id: str
    sales: List[SaleRecord]

@router.post("/sales")
def analyze_sales(data: AnalyticsRequest) -> Dict[str, Any]:
    """
    Receives raw sales transaction records (passed from the Go API or Database)
    and uses pandas to generate meaningful revenue analytics via the Analytics Service.
    """
    sales_data = [s.dict() for s in data.sales]
    return process_sales_insights(data.tenant_id, sales_data)

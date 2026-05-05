import pandas as pd
import numpy as np
from sklearn.linear_model import LinearRegression
from typing import List, Dict, Any

def forecast_demand(product_id: str, historical_data: List[Dict[str, Any]], days_to_predict: int) -> Dict[str, Any]:
    if len(historical_data) < 3:
        return {
            "product_id": product_id,
            "predictions": [],
            "message": "Not enough historical data (minimum 3 data points required)."
        }

    df = pd.DataFrame(historical_data)
    X = df[['day_index']].values
    y = df['sold_qty'].values

    model = LinearRegression()
    model.fit(X, y)

    last_day = int(df['day_index'].max())
    future_X = np.array([[last_day + i] for i in range(1, days_to_predict + 1)])
    
    predictions = model.predict(future_X)
    predictions = [max(0.0, float(pred)) for pred in predictions]

    results = []
    for i, pred in enumerate(predictions):
        results.append({
            "day_index": last_day + i + 1,
            "projected_demand": round(pred)
        })

    return {
        "product_id": product_id,
        "model_used": "Linear Regression (Scikit-Learn)",
        "predictions": results
    }

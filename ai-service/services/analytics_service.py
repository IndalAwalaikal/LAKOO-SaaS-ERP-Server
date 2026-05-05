import pandas as pd
from typing import List, Dict, Any

def process_sales_insights(tenant_id: str, sales_data: List[Dict[str, Any]]) -> Dict[str, Any]:
    if not sales_data:
        return {"total_revenue": 0.0, "top_selling_product": None, "insights": []}

    df = pd.DataFrame(sales_data)
    
    df['date'] = pd.to_datetime(df['date'])
    df['amount'] = pd.to_numeric(df['amount'])
    df['qty'] = pd.to_numeric(df['qty'])

    total_rev = float(df['amount'].sum())
    total_items = float(df['qty'].sum())

    daily_rev = df.groupby(df['date'].dt.date)['amount'].sum().to_dict()
    daily_rev_formatted = {str(k): float(v) for k, v in daily_rev.items()}

    product_sales = df.groupby('product_id')['qty'].sum()
    top_product_id = str(product_sales.idxmax()) if not product_sales.empty else None

    insights = []
    if total_rev > 1000000:
        insights.append("Sangat Bagus! Penjualan mencapai rekor target harian tinggi.")
    else:
        insights.append("Penjualan terpantau stabil berskala kecil. Fokus tingkatkan promosi pada product tertentu.")

    if product_sales.max() > (total_items * 0.5):
        insights.append("Ada indikasi ketergantungan pendapatan tinggi pada satu produk utama.")

    return {
        "tenant_id": tenant_id,
        "total_revenue": total_rev,
        "total_items_sold": total_items,
        "top_product_id": top_product_id,
        "daily_trend": daily_rev_formatted,
        "insights": insights
    }

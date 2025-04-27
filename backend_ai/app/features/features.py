from feast import FeatureStore
import pandas as pd
from datetime import datetime
import numpy as np

brands = pd.read_csv("features/brands.csv")
store = FeatureStore(repo_path="features/transactions/")


def get_client_feature(client_username: str):
    try:
        entity_df = pd.DataFrame.from_dict(
            {
                "client_id": [client_username],
                "timestamp": [
                    datetime.now(),
                ],
            }
        )
        retrieval_job = store.get_historical_features(
            entity_df=entity_df,
            features=["features:" + brand for brand in brands["brands"]],
        )
        feature_data = retrieval_job.to_df().iloc[0]
        return np.array([feature_data[brand] for brand in brands["brands"]])
    except Exception:
        return np.zeros(50)

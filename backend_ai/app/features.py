from datetime import datetime
import pandas as pd
import feast
from feast import FeatureStore, Entity, FeatureView, ValueType
from feast.infra.offline_stores.file_source import FileSource
from feast.types import PrimitiveFeastType

# Initialize a local feature store
feast.init_repo(repo_path="my_repo", template="local")

# Define entities and features
client_entity = Entity(
    name="client_id",
    value_type=ValueType.STRING,
    description="Client identifier",
)

# Define a feature view for brand features
brand_features_source = FileSource(
    path="brand_features.parquet",
    event_timestamp_column="timestamp",
    created_timestamp_column="created_timestamp",
)

brand_features_view = FeatureView(
    name="brand_features",
    entities=["client_id"],
    ttl=None,  # No expiration
    features=[
        *[
            feast.Feature(
                name=f"brand_{brand_id}_feature_{i}", dtype=PrimitiveFeastType.BOOL
            )
            for brand_id in range(1, 6)  # Assuming 5 brands
            for i in range(1, 11)
        ]  # 10 features per brand (50 total)
    ],
    batch_source=brand_features_source,
)


# Create some sample data
def generate_sample_data(num_records=1000):
    data = {
        "client_id": [f"client_{i}" for i in range(num_records)],
        "timestamp": [datetime.now() for _ in range(num_records)],
        "created_timestamp": [datetime.now() for _ in range(num_records)],
    }

    # Add 50 binary features (10 features for 5 brands)
    for brand_id in range(1, 6):
        for feature_num in range(1, 11):
            col_name = f"brand_{brand_id}_feature_{feature_num}"
            data[col_name] = [
                bool(i % 2) for i in range(num_records)
            ]  # Alternating True/False

    return pd.DataFrame(data)


# Generate and save sample data
sample_data = generate_sample_data()
sample_data.to_parquet("brand_features.parquet")

# Apply the feature store configuration
store = FeatureStore(repo_path="my_feature_repo")
store.apply([client_entity, brand_features_view])

# Ingest the data into the feature store
store.materialize_incremental(end_date=datetime.now())

print("Feature store setup complete and data ingested!")


def get_client_feature(client_id: int):
    return np.zeros(50)

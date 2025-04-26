from datetime import datetime
import pandas as pd
import feast
from feast import FeatureStore, Entity, FeatureView, ValueType
from feast.infra.offline_stores.file_source import FileSource
from feast.types import PrimitiveFeastType

# Initialize a local feature store
feast.init_repo(repo_path="my_feature_repo", template="local")

# Define entities and features
client_entity = Entity(
    name="client_id",
    value_type=ValueType.STRING,
    description="Client identifier",
)

# Define a feature view for brand features
brand_features_source = FileSource(
    path="features.parquet",
    event_timestamp_column="timestamp",
)

brands = [
    "Alexander, Mercado and West",
    "Bailey-Johnston",
    "Barnett Group",
    "Barton, Murray and Huber",
    "Brown-Vega",
    "Bryant, Novak and Burns",
    "Carlson Group",
    "Coleman Inc",
    "Cummings, Silva and Galloway",
    "Decker, Carter and Dunn",
    "Duncan and Sons",
    "Duran-Thompson",
    "Evans-Bell",
    "Farley and Sons",
    "Farmer, Pace and Saunders",
    "Fitzgerald, Delgado and Williams",
    "Fox Group",
    "Garner-Martinez",
    "Green PLC",
    "Hall, Price and Mccarthy",
    "Hammond Inc",
    "Hayes-Ali",
    "Ho-Hall",
    "Holden and Sons",
    "Jimenez, Baker and Davenport",
    "Jones, Anderson and Mcintyre",
    "Long, Jones and Carter",
    "Mitchell-Smith",
    "Moore, Peters and Freeman",
    "Moore-Wagner",
    "Parker-Fischer",
    "Peters-Hanna",
    "Peters-Mclaughlin",
    "Roberts, Sparks and Rose",
    "Salinas-Glenn",
    "Sharp Ltd",
    "Simmons, Diaz and Alexander",
    "Small-Anthony",
    "Smith LLC",
    "Smith, Jackson and George",
    "Smith-Hunter",
    "Taylor PLC",
    "Thomas-Thompson",
    "Velasquez, Nicholson and Burton",
    "Vincent, Collins and Arnold",
    "Walton-Stewart",
    "Ward-Olsen",
    "Watts, Fisher and Shah",
    "West PLC",
    "Yoder-Bradley",
]


brand_features_view = FeatureView(
    name="features",
    source=brand_features_source,
    entities=["client_id"],
    ttl=None,
    features=[
        feast.Feature(name=brand_id, dtype=PrimitiveFeastType.BOOL)
        for brand_id in brands
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

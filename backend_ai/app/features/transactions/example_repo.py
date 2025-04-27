import feast
from feast import Entity, FeatureView, Project, ValueType
from feast.infra.offline_stores.file_source import FileSource
from feast.types import PrimitiveFeastType
from brands import brands

# Initialize a local feature store
project = Project(name="transactions")

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

brand_features_view = FeatureView(
    name="features",
    source=brand_features_source,
    entities=[client_entity],
    ttl=None,
    schema=[
        feast.Field(name=brand_id, dtype=PrimitiveFeastType.BOOL) for brand_id in brands
    ],
)

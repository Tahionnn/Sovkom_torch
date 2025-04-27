from datetime import datetime
import feast
from feast import FeatureStore, Entity, FeatureView, Project, ValueType
from feast.infra.offline_stores.file_source import FileSource
from feast.types import PrimitiveFeastType

# Initialize a local feature store
project = Project(name="bin_trans")

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
    created_timestamp_column="created_timestamp",
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
    entities=[client_entity],
    ttl=None,
    schema=[
        feast.Field(name=brand_id, dtype=PrimitiveFeastType.BOOL) for brand_id in brands
    ],
)

from typing import Annotated
from fastapi import BackgroundTasks, FastAPI, File
import numpy as np

from prometheus_client import Gauge

DRIFT_SCORE = Gauge("model_drift_score", "Drift score between current and reference data")

app = FastAPI()




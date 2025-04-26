import mlflow

import numpy as np
import torch


class MultiVAETorch:
    def __init__(self, model_uri: str = None, k=5):
        self.k = k

        if model_uri is None:
            model_uri = "models:/multivae-torch/1"

        if model_uri.startswith("models:/"):
            self.model = mlflow.pytorch.load_model(model_uri)
        else:
            self.model = torch.jit.load(model_uri)

    def __call__(self, feature):
        with torch.no_grad():
            preds = self.model(torch.FloatTensor(feature))
            preds = np.array(preds)
        return np.argpartition(preds, -self.k)[-self.k :].tolist()


mlflow.set_tracking_uri("http://localhost:5000")
model = MultiVAETorch(model_uri="models:/multivae-torch/1")
# model = MultiVAEONNX(model_uri="models:/multivae/1")

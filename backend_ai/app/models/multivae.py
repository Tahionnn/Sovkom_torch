import mlflow
import mlflow.onnx

import numpy as np
import torch
import os

import onnx
import onnxruntime as rt


class MultiVAEONNX:
    def __init__(self, model_uri: str = "models:/multivae/1", k: int = 5):
        self.k = k
        self.session = rt.InferenceSession(
            mlflow.artifacts.download_artifacts(model_uri),
            providers=rt.get_available_providers(),
        )
        self.input_name = self.session.get_inputs()[0].name

    def predict(self, feature):
        preds = self.session.run(None, {self.input_name: feature})[0]
        return np.argpartition(preds, -self.k)[-self.k :].tolist()


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
        return np.argpartition(preds, -self.k)[-self.k :].tolist()


mlflow.set_tracking_uri("http://localhost:5000")
model = MultiVAETorch(model_uri="models:/multivae-torch/1")
#model = MultiVAEONNX(model_uri="models:/multivae/1")

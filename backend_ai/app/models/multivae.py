import numpy as np
import torch
import os

import onnx
import onnxruntime as rt


class MultiVAEONNX:
    def __init__(self, path: str = None,  k: int = 5):
        if path is None:
            current_dir = os.path.dirname(os.path.abspath(__file__))
            path = os.path.join(current_dir, "models", "multivae.onnx")
    
        self.session  = rt.InferenceSession(path, providers=rt.get_available_providers())
        self.input_name = self.session.get_inputs()[0].name
        self.k = k

    def predict(self, feature):
        preds = self.session.run(None, {self.input_name: feature})[0]
        return np.argpartition(preds, -self.k)[-self.k :].tolist()


class MultiVAETorch:
    def __init__(self, path, k=5):
        self.k = k
        self.model = torch.jit.load(path)

    def __call__(self, feature):
        with torch.no_grad():
            preds = self.model(torch.FloatTensor(feature))
        return np.argpartition(preds, -self.k)[-self.k :].tolist()


#model = MultiVAETorch("models\\multivae.pt")
model = MultiVAEONNX()

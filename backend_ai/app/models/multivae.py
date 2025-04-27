import numpy as np
import onnx
import onnxruntime as rt


class MultiVAEONNX:
    def __init__(self, model_uri: str = "models:/multivae/1", k: int = 5):
        self.k = k
        self.session = rt.InferenceSession(
            model_uri,
            providers=rt.get_available_providers(),
        )
        self.input_name = self.session.get_inputs()[0].name

    def predict(self, feature):
        preds = self.session.run(None, {self.input_name: feature})[0]
        return np.argpartition(preds, -self.k)[-self.k :].tolist()


model = MultiVAEONNX(model_uri="models/multivae.onnx")
# model = MultiVAEONNX(model_uri="models:/multivae/1")

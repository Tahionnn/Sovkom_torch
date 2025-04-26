import numpy as np
import torch
#import onnx
#import onnxruntime


class MultiVAEONNX:
    def __init__(self, path):
        self.sess = rt.InferenceSession(path, providers=rt.get_available_providers())
        self.input_name = self.sess.get_inputs()[0].name

    def predict(self, feature):
        return self.sess.run(None, {self.input_name: feature})[0]

class MultiVAETorch:
    def __init__(self, path, k=5):
        self.k = k
        self.model = torch.jit.load(path)

    def __call__(self, feature):
        with torch.no_grad():
            preds = self.model(torch.FloatTensor(feature))
        return np.argpartition(preds, -self.k)[-self.k:].tolist()


model = MultiVAETorch("models\\multivae.pt")
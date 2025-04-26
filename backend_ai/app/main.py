from typing import Annotated
import cv2
from fastapi import FastAPI, File
import numpy as np

from app.oсr_utils import ocr
from app.features import get_client_feature
from models.multivae import model


app = FastAPI()


@app.get("/predict")
def recomend_cashbacks(client_id: int):
    global model
    if model is None:
        raise Exception()
    feature = get_client_feature(client_id)
    return model(feature)


@app.post("/ocr")
def ocr_image(
    image: Annotated[bytes, File()],
):
    nparr = np.fromstring(image, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    return ocr(img, workers=1)


@app.get("/train")
def retrain_model():
    pass

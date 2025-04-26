from fastapi import FastAPI

from features import get_client_feature
from models.multivae import model


app = FastAPI()

@app.get("/predict")
def recomend_cashbacks(client_id: int):
    global model
    if model is None:
        raise Exception()
    feature = get_client_feature(client_id)
    return model(feature)

@app.get("/train")
def retrain_model():
    pass
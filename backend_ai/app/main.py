from typing import Annotated
import cv2
from fastapi import BackgroundTasks, FastAPI, File
import numpy as np
import logging

from app.oсr_utils import send_to_llm, parse_json_garbage
from app.features import get_client_feature
from app.models.multivae import model

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - client_id=%(client_id)s - offer=%(offer_id)s - cashback=%(cashback).2f - model=%(model_version)s",
    handlers=[
        logging.FileHandler("/tmp/log/recommendations.log"),
        logging.StreamHandler(),
    ],
)
logger = logging.getLogger(__name__)

app = FastAPI()


async def ocr_background(img: np.ndarray):
    response = await send_to_llm(
        img,
        """{
  "shop": "Название магазина",
  "items": [
    {
      "name": "Название продукта",
      "price": "Цена за одну единица измерения",
      "count": "Кол-во товара в единицах измерения",
      "measurement": "Единица измерения (например шт. или кг)",
      "overall": "Итоговая цена товара"
    }
  ]
}""",
        "6c3t8lc5ske6tIraZDakla91bZZMZ6Hf",
    )
    if "choices" in response:
        json_result = parse_json_garbage(response["choices"][0]["message"]["content"])
        print(json_result)
    else:
        print("Error mistral")


@app.get("/predict")
def recomend_cashbacks(client_username: str):
    global model
    if model is None:
        raise Exception()
    feature = get_client_feature(client_username)
    recommendations = model(feature)

    for cashback in recommendations:
        logger.info(
            "Cashback recommendation",
            extra={
                "client_username": client_username,
                "cashback": cashback,
                "model_version": "1.0",
            },
        )

    return recommendations


@app.post("/ocr")
def ocr_image(image: Annotated[bytes, File()], background_tasks: BackgroundTasks):
    nparr = np.fromstring(image, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    img = opencv_resize_maxcap(img, max_cap=500)
    background_tasks.add_task(ocr_background, img)
    return


@app.get("/train")
def retrain_model():
    pass

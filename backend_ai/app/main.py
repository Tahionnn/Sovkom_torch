from typing import Annotated
import cv2
from fastapi import BackgroundTasks, FastAPI, File
import numpy as np

from app.oсr_utils import (
    ocr,
    opencv_resize_maxcap,
    filters,
    parse_json_garbage,
    send_to_llm,
)
from app.features import get_client_feature
from app.models.multivae import model


app = FastAPI()


async def ocr_background(img: np.ndarray):
    print("start background ocr")
    results = [f(img) for f in filters]
    results = [" ".join(ocr(img)) for x in results]
    results = "\n".join(results)
    response = await send_to_llm(
        results,
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
        "api_key",
    )
    if "choices" in response:
        json_result = parse_json_garbage(response["choices"][0]["message"]["content"])
        print(json_result)
    else:
        print("Error mistral")


@app.get("/predict")
def recomend_cashbacks(client_id: int):
    global model
    if model is None:
        raise Exception()
    feature = get_client_feature(client_id)
    return model(feature)


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

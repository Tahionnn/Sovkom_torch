import cv2
import numpy as np
import easyocr
from typing import Union, List, Optional, Tuple, Dict, Any
from skimage.filters import threshold_local
import json
import aiohttp


def opencv_resize(image: np.ndarray, ratio: float) -> np.ndarray:
    width = int(image.shape[1] * ratio)
    height = int(image.shape[0] * ratio)
    dim = (width, height)
    return cv2.resize(image, dim, interpolation=cv2.INTER_AREA)


def bw_scanner(image: np.ndarray) -> np.ndarray:
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    T = threshold_local(gray, block_size=21, offset=5, method="gaussian")
    return (gray > T).astype("uint8") * 255


def otsu_scanner(image: np.ndarray) -> Tuple[int, np.ndarray]:
    image = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    thresh_val, binary_image = cv2.threshold(
        image, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU
    )
    return thresh_val, binary_image


def ocr(
    image: Union[str, np.ndarray],
    workers: int = 2,
    gpu: bool = True,
    text_threshold: float = 0.8,
) -> List[str]:
    reader = easyocr.Reader(["ru"], gpu=gpu)

    if isinstance(image, str):
        image = cv2.imread(image)

    result = reader.readtext(
        image, detail=0, text_threshold=text_threshold, workers=workers
    )
    return result


async def send_to_llm(
    text: str, json_format: str, api_key: str, session: aiohttp.ClientSession
) -> Dict[str, Any]:
    prompt = f"""Ты помогаешь составлять json-файлы с содержимым чека из магазина. На вход тебе постпает сырой OCR-текст. Извлеки список покупок, их цены и итоговую цену из чека
    Входные данные: 
    {text}

    Проанализируй текст и верни только валидный JSON в следующем формате:
    {json_format}
    """

    async with session.post(
        "https://api.mistral.ai/v1/chat/completions",
        headers={"Authorization": f"Bearer {api_key}"},
        json={
            "model": "mistral-small",
            "messages": [{"role": "user", "content": prompt}],
        },
    ) as response:
        return await response.json()

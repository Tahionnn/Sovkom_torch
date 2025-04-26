import cv2
import numpy as np
import easyocr
from typing import Union, List, Optional, Tuple, Dict, Any
from skimage.filters import threshold_local
import aiohttp
import json


def opencv_resize_ratio(image: np.ndarray, ratio: float) -> np.ndarray:
    width = int(image.shape[1] * ratio)
    height = int(image.shape[0] * ratio)
    dim = (width, height)
    return cv2.resize(image, dim, interpolation=cv2.INTER_AREA)


def opencv_resize_maxcap(image: np.ndarray, max_cap: int = 300) -> np.ndarray:
    max_size = max(image.shape[1], image.shape[0])
    ratio = max_cap / max_size
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


filters = [lambda x: x, bw_scanner, otsu_scanner]
reader = easyocr.Reader(["ru"], gpu=False)


def ocr(
    image: np.ndarray,
    workers: int = 2,
    text_threshold: float = 0.8,
) -> List[str]:
    global reader
    return reader.readtext(
        image, detail=0, text_threshold=text_threshold, workers=workers
    )


async def send_to_llm(text: str, json_format: str, api_key: str) -> Dict[str, Any]:
    async with aiohttp.ClientSession() as session:
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


def parse_json_garbage(s: str) -> dict[str, Any]:
    s = s.replace("'", '"')
    s = s[next(idx for idx, c in enumerate(s) if c in "{[") :]
    try:
        return json.loads(s)
    except json.JSONDecodeError as e:
        try:
            return json.loads(s[: e.pos])
        except Exception as e:
            print(s)
            raise e

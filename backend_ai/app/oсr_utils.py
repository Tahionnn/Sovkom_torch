import cv2
import io
import numpy as np
from PIL import Image
import easyocr
from typing import Union, List, Optional, Tuple, Dict, Any
from skimage.filters import threshold_local
import aiohttp
import json
import base64


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
    image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    thresh_val, binary_image = cv2.threshold(
        image, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU
    )
    return thresh_val, binary_image


filters = [lambda x: x, bw_scanner, otsu_scanner]
reader = easyocr.Reader(["ru"], gpu=True)


def ocr(
    image: np.ndarray,
    workers: int = 2,
    text_threshold: float = 0.6,
) -> List[str]:
    global reader
    return reader.readtext(
        image, detail=0, text_threshold=text_threshold, workers=workers
    )


async def send_to_llm(
    image_np: np.ndarray,
    prompt_text: str,
    api_key: str,
    model: str = "mistral-small-latest"
) -> dict:
    if image_np.dtype != np.uint8:
        image_np = image_np.astype(np.uint8)
    
    if image_np.shape[2] == 3:
        image_pil = Image.fromarray(image_np[..., ::-1]) 
    else:
        image_pil = Image.fromarray(image_np)
    
    buffer = io.BytesIO()
    image_pil.save(buffer, format="JPEG", quality=90)
    base64_str = base64.b64encode(buffer.getvalue()).decode("utf-8")

    messages = [
        {
            "role": "user",
            "content": [
                {"type": "text", "text": prompt_text},
                {
                    "type": "image_url",
                    "image_url": {
                        "url": f"data:image/jpeg;base64,{base64_str}"
                    }
                }
            ]
        }
    ]
    
    async with aiohttp.ClientSession() as session:
        async with session.post(
            "https://api.mistral.ai/v1/chat/completions",
            headers={"Authorization": f"Bearer {api_key}"},
            json={
                "model": model,
                "messages": messages,
                "response_format": {"type": "json_object"}
            }
        ) as response:
            response.raise_for_status()
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

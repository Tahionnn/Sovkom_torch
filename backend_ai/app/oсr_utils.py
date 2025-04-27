import cv2
import io
import numpy as np
from typing import Any
import aiohttp
import json
import base64


def opencv_resize_maxcap(image: np.ndarray, max_cap: int = 300) -> np.ndarray:
    max_size = max(image.shape[1], image.shape[0])
    ratio = max_cap / max_size
    width = int(image.shape[1] * ratio)
    height = int(image.shape[0] * ratio)
    dim = (width, height)
    return cv2.resize(image, dim, interpolation=cv2.INTER_AREA)


async def send_to_llm(
    image_np: np.ndarray,
    prompt_text: str,
    api_key: str,
    model: str = "mistral-small-latest",
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
                    "image_url": {"url": f"data:image/jpeg;base64,{base64_str}"},
                },
            ],
        }
    ]

    async with aiohttp.ClientSession() as session:
        async with session.post(
            "https://api.mistral.ai/v1/chat/completions",
            headers={"Authorization": f"Bearer {api_key}"},
            json={
                "model": model,
                "messages": messages,
                "response_format": {"type": "json_object"},
            },
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

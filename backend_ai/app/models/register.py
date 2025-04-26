import mlflow
import mlflow.onnx
import torch
import onnx
import onnxruntime as rt
import numpy as np
import os
from mlflow import MlflowClient

def register_onnx_model(model_path: str):
    onnx_model = onnx.load(model_path)  # Загрузка ONNX-файла
    
    with mlflow.start_run():
        # Логирование модели
        mlflow.onnx.log_model(
            onnx_model=onnx_model,
            artifact_path="model",
            registered_model_name="multivae"  # Автоматическая регистрация
        )
        mlflow.log_param("framework", "ONNX")
        mlflow.log_param("top_k", 5)

def register_torch_model(model_path: str):
    # Загрузка модели (пример)
    model = torch.jit.load(model_path)  # Или torch.jit.load
    
    # Конвертация в TorchScript
    scripted_model = torch.jit.script(model)  # Или torch.jit.trace
    
    # Логирование в MLflow
    with mlflow.start_run():
        mlflow.pytorch.log_model(
            pytorch_model=scripted_model,
            artifact_path="model",
            registered_model_name="multivae-torch",
        )
        mlflow.log_param("framework", "PyTorch")
        mlflow.log_param("top_k", 5)

if __name__ == "__main__":
    mlflow.set_tracking_uri("http://localhost:5000")
    model_name = "multivae-torch"
    version = 1  # Ваша версия из лога
    client = MlflowClient()

    model_version = client.get_model_version(model_name, version)
    source = model_version.source

    print(f"Model version {version} source: {source}")
    print("Artifacts:")
    artifacts = client.list_artifacts(run_id=model_version.run_id, path="model")
    for artifact in artifacts:
        print(f"- {artifact.path}")
    register_torch_model("backend_ai/app/models/multivae.pt")  # Путь к вашим весам
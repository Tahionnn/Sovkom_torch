from fastapi import FastAPI
 
app = FastAPI()
 
@app.get("/predict")
def read_root():
    pass
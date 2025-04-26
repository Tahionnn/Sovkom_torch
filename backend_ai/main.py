import logging
import uuid
from logging.handlers import RotatingFileHandler
from venv import create

import uvicorn
from fastapi import FastAPI, Depends, HTTPException, UploadFile
from typing import Annotated
from fastapi.responses import JSONResponse, FileResponse
from pydantic import BaseModel

from db import Dbase

logging.basicConfig(level=logging.INFO, filename="py_log.log",filemode="w")

class LoginRequest(BaseModel):
    username: str
    password: str

db = Dbase()

def get_db():
    return db

app = FastAPI()

@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.get('/user/id/{id}')
async def get_user_id(id: int, db: Annotated[Dbase, Depends(get_db)]):
    # id = 2
    user = db.getUsersId(id)
    print(user)
    if not user:
        print(12312321)
        logging.info(f"FAILED ID: {id}")
        raise HTTPException(status_code=400, detail="There is no such user")
    logging.info(f"SUCCESS ID: {id}")
    return {"id": user[0], 'username': user[1]}

@app.get('/user/username/{username}')
async def get_user_id(username: str, db: Annotated[Dbase, Depends(get_db)]):
    print(username)
    user = db.getUserUsername(username)
    if not user:
        logging.info(f"FAILED  USERNAME: {username}")
        raise HTTPException(status_code=400, detail="There is no such user")
    logging.info(f"SUCCESS USERNAME: {username}")
    return {"id": user[0], 'username': user[1]}

@app.get("/users")
async def read_item(db: Annotated[Dbase, Depends(get_db)]):
    item = db.getUsers()
    dict_item = {'users': []}
    print(dict_item)
    for user in item:
        # print(user.id)
        dict_item['users'].append({'id':  user[0], 'username': user[1]})
    return dict_item


@app.post('/new_user')
async def new_user(credentials: LoginRequest, db: Annotated[Dbase, Depends(get_db)]):
    if db.getUsersId(credentials.username):
        logging.info(f"FAILDER NEW USER: {credentials.username} {credentials.password}")
        raise HTTPException(status_code=400, detail="Such a user has already been created")
    db.newUser(credentials.username, credentials.password)
    logging.info(f"SUCCESS NEW USER: {credentials.username} {credentials.password}")
    raise HTTPException(status_code=400, detail="Success")


@app.post('/check/{id}')
async def new_user(id: int, file: UploadFile, db: Annotated[Dbase, Depends(get_db)]):
    logging.info(f"NEW USER CHECK: {id}")
    data = await file.read()
    with open(f'./photo_check_{id}.png', "wb") as buffer:
        buffer.write(data)  # Записываем в файл
    return {"message": "Hello World"}


if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=8000)
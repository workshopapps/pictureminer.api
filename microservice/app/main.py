from fastapi import FastAPI

from routes import mine_image

app = FastAPI()

app.include_router(mine_image.router)
from fastapi import APIRouter

from app.schemas import ImageUpload

router = APIRouter(prefix='/api/microservice', tags=['/api/microservice'])


@router.get('/content')
async def get_image_content(image: ImageUpload):
    return image
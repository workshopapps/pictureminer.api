from fastapi import APIRouter, UploadFile

router = APIRouter(prefix='/api/microservice', tags=['/api/microservice'])


@router.get('/content')
async def get_image_content(image: UploadFile):
    image_filename, image_body = image.filename, image.file

    # call function to caption/generate text from model sub dir
    # possibly async
    # captioned_text = generate_caption(image_filename, image_body)
    captioned_text = '<dummy captioned text>'

    return {'text_description': captioned_text}

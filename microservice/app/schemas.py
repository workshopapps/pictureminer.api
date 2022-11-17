from pydantic import BaseModel


class ImageUpload(BaseModel):
    image: bytes

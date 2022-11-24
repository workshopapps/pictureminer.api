# Minergram

Find the guides, samples, and references you need to use the platform, generate context, and build applications on the Minergram Platform.

# Getting Started

## Getting Access to the API

Before we jump ahead we assume you hav already created an Minergram Account. If you haven't please visit [our webpage]("https://minergram.netlify.app/") to create an account. Access to the API is automatically allowed with your signup credentials,
no need for API keys.


# Mining your first image.
We've put together some examples on how to call the API using the a CLI (command line interface),an API client like postman and python
though any other programming languge would work too.
To begin calling the API, you must [signup]("https://minergram.netlify.app/") for an account through the website to get access.

## Example using a CLI

Simplest way to test our API is calling through a CLI.
Prerequisites:

- You must have an account
- A computer with internet access

#### Login

To get a JWT access token, you have to first make a login request with your signup credentials.
to make a cURL request, copy and paste or type:

```shell
curl --request POST \
  --url http://44.211.169.234:9000/api/v1/login \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "johndoe@gmail.com",
	"password": "password"
}'
```

replace <code>johndoe@gmail.com</code> with your email and <code>password</code> with your password.
sample response:

```json
{
  "status": "success",
  "code": 200,
  "message": "User login successful",
  "data": {
    "Username": "johndoe",
    "FirstName": "John",
    "LastName": "Doe",
    "Email": "johndoe@gmail.com",
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NjkyMjc2MDEsImlkIjoiT2JqZWN0SUQoXCI2MzdkMTdjOTQxMTg5ZmRjNmJiODkwZWVcIikifQ.sTzq6C2B86w94auIjqjSveJ55E6G8Iwa-E564gUdjrg",
    "TokenType": "bearer",
    "ApiCallCount": 8
  }
}
```

copy the "Token" value, that is the JWT access token we'll be using.

### Example request 1: mine image through url

To make the cURL request, copy and paste or type:

```shell
curl --request POST \
  --url http://44.211.169.234:9000/api/v1/mine-service/url \
  --header 'Content-Type: application/json' \
  --data '{
	"url": "any image url here"
}'
```

Assuming the url provided links to the image below.
<div><img src="https://static1.bigstockphoto.com/7/0/3/large1500/307281196.jpg" width=500></div>
the sample response:

```javascript
{
  "status": "success",
  "code": 201,
  "message": "mine image successful",
  "data": {
    "image_name": "elephants.png",
    "image_path": "https://mined-pictures.s3.amazonaws.com/9d674181d53ac16571d.png",
    "text_content": "A group of elephants eating grass",
    "date_created": "2022-11-23T17:20:51.002216449Z",
    "date_modified": "2022-11-23T17:20:51.002216449Z"
  }
}
```

### Example request 1: mine image through upload

to make the cURL request, copy and paste or type:

```shell
curl --request POST \
  --url http://44.211.169.234:9000/api/v1/mine-service/upload \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NjkyMjc2MDEsImlkIjoiT2JqZWN0SUQoXCI2MzdkMTdjOTQxMTg5ZmRjNmJiODkwZWVcIikifQ.sTzq6C2B86w94auIjqjSveJ55E6G8Iwa-E564gUdjrg' \
  --header 'Content-Type: multipart/form-data' \
  --header 'content-type: multipart/form-data; boundary=---011000010111000001101001' \
  --form image=elephants.jpg
```

Assuming <code>elephants.jpg</code> is the same image as above,  
the sample response:

```json
{
  "status": "success",
  "code": 201,
  "message": "mine image successful",
  "data": {
    "image_name": "elephants.jpg",
    "image_path": "https://mined-pictures.s3.amazonaws.com/981d53ac1d67416571d.png",
    "text_content": "A group of elephants eating grass",
    "date_created": "2022-11-23T17:20:51.002216449Z",
    "date_modified": "2022-11-23T17:20:51.002216449Z"
  }
}
```


## Using an API Client like Postman.

#### Uploading a raw image

- Create a new request with our the url `http://44.211.169.234:9000/api/v1/mine-service/upload`.

  <div><img src="https://firebasestorage.googleapis.com/v0/b/colour-switchboard.appspot.com/o/Screenshot%202022-11-24%20at%2019.04.08.png?alt=media&token=7f4a0e49-02ec-47b7-98e5-b1ca3c3da3d1" width=500></div>
- Set the authorization headers (this token is used to verify your identity. You can get the token at **insert website here**)
  []
  <div><img src="https://firebasestorage.googleapis.com/v0/b/colour-switchboard.appspot.com/o/Screenshot%202022-11-24%20at%2019.04.42.png?alt=media&token=6be1c646-3837-47da-a687-e8f016d084ef" width=500></div>
- ##### Upload the file as `multipart/form-data`

  > set the body as formdata

  <div><img src="https://firebasestorage.googleapis.com/v0/b/colour-switchboard.appspot.com/o/Screenshot%202022-11-24%20at%2019.05.24.png?alt=media&token=fbfbe375-bc69-4cc9-bab0-8fb9c2502fa3" width=500></div>

  > ...and select the `file` option the first row

  <div><img src="https://firebasestorage.googleapis.com/v0/b/colour-switchboard.appspot.com/o/Screenshot%202022-11-24%20at%2019.05.28.png?alt=media&token=03b8eb9c-e5e4-4cf2-857d-ab10fd9bd831" width=500></div>

  > Select a file & enter the field name your API uses for the file field then hit send.

  <div><img src="https://firebasestorage.googleapis.com/v0/b/colour-switchboard.appspot.com/o/Screenshot%202022-11-24%20at%2019.06.39.png?alt=media&token=454bfd7a-dcc7-4bf8-9d31-a0fb5c6829f4" width=500></div>

  you should a response looking some like...

  ```json
    {
  "status": "success",
  "code": 201,
  "message": "mine image successful",
  "data": {
    "image_name": "elephants.jpg",
    "image_path": "https://mined-pictures.s3.amazonaws.com/981d53ac1d67416571d.png",
    "text_content": "A group of elephants eating grass",
    "date_created": "2022-11-23T17:20:51.002216449Z",
    "date_modified": "2022-11-23T17:20:51.002216449Z"
  }
    }
  ```

## Example using Python

The other way to test is to use Python.
Prerequisites:

- You must have an account
- Python installation
- A computer with internet access

### Login

To get a JWT access token, you have to first make a login request with your signup credentials.
to do this with python, copy and paste or type:

NOTE: make sure to run <code>pip install requests</code> first if you don't have it installed.

```python
import requests

url = "http://44.211.169.234:9000/api/v1/login"

payload = "{\n\t\"email\": \"johndoe@gmail.com\",\n\t\"password\": \"password\"\n}"
headers = {'Content-Type': 'application/json'}

response = requests.request("POST", url, data=payload, headers=headers)

print(response.text)
```

replace <code>johndoe@gmail.com</code> with your email and <code>password</code> with your password.
sample response:

```javascript
{
  "status": "success",
  "code": 200,
  "message": "User login successful",
  "data": {
    "Username": "johndoe",
    "FirstName": "John",
    "LastName": "Doe",
    "Email": "johndoe@gmail.com",
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NjkyMjc2MDEsImlkIjoiT2JqZWN0SUQoXCI2MzdkMTdjOTQxMTg5ZmRjNmJiODkwZWVcIikifQ.sTzq6C2B86w94auIjqjSveJ55E6G8Iwa-E564gUdjrg",
    "TokenType": "bearer",
    "ApiCallCount": 8
  }
}
```

copy the "Token" value, that is the JWT access token we'll be using.

### Example request 1: mine image through url

in a python file, copy and paste or type:

```python
import requests

url = "http://44.211.169.234:9000/api/v1/mine-service/url"

payload = "{\n\t\"url\": \"image url here\"\n}"
headers = {
    'Content-Type': "application/json",
    'Authorization': "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NjkyNTEzNTUsImlkIjoiT2JqZWN0SUQoXCI2MzdkMTdjOTQxMTg5ZmRjNmJiODkwZWVcIikifQ.SyW5uq8pgwc2qlFY-PusHTL9HFYUDEoVKMvlTF57e0E"
    }

response = requests.request("POST", url, data=payload, headers=headers)

print(response.text)
```

Assuming the url provided links to the image below.
<div><img src="https://static1.bigstockphoto.com/7/0/3/large1500/307281196.jpg" width=500></div>  
the sample response:

```json
{
  "status": "success",
  "code": 201,
  "message": "mine image successful",
  "data": {
    "image_name": "elephants.png",
    "image_path": "https://mined-pictures.s3.amazonaws.com/9d674181d53ac16571d.png",
    "text_content": "A group of elephants eating grass",
    "date_created": "2022-11-23T17:20:51.002216449Z",
    "date_modified": "2022-11-23T17:20:51.002216449Z"
  }
}
```

### Example request 1: mine image through upload

To make the python request, copy and paste or type:

```python
import requests

url = "http://44.211.169.234:9000/api/v1/mine-service/upload"

payload = "-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"image\"\r\nvalue=elephants.jpg\r\n\r\n-----011000010111000001101001--\r\n"
headers = {
    'Content-Type': "multipart/form-data",
    'Authorization': "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NjkyMjc2MDEsImlkIjoiT2JqZWN0SUQoXCI2MzdkMTdjOTQxMTg5ZmRjNmJiODkwZWVcIikifQ.sTzq6C2B86w94auIjqjSveJ55E6G8Iwa-E564gUdjrg",
    'content-type': "multipart/form-data; boundary=---011000010111000001101001"
    }

response = requests.request("POST", url, data=payload, headers=headers)

print(response.text)
```

Assuming <code>elephants.jpg</code> is the same image as above,  
the sample response:

```json
{
  "status": "success",
  "code": 201,
  "message": "mine image successful",
  "data": {
    "image_name": "elephants.jpg",
    "image_path": "https://mined-pictures.s3.amazonaws.com/981d53ac1d67416571d.png",
    "text_content": "A group of elephants eating grass",
    "date_created": "2022-11-23T17:20:51.002216449Z",
    "date_modified": "2022-11-23T17:20:51.002216449Z"
  }
}
```

## Intergrating in your website using a library like axios.

#### Requirements

- [axios]()
- Some knowledge of html and Javascript.

In your `index.html` file copy the following code

```html
<div className="App">
  <h1>Hello Minergram</h1>

  <h2>Start mining to get some context</h2>

  <input type="file" name="file" onChange="function(e){handleFile(e)}" />

  <button onClick="function(e){handleUpload(e)}">Upload</button>
</div>
```

then in the corresponding `main.js` file configure the following code snippet to suit your frontend/framework.

```js
import  axios  from  "axios";

let file;
handleFile(e) {

    file = e.target.files[0];
}

async  handleUpload(e) {

    console.log(file);

    await  uploadImage(file);

}


const  uploadImage  =  async  file  =>  {

    try  {

        console.log("Upload Image",  file);

        const  formData  =  new  FormData();

        formData.append("image",  file);
        
        
        const  config  =  {

            headers:  {
                "content-type":  "multipart/form-data"
                "authorisation": `Bearer ${your_bearer_token_here}`
            }

        };

        const  API  =  "api/v1/mine_service/upload";

        const  HOST  =  "http://44.211.169.234:9000";

        const  url  =  `${HOST}/${API}`;



        const  result  =  await  axios.post(url,  formData,  config);

        console.log("Result: ",  result);

    }  catch  (error)  {
        console.error(error);
    }

};
```

### Not familiar with any of these terms ?

[Minergram]("https://minergram.netlify.app/") provides a nice GUI to implement all of these. Head over there for a demo

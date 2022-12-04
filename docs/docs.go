// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/mined-images": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "this returns the mined images of all users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "this returns the mined images of all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.MinedImage"
                            }
                        }
                    }
                }
            }
        },
        "/admin/users": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "List all users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "List all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.User"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/health": {
            "get": {
                "description": "Responds with the server status as JSON.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Checks the status of the server",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utility.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "Send a dummy post request to test the status of the server",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Checks the status of the server",
                "parameters": [
                    {
                        "description": "Ping JSON",
                        "name": "ping",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Ping"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utility.Response"
                        }
                    }
                }
            }
        },
        "/forgot-password": {
            "post": {
                "description": "Send a dummy post request to test the status of the server",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Forgot Password"
                ],
                "summary": "Checks the status of the forgot passoword",
                "parameters": [
                    {
                        "description": "Ping JSON",
                        "name": "ping",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PasswordForgot"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utility.Response"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Logs in a User",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Login User",
                "parameters": [
                    {
                        "description": "User Login",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.UserLogin"
                        }
                    }
                }
            }
        },
        "/mine-service/upload": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Send a post request containing a file an receives a response of its context content.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Mine-Service"
                ],
                "summary": "Mines an uploaded image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "image",
                        "name": "os.File",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.MineImageResponse"
                        }
                    }
                }
            }
        },
        "/reset": {
            "post": {
                "description": "Send a post request to reset th password of the user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Resests the password of the user",
                "parameters": [
                    {
                        "description": "Ping JSON",
                        "name": "ping",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PasswordReset"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utility.Response"
                        }
                    }
                }
            }
        },
        "/signup": {
            "post": {
                "description": "Creates an account for a new user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Signs Up a User",
                "parameters": [
                    {
                        "description": "User Signup",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UserSignUp"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.UserResponse"
                        }
                    }
                }
            }
        },
        "/update-user": {
            "patch": {
                "description": "Updates a User's information - email,firstName,lastName,password- Bearer token and email required",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update User",
                "parameters": [
                    {
                        "description": "User Update",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UpdateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.UserLogin"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.MineImageResponse": {
            "type": "object",
            "properties": {
                "date_created": {
                    "type": "string"
                },
                "date_modified": {
                    "type": "string"
                },
                "image_name": {
                    "type": "string"
                },
                "image_path": {
                    "type": "string"
                },
                "text_content": {
                    "type": "string"
                }
            }
        },
        "model.MinedImage": {
            "type": "object",
            "properties": {
                "dateCreated": {
                    "type": "string"
                },
                "dateModified": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "imageKey": {
                    "type": "string"
                },
                "imageName": {
                    "type": "string"
                },
                "imagePath": {
                    "type": "string"
                },
                "textContent": {
                    "type": "string"
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "model.PasswordForgot": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "model.PasswordReset": {
            "type": "object",
            "required": [
                "confirm_password",
                "email",
                "password"
            ],
            "properties": {
                "confirm_password": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "model.Ping": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "model.UpdateUser": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "confirm_password": {
                    "type": "string"
                },
                "current_password": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "new_password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "date_created": {
                    "type": "string"
                },
                "date_updated": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "profile_key": {
                    "type": "string"
                },
                "profile_url": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.UserLogin": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "model.UserResponse": {
            "type": "object",
            "properties": {
                "apiCallCount": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "profileKey": {
                    "type": "string"
                },
                "profileUrl": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "tokenType": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.UserSignUp": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "utility.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "error": {
                    "description": "for errors that occur even if request is successful"
                },
                "extra": {},
                "message": {
                    "type": "string"
                },
                "name": {
                    "description": "name of the error",
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header \"Bearer \u003cadd access token here\u003e\""
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "discripto.hng.tech/api1",
	BasePath:         "/api/v1/",
	Schemes:          []string{"https"},
	Title:            "Minergram",
	Description:      "A picture mining service API in Go using Gin framework.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

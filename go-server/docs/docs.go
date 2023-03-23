// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "autolarry55@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/forgot-password/init/{email}": {
            "post": {
                "description": "An otp code is sent to the email if the user account existed.\nAn otp ID is returned, which must submitted alongside the otpCode sent to the mail to the Forgot password reset endpoint.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Forget password endpoint",
                "operationId": "ForgotPasswordInit",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email address",
                        "name": "email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/authentication.OTPCreationSuccessResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "No Account Found",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/forgot-password/resend/{email}": {
            "post": {
                "description": "An otp code is sent to the email if the user account existed.\nAn otp ID is returned, which must submitted alongside the otpCode sent to the mail to the Forgot password reset endpoint.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Resend OTP code for Forget password",
                "operationId": "ResendForgetPasswordCode",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email address",
                        "name": "email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/authentication.OTPCreationSuccessResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "No Account Found",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/forgot-password/reset": {
            "post": {
                "description": "The Endpoint resets the user password using the otp code sent to the user's email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Complete Forget password reset",
                "operationId": "ResetPassword",
                "parameters": [
                    {
                        "description": "OtpID data",
                        "name": "ForgetAndResetPasswordRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/authentication.ForgetAndResetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "404": {
                        "description": "No Account Found",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "409": {
                        "description": "User already verified",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/games/genres/add": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Adds a new game genre, the slug is generated from the title, and it must be unique",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "games"
                ],
                "summary": "Adds a new game genre",
                "operationId": "addGenre",
                "parameters": [
                    {
                        "description": "addGenre request",
                        "name": "addGenre",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/games.AddGenreRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/games.AddGenreRes"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "409": {
                        "description": "Genre with the same slug already exists",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/login": {
            "post": {
                "description": "Returns a signed JSON Web Token that can be used to talk to secured endpoints",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Login endpoint for all users",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "login request",
                        "name": "loginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/authentication.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/authentication.LoginSuccessResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "426": {
                        "description": "Account is inactive",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/ping": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ping"
                ],
                "summary": "Show the status of server.",
                "operationId": "ping",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/v1/signup": {
            "post": {
                "description": "Creates a new User/Moderator on the system. The Moderator will need to be manually activated by an existing admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Onboarding"
                ],
                "summary": "Signup endpoint for all users and moderators",
                "operationId": "signup",
                "parameters": [
                    {
                        "description": "signup request",
                        "name": "signUpRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/authentication.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "409": {
                        "description": "User already exists",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/verify-account/init/{email}": {
            "post": {
                "description": "An otp code is sent to the email if the user account existed.\nAn otp ID is returned, which must submitted alongside the otpCode sent to the mail to the VerifyAccount endpoint.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Onboarding"
                ],
                "summary": "Initiate user email verification",
                "operationId": "initVerifyAccount",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email address",
                        "name": "email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/authentication.OTPCreationSuccessResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "No Account Found",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "409": {
                        "description": "User already verified",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/verify-account/resend/{email}": {
            "post": {
                "description": "An otp code is sent to the email if the user account existed.\nAn otp ID is returned, which must submitted alongside the otpCode sent to the mail to the VerifyAccount endpoint.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Onboarding"
                ],
                "summary": "Resend OTP code for account verification",
                "operationId": "ResendVerificationCode",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email address",
                        "name": "email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/authentication.OTPCreationSuccessResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "No Account Found",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "409": {
                        "description": "User already verified",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/verify-account/verify": {
            "post": {
                "description": "The Endpoint verifies the user account using the otp code sent to the user's email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Onboarding"
                ],
                "summary": "Complete account verification",
                "operationId": "VerifyAccount",
                "parameters": [
                    {
                        "description": "OtpID data",
                        "name": "verifyAccountRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/authentication.VerifyAccountRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/main.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "404": {
                        "description": "No Account Found",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "409": {
                        "description": "User already verified",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.JSONErrorRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "authentication.ForgetAndResetPasswordRequest": {
            "type": "object",
            "required": [
                "confirmPassword",
                "email",
                "otpCode",
                "password",
                "tokenID"
            ],
            "properties": {
                "confirmPassword": {
                    "type": "string",
                    "minLength": 8
                },
                "email": {
                    "type": "string"
                },
                "otpCode": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "tokenID": {
                    "type": "string"
                }
            }
        },
        "authentication.LoginRequest": {
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
        "authentication.LoginSuccessResponse": {
            "type": "object",
            "properties": {
                "jwt": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/authentication.UserDetail"
                }
            }
        },
        "authentication.OTPCreationSuccessResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "otpId": {
                    "type": "string"
                }
            }
        },
        "authentication.SignUpRequest": {
            "type": "object",
            "required": [
                "email",
                "firstName",
                "lastName",
                "password",
                "phone",
                "role"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "role": {
                    "type": "string",
                    "enum": [
                        "user",
                        "moderator"
                    ]
                }
            }
        },
        "authentication.UserDetail": {
            "type": "object",
            "required": [
                "firstName",
                "lastName",
                "role"
            ],
            "properties": {
                "firstName": {
                    "type": "string"
                },
                "isVerified": {
                    "type": "boolean"
                },
                "lastName": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "role": {
                    "type": "string",
                    "enum": [
                        "user",
                        "admin",
                        "moderator"
                    ]
                }
            }
        },
        "authentication.VerifyAccountRequest": {
            "type": "object",
            "required": [
                "email",
                "otpCode",
                "tokenID"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "otpCode": {
                    "type": "string"
                },
                "tokenID": {
                    "type": "string"
                }
            }
        },
        "games.AddGenreRequest": {
            "type": "object",
            "required": [
                "desc",
                "title"
            ],
            "properties": {
                "desc": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "games.AddGenreRes": {
            "type": "object",
            "properties": {
                "slug": {
                    "type": "string"
                }
            }
        },
        "main.JSONErrorRes": {
            "type": "object",
            "properties": {
                "error": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "main.JSONResult": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Game Review API",
	Description:      "This is an Api AuthService for Cool Game Review Api.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

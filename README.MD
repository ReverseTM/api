# Startup Instructions

## Steps to Run

1. Navigate to the `deployments` folder.

2. In the terminal, execute the following command:

```bash
docker-compose up
```

# Configuration Setup
If needed, you can create your own configuration file in the configs folder and specify the path to this file in the CONFIG_PATH environment variable for the go-app service.

## Configuration File Structure
The configuration file should be in YAML format and have the following structure:

  ```yaml
  env: local  # Application environment. Possible values: local, dev, prod (default is local)

  jwt:
    secret: "secret"  # Key for encrypting JWT tokens (default is "secret")
    token_ttl: "1h"   # Time-to-live for JWT tokens (default is 1 hour)

  storage:
    storage_path: "database_connection_address"  # Address for connecting to the database (mandatory parameter)
    space_name: "database_space_name"  # Name of the space in the database (mandatory parameter)

  http_server:
    address: ":8080"  # Address and port for the server to listen on (default is :8080)

  ```

# API Documentation

## Overview

This document provides details about the API endpoints.

- **`/api/login`**: Allows users to authenticate by providing a username and password and, upon successful authentication, receive an authentication token.
- **`/api/read`**: Enables authorized users to read data from the server.
- **`/api/write`**: Allows authorized users to write data to the server.

Each endpoint requires proper authentication and data formatting to ensure successful interactions with the API.

## Endpoints

### Login

**Endpoint:** `POST /api/login`

This endpoint is used to authenticate users. Clients must provide valid credentials (username and password) to receive an authentication token.

#### Request

- **Content-Type:** `application/json`

- **Request Body:**

The request body must be a JSON object with the following fields:

```json
  {
    "username": "string",
    "password": "string"
  }
```
  
- username (string, required): The username for authentication.
- password (string, required): The password for authentication.

#### Example Request:

```http
POST /api/login HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "username": "user123",
    "password": "password123"
 }
```

#### Response

##### Success (HTTP 200 OK):
If the authentication is successful, the server responds with a JSON object containing the authentication token.

- **Content-Type:** `application/json`
- **Response body:**

```json
  {
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
 ```

- token (string): The authentication token that can be used for authorization in subsequent requests.

##### Failure (HTTP 400 Bad Request):
If the request is malformed or missing required fields, the server will return an error with a 400 status code.

- **Content-Type:** `application/json`
- **Response Body:**

```json
  {
    "error": "invalid request"
  }
```
  
- error (string): An error message indicating that the request was invalid or missing required parameters.

##### Failure (HTTP 401 Unauthorized):
If the provided credentials are incorrect or authentication fails, the server will return an error with a 401 status code.

- **Content-Type:** `application/json`
- **Response Body:**
  
```json
  {
     "error": "invalid credentials"
  }
```

- error (string): An error message indicating that authentication failed due to invalid credentials.

### Read

**Endpoint:** `POST /api/read`

This endpoint is used to retrieve data based on a set of keys. Clients need to provide the necessary keys in the request body and include a Bearer token in the Authorization header to access the endpoint.

#### Request

- **Content-Type:** `application/json`
- **Authorization Header:** `Bearer <token>`

- **Request Body:**

The request body must be a JSON object with the following field:

```json
  {
    "keys": ["string"]
  }
```

- keys (array of strings, required): A list of keys for which data needs to be retrieved.

#### Example Request:

```http
POST /api/read HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer your_jwt_token

{
  "keys": ["key1", "key2"]
}
```

#### Response

##### Success (HTTP 200 OK):
If the data retrieval is successful, the server responds with a JSON object containing the retrieved data.

- **Content-Type:** `application/json`
- **Response Body:**

```json
{
  "data": {
      "key1": "value1",
      "key2": "value2"
  }
}
```

- data (array of objects): A list of key-value pairs where each object represents the retrieved data for a specific key.
- key (string): The key for which data was retrieved.
- value (string): The value associated with the key.

##### Failure (HTTP 400 Bad Request):
If the request body is malformed or missing required fields, the server will return an error with a 400 status code.

- **Content-Type:** `application/json`
- **Response Body:**

```json
{
  "error": "invalid request"
}
```

- error (string): An error message indicating that the request was invalid or missing required parameters.

##### Failure (HTTP 401 Unauthorized):
If the Authorization header is missing, invalid, or if the Bearer token is invalid or expired, the server will return an error with a 401 status code.

- **Content-Type:** `application/json`
- **Response Body:**

```json
{
  "error": "auth header is missing"
}
```

- if authorization header is missing.

```json
{
   "error": "invalid authorization header format"
}
```

- if authorization header is invalid

```json
{
   "error": "token signature is invalid"
}
```

- if token signature is invalid

```json
{
   "error": "invalid authorization token"
}
```

- if token is invalid

##### Failure (HTTP 500 Internal Server Error):
If there is an internal server error while processing the request, the server will return an error with a 500 status code.

- **Content-Type:** `application/json`
- **Response Body:**

```json
{
  "error": "internal server error"
}
```

- error (string): An error message indicating that an internal server error occurred while processing the request.

### Write

**Endpoint:** `POST /api/write`

- This endpoint is used to write data to the server. Clients need to provide the necessary data in the request body and include a Bearer token in the Authorization header to access the endpoint.

#### Request

- **Content-Type:** `application/json`
- **Authorization Header:** `Bearer <token>`

- **Request Body:**

The request body must be a JSON object with the following field:

```json
  {
    "data": {
        "key1": "value1",
        "key2": "value2"
    }
  }
```

- data (array of objects, required): A list of key-value pairs to be written to the server.
- key (string): The key for the data entry.
- value (any): The value associated with the key

#### Example Request:

```http
POST /api/write HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer your_jwt_token

{
  "data": {
      "key1": "value1",
      "key2": 10
    }
}
```

#### Response

##### Success (HTTP 200 OK):
If the data is successfully written to the server, the server responds with a confirmation message.

- **Content-Type:** `application/json`
- **Response Body:**

```json
{
  "status": "success"
}
```

- response (string): A confirmation message indicating that the data was successfully written.

##### Failure (HTTP 400 Bad Request):
If the request body is malformed or missing required fields, the server will return an error with a 400 status code.

- **Content-Type:** `application/json`
- **Response Body:**

```json
{
"error": "invalid request"
}
```

- error (string): An error message indicating that the request was invalid or missing required parameters.

##### Failure (HTTP 401 Unauthorized):
If the Authorization header is missing, invalid, or if the Bearer token is invalid or expired, the server will return an error with a 401 status code.

- **Content-Type:** `application/json`
- **Response Body:**

```json
{
  "error": "auth header is missing"
}
```

- if authorization header is missing.

```json
{
   "error": "invalid authorization header format"
}
```

- if authorization header is invalid

```json
{
   "error": "token signature is invalid"
}
```

- if token signature is invalid

```json
{
   "error": "invalid authorization token"
}
```

- if token is invalid

##### Failure (HTTP 500 Internal Server Error):
If there is an internal server error while processing the request, the server will return an error with a 500 status code.

- **Content-Type:** `application/json`
- **Response Body:**

```json
{
  "error": "internal server error"
}
```

- error (string): An error message indicating that an internal server error occurred while processing the request.
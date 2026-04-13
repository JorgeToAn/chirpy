# Chirpy

Chirpy is social network where you can share all your thoughts with anyone around the world!

## Features

- REST API
- File server
- User authentication via passwords and JWTs
- Webhooks to handle when a user upgrades to Chirpy Red

## Installation

### Requirements

You need the following programs to run Chirpy:

- Go v1.25.0 or newer
- PostgreSQL
- [Goose](https://github.com/pressly/goose) v3.26.0 or newer

### Setup

Copy this repository with the following command or download the ZIP file directly:

```bash
git clone git@github.com:JorgeToAn/chirpy.git
```

Create an empty database for Chirpy and apply the migrations using goose (you need to run it in the sql/schema directory or use the -d flag):

```bash
goose postgres "postgres://user:password@host:port/chirpy" up
```

Create an .env file that contains the following:

- **JWT_SECRET**: A secret used to sign and verify the JWTs you hand to users for authentication, you can create one with `openssl rand -base64 64`
- **POLKA_KEY**: A Polka API key used to verify webhook requests are sent from Polka themselves
- **DB_URL**: A connection URL to a Postgres database
- **PLATFORM** (optional): Used to indicate if the environment is for development if the value is "dev"

Finally, build a binary with `go build` and run it!

## Development

If you need to create new queries to expand the functionalities of this project, I recommend using [sqlc](https://sqlc.dev/).

## API Documentation

### Health

Shows whether the server is healthy and ready.

**Endpoint:** `GET /api/healthz`

**Authentication:** Not required

### Reset environment

**NOT FOR PRODUCTION.**

Resets metrics and clears all users from the database.

**Requires** that the PLATFORM is `dev`.

**Endpoint:** `POST /admin/reset`

**Authentication:** Not required

### Get Metrics

Shows how many hits have been made to the file server.

**Endpoint:** `GET /admin/metrics`

**Authentication:** Not required

### Create User

Create a new user account.

**Endpoint:** `POST /api/users`

**Authentication:** Not required

**Request Body:**

```json
{
    "email": "jesse@breakingbad.com",
    "password": "123456"
}
```

**Request Schema:**

| Field | Type | Required | Description | Constraints |
| ----- | ---- | -------- | ----------- | ----------- |
| email | string | Yes | User's email address | Unique |
| password | string | Yes | User's password | None |

**Response:**

```json
{
  "id": "725f8d4a-08b3-4486-81ee-ff454600eb1a",
  "created_at": "2026-04-12T16:34:42.795796Z",
  "updated_at": "2026-04-12T16:34:42.795796Z",
  "email": "jesse@breakingbad.com",
  "is_chirpy_red": false
}
```

### Update User Credentials

Update user's email and password.

**Endpoint:** `PUT /api/users`

**Authentication:** Required

**Request Body:**

```json
{
    "email": "jesse@elcamino.com",
    "password": "7890"
}
```

**Request Schema:**

| Field | Type | Required | Description | Constraints |
| ----- | ---- | -------- | ----------- | ----------- |
| email | string | Yes | User's email address | Unique |
| password | string | Yes | User's password | None |

**Response:**

```json
{
  "id": "725f8d4a-08b3-4486-81ee-ff454600eb1a",
  "created_at": "2026-04-12T16:34:42.795796Z",
  "updated_at": "2026-04-12T16:43:03.744785Z",
  "email": "jesse@elcamino.com",
  "is_chirpy_red": false
}
```

### Log In

Log in with an email and password.

**Endpoint:** `POST /api/login`

**Authentication:** Not required

**Request Body:**

```json
{
    "email": "jesse@breakingbad.com",
    "password": "123456"
}
```

**Request Schema:**

| Field | Type | Required | Description | Constraints |
| ----- | ---- | -------- | ----------- | ----------- |
| email | string | Yes | User's email address | Unique |
| password | string | Yes | User's password | None |

**Response:**

```json
{
  "id": "725f8d4a-08b3-4486-81ee-ff454600eb1a",
  "created_at": "2026-04-12T16:34:42.795796Z",
  "updated_at": "2026-04-12T16:34:42.795796Z",
  "email": "jesse@breakingbad.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiNzI1ZjhkNGEtMDhiMy00NDg2LTgxZWUtZmY0NTQ2MDBlYjFhIiwiZXhwIjoxNzc2MDQwNzI1LCJpYXQiOjE3NzYwMzcxMjV9.h0V0zPImLg3TP7dNIKqJIUJo--E_y6XMbFoOCXsvITo",
  "refresh_token": "f806400fa3def3f821bda164df9db6937497767f4ec2576c97aec713211a0be7"
}
```

### Refresh Access

Use user's refresh token to create a new access token.

**Endpoint:** `POST /api/refresh`

**Authentication:** Required

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiNzI1ZjhkNGEtMDhiMy00NDg2LTgxZWUtZmY0NTQ2MDBlYjFhIiwiZXhwIjoxNzc2MDQxNDQyLCJpYXQiOjE3NzYwMzc4NDJ9.ea2h7F5kJVEDkKWY8WuFevbMC_mcYIPB8-ZMQ58TASg"
}
```

### Revoke Refresh Token

Revoke a refresh token so it can no longer be used.

***Endpoint:** `POST /api/revoke`

**Authentication:** Required

### Create Chirp

Create a new Chirp.

**Endpoint:** `POST /api/chirps`

**Authentication:** Required

**Request Body:**

```json
{
    "body": "Hello brave new world!"
}
```

**Request Schema:**

| Field | Type | Required | Description | Constraints |
| ----- | ---- | -------- | ----------- | ----------- |
| body | string | Yes | Chirp's body | None |

**Response:**

```json
{
  "id": "f1a12ddb-5527-439e-bd1d-59f8a6c1144a",
  "created_at": "2026-04-12T16:56:36.772375Z",
  "updated_at": "2026-04-12T16:56:36.772375Z",
  "body": "Hello brave new world!",
  "user_id": "725f8d4a-08b3-4486-81ee-ff454600eb1a"
}
```

### Get Chirps

Get a collection of chirps.

**Endpoint:** `GET /api/chirps`

**Authentication:** Not required

**Query Parameters:**

| Parameter | Type | Required | Description |
| --------- | ---- | -------- | ----------- |
| author_id | string | No | Filters to chirps made by the author with the specified ID |
| sort | string | No | Orders chirps by creation date `asc` (default) or `desc` |

**Response:**

```json
[
  {
    "id": "f1a12ddb-5527-439e-bd1d-59f8a6c1144a",
    "created_at": "2026-04-12T16:56:36.772375Z",
    "updated_at": "2026-04-12T16:56:36.772375Z",
    "body": "Hello brave new world!",
    "user_id": "725f8d4a-08b3-4486-81ee-ff454600eb1a"
  },
  {
    "id": "7d0fb651-6a93-450b-8876-869e0aaf78f8",
    "created_at": "2026-04-12T17:49:13.370737Z",
    "updated_at": "2026-04-12T17:49:13.370737Z",
    "body": "Today is gonna be the day",
    "user_id": "e6da1f89-7f01-484b-8410-950d3117a44e"
  }
]
```

### Get Chirp by ID

Get a chirp with a specific ID.

**Endpoint:** `GET /api/chirps/{id}`

**Authentication:** Not required

**Path Parameters:**

| Parameter | Type | Required | Description |
| --------- | ---- | -------- | ----------- |
| id | string | Yes | The chirp's ID |

**Response:**

```json
{
  "id": "f1a12ddb-5527-439e-bd1d-59f8a6c1144a",
  "created_at": "2026-04-12T16:56:36.772375Z",
  "updated_at": "2026-04-12T16:56:36.772375Z",
  "body": "Hello brave new world!",
  "user_id": "725f8d4a-08b3-4486-81ee-ff454600eb1a"
}
```

### Delete Chirp

Delete a chirp with the specified ID.

**Endpoint:** `DELETE /api/chirps/{chirpID}`

**Authentication:** Required

**Path Parameters:**

| Parameter | Type | Required | Description |
| --------- | ---- | -------- | ----------- |
| id | string | Yes | The chirp's ID |

### Polka Webhook

Webhook handler for Polka events

**Endpoint:** `POST /api/polka/webhooks`

**Authentication:** Required

**Request Body:**

```json
{
    "event": "user.upgraded",
    "data": {
        "user_id": "d08f502d-a344-411c-b6e7-1342f411b885"
    }
}
```

**Request Schema:**

| Field | Type | Required | Description | Constraints |
| ----- | ---- | -------- | ----------- | ----------- |
| event | string | Yes | Event's name | None |
| data.user_id | string | Yes | User's ID | None |

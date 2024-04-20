# eau-de-go

## What is eau-de-go?
`eau-de-go`is a template web backend project written in Go. 
It is designed to be a starting point for web projects that require a backend. 
We named it "eau-de-go" because we hope that it will make your Go project flow like water.

It is built on top of the workflow proposed by [sqlc](https://docs.sqlc.dev/en/latest/), 
and roughly follows domain-driven design with repository pattern.

A service layer is also provided to provide more complex business logic implementation, 
especially since we are relying on the generated code for repository operations; 
so you can implement custom business logic in the service layer without changing the repository code generated by sqlc.

In the transport layer, we use [mux](https://github.com/gorilla/mux) for routing and for query parsing, 
which interacts with the service layer to handle requests.

Stateless authentication is also implemented using JWT, including middleware to protect resource routes.
The JWT signing is done using asymmetric keys, so that it may support microservices architecture better.
A in-memory JWT key store is included for development, with cloud-based key store support in the future.

A set of API endpoints are included for sign up, sign in, and token refresh mechanism, 
as well as sample API endpoints for managing user resources.

A centralized settings package is also included to manage application settings, 
and allows configuration using environment variables.


### Features
- Makefile for common tasks
- Docker configuration
- Centralized environment variables and settings
- SQLC for repository
- Service layer for business logic
- Mux for routing
- JWT for authentication
- Auth API endpoints
- Basic user model and API endpoints
- User email verification


## Development

### Quick Setup
Set up environment variables:
```bash
make dot-env
```
This will create a `.env` file and a `.env.docker` file in the root directory, 
make sure the environment variables are set correctly for your development environment.

Set up the database and run the migrations:
```bash
make run-migrations
```

Run the server:
```bash
make run-server
```

To run all tests:
```bash
make run-tests
```

### Docker
For development with Docker, you can use the following commands:
```bash
make docker-build
make docker-up
make docker-down
```

## Authentication
The included authentication mechanism is a stateless JWT-based authentication. 
Upon successful sign in, the user is provided with a JWT refresh token and JWT access token.

The access token is used to authenticate requests to protected resources, 
and must be included in the `Authorization` header of the request as a `Bearer` token. 

The access token is short-lived, and must be refreshed using the refresh token. 
The refresh token is implicitly included in the response of the sign in request as a `HttpOnly` cookie, 
and is only included in requests for obtaining a new access token.


### Auth API endpoints
The following API endpoints are included for authentication, for usage examples see the included [scratch file](docs/api.http).
- `POST /auth/sign-up` - Sign up a new user
- `POST /auth/login` - Sign in a user
- `POST /auth/token-refresh` - Refresh the access token

## Email
Email helper is included to send emails using SMTP. To configure the email settings, set the following environment variables:
- `EMAIL_HOST` - The SMTP server host
- `EMAIL_PORT` - The SMTP server port
- `EMAIL_HOST_USER` - The SMTP server username
- `EMAIL_HOST_PASSWORD` - The SMTP server password

## User
Some basic user features are included in this template.

### User model
SQLC is used to generate model using migration files located in ["schemata" directory](schemata).
The user model is included as a sample model, and includes the following fields:
- `id` - The user's unique identifier (UUID)
- `email` - The user's email address
- `email_verified` - Whether the user's email address is verified
- `password` - The user's password hash
- `last_login` - The user's last login time
- `first_name` - The user's first name
- `last_name` - The user's last name
- `is_staff` - Whether the user is a staff member // TODO: router middleware
- `is_active` - Whether the user is active // TODO: log in and refresh service logic
- `date_joined` - The user's date of joining // TODO: sign up service logic

### User queries
SQLC is used to generate repository functions using SQL queries located in ["sqlc/queries" directory](sqlc/queries).

### User API endpoints
- `PATCH /api/user/me` - Update the current user's details
- `POST /api/user/me/change-password` - Change the current user's password
- `POST /api/user/send-email-verification` - Send email verification email
- `POST /api/user/verify-email` - Verify email

## Miscellaneous commands
```bash
migrate create -ext sql -dir schemata <migration_name> // Create a new migration
sqlc generate // Generate/update repository models and functions
```

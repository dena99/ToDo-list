# Todo App

A todo application built with Go, featuring user registration and authentication.

## Features

- User registration with secure password hashing
- Username uniqueness validation
- Input validation for user data
- RESTful API for user management
- Possibility of adding deadline for tasks
- Task expiration after deadline

## Technologies Used

- Go
- GORM (for database operations)
- bcrypt (for password hashing)
- net/http (for HTTP server and client)
- encoding/json (for JSON parsing and encoding)

## Getting Started

1. Clone the repository
2. Install dependencies
3. Set up your database configuration
4. Run the application

## API Endpoints

- POST /register: Register a new user

# MyUtilityX Backend

MyUtilityX is a versatile backend service built with Go, offering various utility features including link shortening, file management, and contact form handling. This repository contains the backend application. The corresponding front-end application can be found at [https://github.com/omarammura-dev/mux-ui-react](https://github.com/omarammura-dev/mux-ui-react).

## 🚀 Features

- **Link Shortener**: Create and manage shortened URLs
- **File Management**: Upload, retrieve, and manage files securely
- **Contact Form**: Handle and store contact form submissions
- **User Authentication**: Secure user authentication using JWT
- **Email Services**: Send emails using SendGrid integration

## 🛠️ Technologies Used

- [Go](https://golang.org/): Primary programming language
- [GraphQL](https://graphql.org/) with [gqlgen](https://github.com/99designs/gqlgen): Flexible API queries
- [MongoDB](https://www.mongodb.com/): Primary database
- [Docker](https://www.docker.com/): Containerization and deployment
- [gRPC](https://grpc.io/): Efficient file handling operations
- [Gin](https://github.com/gin-gonic/gin): Web framework for routing and middleware
- [JWT](https://jwt.io/): Secure authentication and authorization
- [SendGrid](https://sendgrid.com/): Email services
- [GitLab CI/CD](https://docs.gitlab.com/ee/ci/): Continuous integration and deployment

## 📁 Project Structure

├── graph/ # GraphQL schema and resolvers
├── models/ # Data models
├── routes/ # HTTP routing (Gin)
├── utils/ # Utility functions
├── repository/ # Data access layer
├── mailS/ # Email service implementation
├── gRPC/ # gRPC client implementation

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- MongoDB
- SendGrid API key

### Installation

1. Clone the repository:

   ```
   git clone https://github.com/omarammura-dev/myutilityx.git
   cd myutilityx
   ```

2. Set up environment variables:
   Copy the `.env.example` file to `.env` and fill in the required values:

   ```
   cp .env.example .env
   ```

3. Build and run the Docker containers:

   ```
   docker-compose up --build
   ```

4. The server should now be running at `http://localhost:8080`

## 🔧 Configuration

The following environment variables are required:

- `MONGO_URL`: MongoDB connection string
- `MONGO_DB_NAME`: Name of the MongoDB database
- `SENDGRID_API_KEY`: API key for SendGrid email service
- `API_URL`: Base URL for the API
- `UI_URL`: Base URL for the user interface
- `JWT_SECRET_KEY`: Secret key for JWT token generation and verification

## 📖 API Documentation

The API is built using both GraphQL and REST.

### GraphQL API
The GraphQL schema is defined in `graph/schema.graphqls`. For detailed GraphQL API documentation, run the server and visit the GraphQL playground at `http://localhost:8080/playground`.

### REST API
The REST endpoints are defined in the `routes` package. For detailed REST API documentation,

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 📞 Contact

Your Name - [@omarammurame](https://twitter.com/omarammurame) - contact@omarammura.me

Project Link: [https://github.com/omarammura-dev/my-utility-x-Backend](https://github.com/omarammura-dev/my-utility-x-Backend)

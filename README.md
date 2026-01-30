# APIHub - API Management Platform

A comprehensive API management platform similar to Apifox/Postman, built with modern technologies.

## Tech Stack

### Frontend
- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **shadcn/ui** - UI components
- **TailwindCSS** - Styling
- **Zustand** - State management
- **React Router** - Routing
- **Axios** - HTTP client

### Backend
- **Go 1.21** - Backend framework
- **Gin** - HTTP router
- **PGX** - PostgreSQL driver
- **JWT** - Authentication

### Database
- **Supabase** - PostgreSQL database with additional features
- **Docker Compose** - Local development

## Features

- 🔐 **User Authentication** - JWT-based authentication system
- 📁 **Project Management** - Create and manage API projects
- 🗂️ **Collection Organization** - Group APIs into collections
- 📝 **API Endpoint Management** - Define and document REST APIs
- 🧪 **API Testing** - Send HTTP requests and view responses
- 🔧 **Environment Variables** - Manage multiple environments with variables
- 🎨 **Modern UI** - Clean and intuitive interface with shadcn/ui

## Project Structure

```
.
├── backend/                 # Go backend
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── config/        # Configuration
│   │   ├── handler/       # HTTP handlers
│   │   ├── middleware/    # Middleware
│   │   └── model/         # Data models
│   ├── pkg/               # Public packages
│   │   ├── auth/          # Authentication utilities
│   │   ├── database/      # Database connection
│   │   ├── logger/        # Logging
│   │   └── response/      # Response helpers
│   ├── migrations/        # Database migrations
│   └── Dockerfile
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── lib/          # Utilities and API client
│   │   ├── pages/        # Page components
│   │   ├── store/        # State management
│   │   └── types/        # TypeScript types
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml     # Docker compose configuration
└── README.md

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for local frontend development)
- Go 1.21+ (for local backend development)

### Quick Start with Docker

1. Clone the repository:
```bash
git clone <repository-url>
cd claude_helloworld
```

2. Start all services:
```bash
docker-compose up -d
```

3. Access the application:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Supabase Studio: http://localhost:3000

### Local Development

#### Backend Development

1. Install dependencies:
```bash
cd backend
go mod download
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Run database migrations:
```bash
psql -h localhost -U postgres -d postgres -f migrations/001_init.up.sql
```

4. Run the backend:
```bash
go run cmd/api/main.go
```

#### Frontend Development

1. Install dependencies:
```bash
cd frontend
npm install
```

2. Start the development server:
```bash
npm run dev
```

3. Build for production:
```bash
npm run build
```

## API Documentation

### Authentication

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

### Projects

#### List Projects
```http
GET /api/v1/projects
Authorization: Bearer <token>
```

#### Create Project
```http
POST /api/v1/projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "My API Project",
  "description": "Project description"
}
```

### Collections

#### List Collections
```http
GET /api/v1/projects/:project_id/collections
Authorization: Bearer <token>
```

#### Create Collection
```http
POST /api/v1/projects/:project_id/collections
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "User APIs",
  "description": "User management endpoints"
}
```

### Endpoints

#### List Endpoints
```http
GET /api/v1/collections/:collection_id/endpoints
Authorization: Bearer <token>
```

#### Create Endpoint
```http
POST /api/v1/collections/:collection_id/endpoints
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Get User",
  "method": "GET",
  "url": "https://api.example.com/users/:id",
  "description": "Retrieve user by ID"
}
```

### Test Request

#### Send HTTP Request
```http
POST /api/v1/test/request
Authorization: Bearer <token>
Content-Type: application/json

{
  "method": "GET",
  "url": "https://api.example.com/users",
  "headers": {
    "Authorization": "Bearer token"
  }
}
```

## Environment Variables

### Backend
- `ENVIRONMENT` - development|production
- `SERVER_PORT` - Server port (default: 8080)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - JWT signing secret
- `FRONTEND_URL` - Frontend URL for CORS

### Frontend
- `VITE_API_URL` - Backend API URL

## Roadmap

- [ ] API request history
- [ ] Automated API documentation generation
- [ ] Mock server
- [ ] Import/Export (Postman, OpenAPI)
- [ ] Team collaboration features
- [ ] API monitoring and analytics
- [ ] Webhook testing
- [ ] GraphQL support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Acknowledgments

- Built with [shadcn/ui](https://ui.shadcn.com/)
- Inspired by [Apifox](https://apifox.com/) and [Postman](https://www.postman.com/)

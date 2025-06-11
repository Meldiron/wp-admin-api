# WordPress Admin API

A Go-based web application for managing WordPress development environments with debug mode toggle functionality. Built with Fiber framework and HTMX for a modern, interactive web interface.

[![Go Lint](https://github.com/meldiron/wp-admin-api/actions/workflows/lint.yml/badge.svg)](https://github.com/meldiron/wp-admin-api/actions/workflows/lint.yml)
[![Tests & Build](https://github.com/meldiron/wp-admin-api/actions/workflows/tests.yml/badge.svg)](https://github.com/meldiron/wp-admin-api/actions/workflows/tests.yml)
[![Docker Release](https://github.com/meldiron/wp-admin-api/actions/workflows/release.yml/badge.svg)](https://github.com/meldiron/wp-admin-api/actions/workflows/release.yml)

## Features

- 🔐 **Authentication System**: Secure login with session management
- 🐛 **Debug Mode Toggle**: Enable/disable WordPress debug mode across multiple servers
- 🖥️ **Web Dashboard**: Clean, responsive interface built with HTMX and Tailwind CSS
- 🔒 **Session Management**: SQLite-based session storage
- 🐳 **Docker Support**: Full containerization with Docker Compose and multi-arch releases
- 🚀 **High Performance**: Built with Go Fiber framework for speed and efficiency
- 📦 **Multi-Architecture**: Docker images built for AMD64 and ARM64 platforms

## Architecture

The application is structured as a modern web application with the following components:

- **Backend**: Go with Fiber framework for high-performance HTTP handling
- **Frontend**: Server-side rendered HTML with HTMX for dynamic interactions
- **Database**: SQLite for session storage
- **Styling**: Tailwind CSS with Launch.css UI framework
- **Authentication**: Form-based authentication with session cookies

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose (for containerized deployment)
- SQLite (included in the build)

## Installation

### Option 1: Docker Hub Image (Recommended)

1. Create environment file:
```bash
curl -o .env https://raw.githubusercontent.com/meldiron/wp-admin-api/main/.env.example
```

2. Edit `.env` with your configuration:
```env
USERS=user:password,admin=admin
SERVERS=First app:./mock/app1,Second app:./mock/app2
```

3. Run with Docker:
```bash
docker run -d \
  --name wp-admin-api \
  -p 3000:3000 \
  --env-file .env \
  -v $(pwd)/mock:/app/mock \
  meldiron/wp-admin-api:latest
```

### Option 2: Docker Compose (Development)

1. Clone the repository:
```bash
git clone https://github.com/meldiron/wp-admin-api.git
cd wp-admin-api
```

2. Copy the environment file and configure:
```bash
cp .env.example .env
```

3. Edit `.env` with your configuration:
```env
USERS=user:password,admin=admin
SERVERS=First app:./mock/app1,Second app:./mock/app2
```

4. Run with Docker Compose:
```bash
docker-compose up -d
```

### Option 3: Local Development

1. Clone and navigate to the project:
```bash
git clone https://github.com/meldiron/wp-admin-api.git
cd wp-admin-api
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run ./src/server.go
```

The application will be available at `http://localhost:3000`.

## Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `USERS` | Comma-separated list of username:password pairs | `admin:secret,user:password` |
| `SERVERS` | Comma-separated list of server_name:path pairs | `Production:/var/www/wp1,Staging:/var/www/wp2` |

### Server Configuration

Each server entry in the `SERVERS` environment variable should point to a directory containing a WordPress installation with a `docker-compose.yml` file. The application will:

1. Read the `docker-compose.yml` file to check current debug status
2. Toggle `WORDPRESS_DEBUG` and `WORDPRESS_DEBUG_LOG` values between `true` and `false`
3. Update the file to persist the changes

## API Documentation

### Authentication Endpoints

#### Create Session (Login)
```http
POST /v1/sessions
Content-Type: application/x-www-form-urlencoded

username=admin&password=secret
```

#### Delete Session (Logout)
```http
DELETE /v1/sessions
```

### Debug Management Endpoints

#### Toggle Debug Mode
```http
POST /v1/actions/debug
Content-Type: application/x-www-form-urlencoded

server=MyWebApp2
```

### Health Check Endpoint

#### Health Check
```http
GET /health
```

Returns application health status and timestamp:
```json
{
  "status": "healthy",
  "time": "2024-01-01T12:00:00Z",
  "service": "wp-admin-api"
}
```

If unhealthy (e.g., session store unavailable):
```json
{
  "status": "unhealthy",
  "error": "session store unavailable",
  "time": "2024-01-01T12:00:00Z"
}
```

## Development

### Project Structure

```
wp-admin-api/
├── src/
│   ├── config/
│   │   └── users.go          # User authentication logic
│   ├── resources/
│   │   └── wordpress.go      # WordPress server management
│   ├── index.html           # Main HTML template
│   └── server.go            # Main application entry point
├── public/                  # Static assets (CSS, JS, images)
├── mock/                    # Mock WordPress installations for testing
├── .github/workflows/       # GitHub Actions CI/CD
├── Dockerfile              # Container definition
├── docker-compose.yml      # Development environment
├── go.mod                  # Go module definition
└── .golangci.yml          # Linting configuration
```

### Running Tests

```bash
# Run all tests
go test -cover -v ./...
```

### Linting

The project uses `golangci-lint` for comprehensive code analysis:

```bash
# Run linter locally (requires golangci-lint installation)
golangci-lint run --verbose

# Fix auto-fixable issues
golangci-lint run --fix
```

```
# Run code formatter
go fmt ./...
```

### Building

```bash
# Build Docker image
docker build -t wp-admin-api .
```

## Deployment

### Production Deployment with Docker

1. Create a production environment file:
```bash
cp .env.example .env.production
```

2. Configure production settings:
```env
USERS=admin:your_secure_password
SERVERS=WP1:/var/www/wordpress1,WP2:/var/www/wordpress2
```

3. Deploy with Docker Swarm:
```bash
docker swarm init
WP_ADMIN_API_VERSION=0.1.2 docker stack deploy -c docker-compose.prod.yml --detach=false wp-admin-api
```

### Health Monitoring

The application includes a built-in health check endpoint at `/health` that:
- Validates session store connectivity
- Returns JSON status with timestamp
- Is used by Docker Compose for container health monitoring

Docker Compose health check configuration:
- **Interval**: 30 seconds
- **Timeout**: 10 seconds  
- **Retries**: 3 attempts
- **Start Period**: 40 seconds

Monitor container health:
```bash
docker ps  # Shows health status in STATUS column
docker inspect wp-admin-api | grep -A 10 Health  # Detailed health info
```

### Security Considerations

- **Password Security**: Consider implementing password hashing (currently TODO)
- **HTTPS**: Use a reverse proxy (nginx, Traefik) for SSL termination
- **CSRF Protection**: Enable CSRF middleware for production use
- **CORS**: Configure CORS settings based on your domain requirements

## CI/CD

The project includes GitHub Actions workflows:

- **Linting** (`lint.yml`): Runs code quality checks on pull requests
- **Build & Test** (`build.yml`): Builds and tests the application, including Docker image testing
- **Docker Release** (`release.yml`): Builds and publishes multi-architecture Docker images on GitHub releases

### Docker Images

Multi-architecture Docker images are automatically built and published to Docker Hub when a new release is created:

- **Registry**: `meldiron/wp-admin-api`
- **Architectures**: `linux/amd64`, `linux/arm64`
- **Tags**: 
  - `latest` (latest release)
  - Version tags (e.g., `v1.0.0`, `1.0.0`, `1.0`, `1`)

### Creating a Release

1. Create and push a git tag:
```bash
git tag v1.0.0
git push origin v1.0.0
```

2. Create a GitHub release from the tag
3. The Docker image will be automatically built and published

### Required Secrets

For the release workflow to work, configure these GitHub repository secrets:

- `DOCKER_USERNAME`: Your Docker Hub username
- `DOCKER_PASSWORD`: Your Docker Hub password or access token

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linting:
   ```bash
   go test ./...
   golangci-lint run
   ```
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Ensure all linting checks pass
- Add tests for new functionality
- Update documentation as needed

## Troubleshooting

### Common Issues

**Application won't start:**
- Check that port 3000 is available
- Verify environment variables are set correctly

**Debug toggle not working:**
- Verify server paths in `SERVERS` environment variable
- Check that `docker-compose.yml` files exist in specified paths

**Authentication issues:**
- Verify `USERS` environment variable format
- Check that usernames don't contain colons or commas

**Health check failing:**
- Verify the application is running on port 3000
- Check that SQLite session store is accessible
- Ensure curl is available in the container
- Review Docker health check logs: `docker inspect wp-admin-api`

### Logging

1. Check container logs:
```bash
docker logs wp-admin-api
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Password hashing implementation
- [ ] CSRF protection
- [ ] Multi-user role system
- [ ] Email notifications for debug changes
- [ ] API rate limiting
- [ ] Advanced logging and monitoring

## Support

For support, please open an issue on GitHub or contact the maintainer.

---

**Built with ❤️ using Go, Fiber, and HTMX**

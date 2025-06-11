# WordPress Admin API

A Go-based web application for managing WordPress development environments with debug mode toggle functionality. Built with Fiber framework and HTMX for a modern, interactive web interface.

[![Go Lint](https://github.com/meldiron/wp-admin-api/actions/workflows/lint.yml/badge.svg)](https://github.com/meldiron/wp-admin-api/actions/workflows/lint.yml)
[![Tests & Build](https://github.com/meldiron/wp-admin-api/actions/workflows/build.yml/badge.svg)](https://github.com/meldiron/wp-admin-api/actions/workflows/build.yml)

## Features

- 🔐 **Authentication System**: Secure login with session management
- 🐛 **Debug Mode Toggle**: Enable/disable WordPress debug mode across multiple servers
- 🖥️ **Web Dashboard**: Clean, responsive interface built with HTMX and Tailwind CSS
- 🔒 **Session Management**: SQLite-based session storage
- 🐳 **Docker Support**: Full containerization with Docker Compose
- 🚀 **High Performance**: Built with Go Fiber framework for speed and efficiency

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

### Option 1: Docker (Recommended)

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

### Option 2: Local Development

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

3. Deploy with Docker Compose:
```bash
docker-compose -f docker-compose.yml --env-file .env.production up -d
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
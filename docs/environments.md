# Environment Configuration

This application supports multiple environments using environment-specific configuration files.

## Environments

- **Development** (`development`): Local development with relaxed security
- **Staging** (`staging`): Pre-production testing environment
- **Production** (`production`): Live production environment with strict security

## Setup

### 1. Choose Environment

Set the `APP_ENV` environment variable:

```bash
# Development (default)
export APP_ENV=development

# Staging
export APP_ENV=staging

# Production
export APP_ENV=production
```

### 2. Create Environment File

Copy the example file and create environment-specific configuration:

```bash
# For development
cp .env.example .env.local

# For staging
cp .env.example .env.staging

# For production
cp .env.example .env.production
```

### 3. Configure Each Environment

Edit the appropriate `.env.*` file with your environment-specific settings.

**Important Security Notes:**

- The `CSRF_KEY` MUST be exactly 32 bytes (characters)
- Generate unique CSRF keys for each environment
- Never commit actual `.env.*` files to version control
- In production, consider using a secrets manager instead of .env files

### Generating Secure CSRF Keys

Generate a secure 32-byte key using:

```bash
# Using openssl
openssl rand -base64 24 | head -c 32

# Using python
python3 -c "import secrets; print(secrets.token_urlsafe(24)[:32])"

# Or any random 32-character string
```

## Environment-Specific Behavior

### Development

- CSRF secure cookie: **disabled** (allows HTTP)
- Session secure cookie: **disabled** (allows HTTP)
- Database SSL: can be disabled
- Trusted origins: localhost, 127.0.0.1
- Validation: relaxed

### Staging

- CSRF secure cookie: **enabled** (requires HTTPS)
- Session secure cookie: **enabled** (requires HTTPS)
- Database SSL: should be enabled
- Trusted origins: from `STAGING_DOMAIN` env var
- Validation: strict

### Production

- CSRF secure cookie: **enabled** (requires HTTPS)
- Session secure cookie: **enabled** (requires HTTPS)
- Database SSL: **required**
- Trusted origins: from `PRODUCTION_DOMAIN` env var
- Validation: **very strict** - app won't start if security requirements aren't met

## Running the Application

```bash
# Development (default)
APP_ENV=development go run ./cmd/web

# Staging
APP_ENV=staging go run ./cmd/web

# Production
APP_ENV=production go run ./cmd/web
```

## Configuration Reference

### Server Configuration

- `SERVER_PORT`: Port to listen on (default: 8080)
- `SERVER_READ_TIMEOUT`: Request read timeout in seconds (default: 5)
- `SERVER_WRITE_TIMEOUT`: Response write timeout in seconds (default: 10)
- `SERVER_IDLE_TIMEOUT`: Idle connection timeout in seconds (default: 120)

### Database Configuration

- `DB_HOST`: Database host
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSLMODE`: SSL mode (disable, require, verify-ca, verify-full)

### Session Configuration

- `SESSION_COOKIE_NAME`: Name of session cookie (default: firecrest_session)
- `SESSION_LIFETIME_HRS`: Session lifetime in hours (default: 12)

### CSRF Configuration

- `CSRF_KEY`: 32-byte secret key for CSRF token generation (required)

### Domain Configuration

- `STAGING_DOMAIN`: Your staging domain (e.g., staging.example.com)
- `PRODUCTION_DOMAIN`: Your production domain (e.g., example.com)

## Docker/Container Deployment

For containerized deployments, pass environment variables:

```bash
docker run -e APP_ENV=production \
  -e CSRF_KEY="your-32-byte-key" \
  -e DB_HOST="prod-db" \
  -e PRODUCTION_DOMAIN="example.com" \
  your-app:latest
```

## Kubernetes Secrets

For Kubernetes, use secrets for sensitive values:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: firecrest-secrets
type: Opaque
stringData:
  csrf-key: "your-32-byte-secret-key-here!!!"
  db-password: "your-db-password"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: firecrest-config
data:
  APP_ENV: "production"
  SERVER_PORT: "8080"
  PRODUCTION_DOMAIN: "example.com"
```

## Troubleshooting

### "CSRF_KEY must be exactly 32 bytes"

- Ensure your CSRF_KEY is exactly 32 characters long
- Check for extra whitespace or newlines

### "CSRF secure cookie must be enabled in production"

- This is a safety check - production requires HTTPS
- Ensure `APP_ENV=production` is set correctly

### "database SSL should be enabled in production"

- Set `DB_SSLMODE=require` or higher in production
- This is a security requirement for production databases

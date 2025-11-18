### Quick Start

**Start the PostgreSQL container:**

```bash
docker-compose up -d
```

**Check if the container is running:**

```bash
docker-compose ps
```

**View logs:**

```bash
docker-compose logs -f postgres
```

**Stop the container:**

```bash
docker-compose down
```

**Stop and remove data:**

```bash
docker-compose down -v
```

### Database Connection

The database is accessible at `localhost:5432` with the following credentials:

- **Username:** `postgres`
- **Password:** `postgres`
- **Database:** `firecrest`

### Loading the Schema

**Load your schema:**

```bash
docker-compose exec -T postgres psql -U postgres -d firecrest < schema.sql
```

### Connecting to PostgreSQL

**Using psql from your host (if installed):**

```bash
psql -h localhost -U postgres -d firecrest
```

**Or exec into the container:**

```bash
docker-compose exec postgres psql -U postgres -d firecrest
```

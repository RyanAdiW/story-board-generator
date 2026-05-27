# Story Board Generator Backend (Milestone 2)

This repository now includes Milestone 1 and Milestone 2 backend basics:

- Echo HTTP server
- `POST /api/v1/storyboards` endpoint
- multipart form validation
- local product image upload
- local metadata persistence
- RabbitMQ-backed async job queue
- worker service for background processing
- job polling endpoint with real status transitions

## Architecture Layout

The project now follows this layered structure:

- `cmd/api`, `cmd/worker`
- `internal/app`
- `internal/domain`
- `internal/ports`
- `internal/adapters/http`
- `internal/adapters/postgres`
- `internal/adapters/redis`
- `internal/adapters/storage`
- `internal/adapters/ai`
- `internal/adapters/renderer`
- `internal/worker`

## API Endpoints

- `GET /health`
- `POST /api/v1/storyboards`
- `GET /api/v1/storyboards/{project_id}`
- `GET /api/v1/storyboards/{project_id}/jobs/{job_id}`

## Create Storyboard Request

`POST /api/v1/storyboards` with `multipart/form-data`

Required fields:

- `title`
- `style`
- `platform`
- `format`
- `total_duration_seconds`
- `product_images[]` (1-10 files, png/jpg/jpeg/webp)

Success response:

```json
{
  "project_id": "string",
  "job_id": "string",
  "status": "pending"
}
```

## Environment Variables

Environment configuration is now in `.env`:

- `APP_PORT`
- `UPLOAD_DIR`
- `DATA_DIR`
- `RABBITMQ_URL`
- `RABBITMQ_QUEUE`

## Run

```bash
make tidy
make run-api
```

Run worker in a separate terminal:

```bash
make run-worker
```

`make` expects `go` to be available in your shell PATH.

## Docker Compose

```bash
docker compose up --build
```

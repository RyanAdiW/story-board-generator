# Story Board Generator Backend (Milestone 5)

This repository now includes Milestone 1, Milestone 2, Milestone 3, Milestone 4, and Milestone 5 backend basics:

- Echo HTTP server
- `POST /api/v1/storyboards` endpoint
- multipart form validation
- local product image upload
- local metadata persistence
- RabbitMQ-backed async job queue
- worker service for background processing
- job polling endpoint with real status transitions
- scene metadata generation in worker (OpenAI with local fallback)
- generated scenes included in storyboard result response
- uploaded product images are attached as vision inputs in the OpenAI scene generation request (up to 4 images, max 2MB each)
- scene image generation per prompt (OpenAI image model with local fallback image)
- generated scene image assets are stored and linked back to each scene
- final storyboard layout renderer generates a single PNG output
- final storyboard output is stored as `final_storyboard_image` and returned as `final_image_url`

## Architecture Layout

The project now follows this layered structure:

- `cmd/api`, `cmd/worker`
- `internal/app`
- `internal/domain`
- `internal/ports`
- `internal/adapters/http`
- `internal/adapters/postgres`
- `internal/adapters/rabbitmq`
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
- `OPENAI_API_KEY`
- `OPENAI_TEXT_MODEL`
- `OPENAI_IMAGE_MODEL`
- `OPENAI_IMAGE_SIZE`
- `OPENAI_IMAGE_QUALITY`

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

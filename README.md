# Story Board Generator Backend (Milestone 1)

This repository now includes Milestone 1 backend basics:

- Echo HTTP server
- `POST /api/v1/storyboards` endpoint
- multipart form validation
- local product image upload
- local metadata persistence

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

## Run

```bash
make tidy
make run
```

`make` expects `go` to be available in your shell PATH.

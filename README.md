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

- `APP_PORT` (default `8080`)
- `UPLOAD_DIR` (default `./uploads`)
- `DATA_DIR` (default `./data`)

## Run

```bash
go mod tidy
go run ./cmd/api
```

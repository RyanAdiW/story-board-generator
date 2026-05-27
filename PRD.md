# PRD — AI Storyboard Image Generator

## 1. Overview

AI-powered storyboard generator that transforms product images + creative direction into a cinematic storyboard sheet for ads, TikTok/Reels content, and AI video production.

The system generates:
- storyboard scenes
- cinematic product shots
- camera directions
- motion descriptions
- sound design
- on-screen text
- visual layout poster

Final output is a single high-quality storyboard image similar to professional advertising pitch boards.

---

# 2. Goal

Build a tool that helps creators instantly generate:
- ad storyboard concepts
- TikTok/Reels commercial layouts
- AI-video-ready visual plans

without needing:
- designers
- art directors
- video editors

---

# 3. Primary Use Cases

## Use Case 1 — Product Ad Storyboard
User uploads product images and generates:
- cinematic ad storyboard
- multiple scenes
- CTA-ready layout

Example:
- RC drift car
- perfume
- sneakers
- gaming mouse

---

## Use Case 2 — AI Video Previsualization
Creators generate storyboard references before:
- Kling
- Veo
- Runway
- Sora
- Pika
- Hailuo

---

## Use Case 3 — Affiliate Marketing Content
Generate TikTok Shop / Shopee affiliate storyboard posters.

---

# 4. Input

## Required Inputs

### 4.1 Product Images
User uploads:
- 1–10 product images

Supported:
- PNG
- JPG
- WEBP

Purpose:
- maintain product consistency
- use as visual reference

---

### 4.2 Style

Examples:
- Cyberpunk
- Luxury Minimalist
- Anime
- Pixar 3D
- Cinematic
- GTA Style
- Hyperrealistic
- Documentary
- Apple Commercial
- Futuristic Neon

Can be:
- predefined presets
- custom text prompt

---

### 4.3 Total Duration

Examples:
- 15 seconds
- 30 seconds
- 45 seconds

Used to:
- calculate scene count
- pacing

### 4.4 Niche

Examples:
- Sports
- Art
- Toys

Purpose:
- maintain product consistency

---

### 4.4 Platform

Examples:
- TikTok
- Instagram Reels
- YouTube Shorts
- Shopee Video Ads

Used to:
- optimize pacing
- CTA style
- layout decisions

---

### 4.5 Format

Examples:
- 9:16
- 16:9
- 1:1

Used for:
- storyboard layout
- image generation aspect ratio

---

# 5. Output

## Final Output

A single generated storyboard image containing:

- title/header
- scene panels
- cinematic generated frames
- duration per scene
- camera movement
- sound direction
- on-screen text
- notes
- color palette
- visual tone
- CTA section

Output format:
- PNG
- JPG

Resolution target:
- 2048px+
- print-quality capable

---

# 6. Scene Structure

Each storyboard scene contains:

## Scene Number
Example:
- Scene 1
- Scene 2

---

## Time Range
Example:
- 0:00–0:03

---

## Visual
Description of scene visuals.

---

## Camera
Example:
- close up
- low angle
- tracking shot
- macro shot

---

## Motion
Example:
- drifting
- rotating
- slow push in

---

## Sound FX
Example:
- bass drop
- engine rev
- cinematic whoosh

---

## On Screen Text
Example:
- READY TO DRIFT?
- FULL CONTROL
- PREMIUM BUILD

---

## Notes
Optional production direction.

---

# 7. AI Pipeline

## Step 1 — Analyze Product Images

AI extracts:
- product shape
- dominant colors
- product category
- visual identity

---

## Step 2 — Generate Creative Direction

Based on:
- style
- platform
- duration

AI generates:
- visual theme
- lighting style
- tone
- cinematic direction

---

## Step 3 — Generate Storyboard Script

AI creates:
- scene breakdown
- durations
- shot sequencing
- CTA flow

---

## Step 4 — Generate Scene Images

For each scene:
- generate cinematic frame
- maintain product consistency
- maintain style consistency

---

## Step 5 — Compose Storyboard Layout

Auto-layout system:
- places panels
- typography
- scene labels
- cinematic blocks

Final export:
- single storyboard sheet

---

# 8. Suggested Architecture

## Backend
- Go (Echo)

---

## AI Services

### LLM
Used for:
- scene generation
- copywriting
- storyboard logic

Possible:
- OpenAI

---

### Image Generation
Used for:
- cinematic storyboard frames

Possible:
- GPT Image 

---

## Storage
- S3 / Cloudflare R2

---

# 9. Suggested Internal Workflow

```txt
Upload Product Images
        ↓
Analyze Product
        ↓
Generate Creative Direction
        ↓
Generate Storyboard Scenes
        ↓
Generate Scene Images
        ↓
Generate Layout Composition
        ↓
Export Final Storyboard Image
```

# 10. MVP Scope

## MVP Features

- upload product images
- choose style preset
- choose platform
- choose duration
- auto-generate storyboard
- export PNG

---

## Non-MVP (Later)

- editable scenes
- regenerate single scene
- video generation integration
- template marketplace
- collaborative editing
- animation preview
- direct TikTok export

---

# 11. Main Technical Challenges

## Product Consistency
Generated images must keep:
- same product
- same colors
- same proportions

---

## Layout Quality
Need professional:
- typography
- spacing
- visual hierarchy

---

## Scene Cohesion
Storyboard must feel:
- cinematic
- progressive
- not random

---

# 12. Future Vision

Positioning:
> “AI Storyboard Director”

Not just:
> “AI image generator”

The product generates:
- cinematic planning
- ad direction
- visual storytelling
- production-ready concepts

for creators, agencies, brands, and AI filmmakers.


---

# 13. Backend Architecture

## 13.1 Backend Goal

The backend is responsible for handling the full storyboard generation workflow:

- receive user input
- upload and store product images
- analyze product references
- generate storyboard metadata
- generate scene prompts
- generate scene images
- compose the final storyboard sheet
- expose generation status to frontend
- store and retrieve final results

The backend should be designed as an async job-based system because image generation and layout rendering can take time.

---

## 13.2 Recommended Backend Stack

Primary recommendation:

- Language: Go
- HTTP Framework: Echo
- Database: PostgreSQL
- Cache / Queue: Redis
- Object Storage: Cloudflare R2 / AWS S3
- Background Worker: Go worker process
- Image Processing: Go image library or external renderer service
- AI Provider: OpenAI / other image model provider

Optional supporting tools:

- Docker Compose for local development
- OpenTelemetry for tracing
- structured logging with slog
- database migration with Goose / Atlas / golang-migrate

---

## 13.3 High Level Architecture

```txt
Frontend
   |
   | REST API
   v
Backend API Service
   |
   | create generation job
   v
PostgreSQL  <------> Redis Queue
   |                 |
   |                 v
   |            Worker Service
   |                 |
   |                 v
   |        AI Services / Image Provider
   |                 |
   |                 v
   |        Object Storage
   |                 |
   v                 v
Final Storyboard Result
```

---

## 13.4 Core Backend Components

## API Service

Responsibilities:

- handle HTTP requests
- validate input
- receive uploaded product images
- create storyboard generation jobs
- return job status
- return final storyboard result
- handle authentication later if needed

The API service should not run heavy generation directly inside the request-response lifecycle.

---

## Worker Service

Responsibilities:

- consume generation jobs from Redis queue
- analyze uploaded product images
- generate storyboard scenes
- call AI image generation API
- store generated scene images
- compose final storyboard image
- update job status in database

This allows long-running tasks to run safely without blocking the API.

---

## PostgreSQL

Used for persistent data:

- users
- projects
- uploaded product images
- storyboard jobs
- storyboard scenes
- generated assets
- job status
- error logs

---

## Redis

Used for:

- background job queue
- temporary job state
- rate limiting
- retry control
- idempotency key support

---

## Object Storage

Used for storing:

- uploaded product images
- generated scene images
- final storyboard image

Recommended storage path pattern:

```txt
/storyboards/{project_id}/inputs/{image_id}.png
/storyboards/{project_id}/scenes/{scene_id}.png
/storyboards/{project_id}/outputs/final.png
```

---

# 14. Suggested Backend Folder Structure

Recommended Go backend structure:

```txt
storyboard-generator/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── app/
│   │   ├── storyboard_service.go
│   │   ├── generation_service.go
│   │   └── asset_service.go
│   │
│   ├── domain/
│   │   ├── storyboard.go
│   │   ├── scene.go
│   │   ├── asset.go
│   │   └── job.go
│   │
│   ├── ports/
│   │   ├── ai_client.go
│   │   ├── storage.go
│   │   ├── queue.go
│   │   └── repository.go
│   │
│   ├── adapters/
│   │   ├── http/
│   │   │   ├── handler.go
│   │   │   ├── routes.go
│   │   │   └── dto.go
│   │   │
│   │   ├── postgres/
│   │   │   ├── storyboard_repository.go
│   │   │   ├── job_repository.go
│   │   │   └── asset_repository.go
│   │   │
│   │   ├── redis/
│   │   │   └── queue.go
│   │   │
│   │   ├── storage/
│   │   │   └── s3_storage.go
│   │   │
│   │   ├── ai/
│   │   │   ├── openai_client.go
│   │   │   └── image_client.go
│   │   │
│   │   └── renderer/
│   │       └── storyboard_renderer.go
│   │
│   └── worker/
│       ├── processor.go
│       └── jobs.go
│
├── migrations/
├── scripts/
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

---

# 15. Backend Data Model

## 15.1 projects

Stores a user's storyboard project.

```sql
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    style TEXT NOT NULL,
    platform TEXT NOT NULL,
    format TEXT NOT NULL,
    total_duration_seconds INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## 15.2 assets

Stores uploaded and generated image assets.

```sql
CREATE TABLE assets (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id),
    asset_type TEXT NOT NULL,
    file_url TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Asset types:

- input_product_image
- generated_scene_image
- final_storyboard_image

---

## 15.3 storyboard_jobs

Stores the generation job status.

```sql
CREATE TABLE storyboard_jobs (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id),
    status TEXT NOT NULL,
    current_step TEXT,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP
);
```

Job status values:

- pending
- processing
- completed
- failed

---

## 15.4 storyboard_scenes

Stores generated storyboard scene metadata.

```sql
CREATE TABLE storyboard_scenes (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id),
    scene_number INT NOT NULL,
    start_second INT NOT NULL,
    end_second INT NOT NULL,
    visual_description TEXT NOT NULL,
    camera_direction TEXT NOT NULL,
    motion_description TEXT NOT NULL,
    sound_fx TEXT NOT NULL,
    on_screen_text TEXT NOT NULL,
    notes TEXT,
    image_prompt TEXT NOT NULL,
    image_asset_id UUID REFERENCES assets(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

# 16. Backend API Design

## 16.1 Create Storyboard Project

```http
POST /api/v1/storyboards
Content-Type: multipart/form-data
```

Form data:

```txt
title
style
platform
format
total_duration_seconds
product_images[]
```

Response:

```json
{
  "project_id": "uuid",
  "job_id": "uuid",
  "status": "pending"
}
```

---

## 16.2 Get Job Status

```http
GET /api/v1/storyboards/{project_id}/jobs/{job_id}
```

Response:

```json
{
  "job_id": "uuid",
  "project_id": "uuid",
  "status": "processing",
  "current_step": "generating_scene_images",
  "error_message": null
}
```

---

## 16.3 Get Storyboard Result

```http
GET /api/v1/storyboards/{project_id}
```

Response:

```json
{
  "project_id": "uuid",
  "title": "Storyboard Iklan RC Drift",
  "style": "Cyberpunk Neon",
  "platform": "TikTok",
  "format": "9:16",
  "total_duration_seconds": 30,
  "scenes": [],
  "final_image_url": "https://storage.example.com/storyboards/project-id/outputs/final.png"
}
```

---

## 16.4 Regenerate Storyboard

```http
POST /api/v1/storyboards/{project_id}/regenerate
```

Used to rerun generation using the same uploaded product images and configuration.

---

## 16.5 Regenerate Single Scene

```http
POST /api/v1/storyboards/{project_id}/scenes/{scene_id}/regenerate
```

This can be added after MVP.

---

# 17. Generation Job Flow

```txt
1. User submits form + product images
2. API uploads product images to object storage
3. API creates project record
4. API creates storyboard job with status = pending
5. API pushes job_id to Redis queue
6. Worker consumes job
7. Worker updates job status = processing
8. Worker analyzes product images
9. Worker generates storyboard scene metadata
10. Worker generates image prompt for each scene
11. Worker calls image generation API per scene
12. Worker uploads generated scene images to object storage
13. Worker composes final storyboard sheet
14. Worker uploads final image
15. Worker updates job status = completed
16. Frontend displays final storyboard image
```

---

# 18. Worker Processing Steps

Recommended step names:

```txt
uploading_assets
analyzing_product
generating_creative_direction
generating_scenes
generating_scene_prompts
generating_scene_images
rendering_storyboard_layout
uploading_final_output
completed
```

These step names can be stored in `storyboard_jobs.current_step`.

---

# 19. AI Prompting Strategy

## 19.1 Product Analysis Prompt

Input:

- product image URLs
- style
- platform
- format

Output:

```json
{
  "product_category": "RC drift car",
  "dominant_colors": ["silver", "blue", "black"],
  "product_features": ["remote controller", "LED headlights", "sport body", "spoiler"],
  "visual_identity": "futuristic racing toy with premium neon look"
}
```

---

## 19.2 Storyboard Scene Generation Prompt

Input:

- product analysis
- style
- platform
- total duration
- format

Output:

```json
{
  "title": "Storyboard Iklan RC Drift",
  "subtitle": "Mobil Drift Remote Super Cepat & Bisa Ngedrift",
  "scenes": [
    {
      "scene_number": 1,
      "start_second": 0,
      "end_second": 3,
      "visual_description": "Product hero shot with smoke and neon reflections.",
      "camera_direction": "Low angle close-up, slow push in.",
      "motion_description": "Headlights turn on, smoke moves slowly behind the car.",
      "sound_fx": "Engine start, bass drop.",
      "on_screen_text": "READY TO DRIFT?",
      "notes": "Hook must appear in first 2 seconds.",
      "image_prompt": "Cinematic neon product hero shot of RC drift car..."
    }
  ]
}
```

---

## 19.3 Scene Image Prompt Rules

Every scene image prompt should include:

- product consistency instruction
- selected style
- platform format
- lighting direction
- camera angle
- motion feel
- no text inside image unless intentionally needed
- high detail
- clean composition

Example:

```txt
Use the uploaded product image as the main reference. Keep the same product shape, color, decals, and proportions. Create a cinematic cyberpunk neon product advertising shot. Vertical 9:16 frame. Low angle close-up. Wet reflective floor. Blue neon rim light. Smoke in the background. Ultra realistic, premium commercial photography, shallow depth of field. No logos, no watermark, no extra text.
```

---

# 20. Storyboard Layout Renderer

The renderer combines:

- header
- storyboard title
- product subtitle
- metadata badges
- scene rows
- generated scene image panels
- scene descriptions
- color palette
- sound design block
- platform CTA block

Recommended implementation options:

## Option A — HTML to Image

Generate an HTML template and convert to PNG using:

- Playwright
- Puppeteer
- chromedp

Pros:

- easiest layout control
- supports CSS
- good typography
- easier to iterate

Recommended for MVP.

---

## Option B — Pure Go Image Composition

Use Go image libraries.

Pros:

- no browser dependency

Cons:

- slower to develop
- harder typography/layout

Not recommended for MVP.

---

# 21. Queue and Retry Strategy

Each job should support:

- max retry count: 3
- exponential backoff
- failed state with error message
- safe re-run
- idempotency key for duplicate submissions

Recommended Redis queue options:

- Asynq
- custom Redis list/stream
- Watermill

For MVP, Asynq is recommended because it already supports retries and scheduled tasks.

---

# 22. Error Handling

Common failure points:

- invalid uploaded image
- unsupported format
- image upload failure
- AI provider timeout
- AI provider quota limit
- image generation failure
- renderer failure

Each error should:

- update job status to failed
- store readable error message
- keep uploaded product images
- allow user to retry

---

# 23. Observability

Backend should include:

- request ID
- structured logs
- job ID in every worker log
- project ID in every generation log
- duration tracking per step
- OpenTelemetry tracing later

Example log fields:

```json
{
  "level": "info",
  "message": "scene image generated",
  "project_id": "uuid",
  "job_id": "uuid",
  "scene_number": 1,
  "duration_ms": 8400
}
```

---

# 24. MVP Backend Milestones

## Milestone 1 — Basic API

- create Echo server
- create project endpoint
- validate input
- upload image locally first
- store project and asset metadata

---

## Milestone 2 — Async Job

- create Redis queue
- create worker service
- process job status
- add polling endpoint

---

## Milestone 3 — AI Scene Metadata

- call LLM to generate storyboard scenes
- save scenes to database
- return scene data

---

## Milestone 4 — Scene Image Generation

- call image provider
- generate images per scene
- save generated image URLs

---

## Milestone 5 — Final Layout Rendering

- create HTML storyboard template
- render final PNG
- upload/store output
- return final image URL

---

## Milestone 6 — Polish

- retry support
- error display
- regenerate button
- simple dashboard/history

---

# 25. Environment Variables

```env
APP_ENV=local
APP_PORT=8080

DATABASE_URL=postgres://postgres:postgres@localhost:5432/storyboard?sslmode=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

STORAGE_PROVIDER=r2
S3_BUCKET=storyboard-generator
S3_REGION=auto
S3_ENDPOINT=
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=

OPENAI_API_KEY=
OPENAI_TEXT_MODEL=
OPENAI_IMAGE_MODEL=
OPENAI_IMAGE_SIZE=1024x1536
OPENAI_IMAGE_QUALITY=medium

PUBLIC_BASE_URL=http://localhost:8080
```

---

# 26. Local Development With Docker Compose

Recommended local services:

- PostgreSQL
- Redis
- API service
- Worker service

Example services:

```yaml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: storyboard
    ports:
      - "5432:5432"

  redis:
    image: redis:7
    ports:
      - "6379:6379"

  api:
    build: .
    command: go run ./cmd/api
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis

  worker:
    build: .
    command: go run ./cmd/worker
    depends_on:
      - postgres
      - redis
```

---

# 27. Backend MVP Acceptance Criteria

The backend MVP is considered done when:

- user can create a storyboard project with product images
- API returns project_id and job_id
- worker processes the job asynchronously
- job status can be polled
- scenes are generated and saved
- scene images are generated and saved
- final storyboard PNG is generated
- final image URL is returned from API
- failed jobs store clear error messages
- generation can be retried
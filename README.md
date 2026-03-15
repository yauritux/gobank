# GoBank

A simple Bank account microservice application built on top of Golang used as a simple demo on how to utilize new feature in Mux standard package Go 1.22 and above and how to deploy it to [Render](https://render.com) using various approaches.

## Running the App/Service in Local

1. `docker-compose up`

## Testing the Endpoints

1. Create / Register a new Account endpoint
```
curl -X POST -H 'Content-Type: application/json' -d '{"first_name": "M Yauri M", "last_name": "Attamimi", "account_number": "7830371235", "balance": 500000}' localhost:8080/api/accounts
```
2. Fetch all Accounts endpoint
```
curl -X GET -H 'Accept: application/json' localhost:8080/api/accounts
```
3. Healthcheck endpoint
```
curl -X GET localhost:8080/api/health

## Prerequisites (Render account + billing)

1. Create a Render account (if you haven't done so)
- Go to https://render.com and sign up (you can use GitHub login).

2. Choose a pricing plan
- Render provides multiple workspace plans (including a free tier).
- You can deploy some services on the free tier without entering billing information.

3. Add billing information (required for this project)
- This project uses a Render Postgres database.
- Render Postgres is no longer available on the free plan for new databases, so you must add billing/payment information to your Render workspace before provisioning Postgres.
- In the Render Dashboard:
  - Go to your workspace settings
  - Open Billing
  - Add a payment method

## Deploying to Render

This repository supports multiple deployment approaches on Render.
In all approaches below:
- The web service is named: `gobank-service`
- The health endpoint is: `GET /api/health`
- The API endpoints are under: `/api/*`

After your service is live, you can test:
- `https://<YOUR_SERVICE_HOST>/api/health`
- `https://<YOUR_SERVICE_HOST>/api/accounts`

### 1) Deploy from Render Dashboard (GitHub-connected, simplest)

This approach lets Render build and deploy directly from your GitHub repository.

1. Push this repository to your GitHub (if you haven't)
- Create a GitHub repo
- Push your code to it

2. Create the Postgres database on Render
- In Render Dashboard, create a new Postgres database
- Name it (recommended): `gobank-postgres`
- Choose a supported plan (example used in this repo: `Basic-1gb`)
- Note: This step requires billing info in your workspace

3. Create the web service on Render
- In Render Dashboard, create a new Web Service
- Choose "Build and deploy from a Git repository"
- Connect and authorize GitHub (if prompted)
- Select your gobank repository and the desired branch
- Runtime: Docker (because this repo provides a Dockerfile)
- Render will build using the Dockerfile at the repo root

4. Configure environment variables for the web service
Set the following environment variables (using the values from your Render Postgres database):
- `PG_HOST` (database host)
- `PG_DB` (database name)
- `PG_USER` (database user)
- `PG_PASSWORD` (database password)
- `PG_CONTAINER_PORT` (database port)
- `PG_SSL_MODE` (recommended: `require`)
Also set:
- `APP_CONTAINER_PORT` = `8080`

5. Deploy
- Trigger a deploy from the Dashboard (or enable auto-deploy on git push).

6. Initialize database schema (required)
Render managed Postgres won't run your local docker init scripts.
You must apply the schema manually once:
- Use the SQL in scripts/schema.sql
- Connect to the DB and run it with `psql` (see "Database schema initialization" below).

### 2) Deploy by pushing a prebuilt Docker image to a registry (GHCR)

This approach builds the Docker image yourself, pushes it to a registry, then Render pulls it.

1. Build the Docker image locally
From the repo root:
- Build:
  - `docker build -t ghcr.io/<YOUR_GH_USERNAME>/gobank:latest .`

2. Login to GHCR and push
- Authenticate to GitHub Container Registry (GHCR) and push:
  - `docker push ghcr.io/<YOUR_GH_USERNAME>/gobank:latest`

3. Create Postgres on Render
- Same as approach (1), create `gobank-postgres` first.

4. Create a web service that deploys from an image
- In Render Dashboard, create a new Web Service
- Choose "Deploy an existing image from a registry"
- Image URL:
  - `ghcr.io/<YOUR_GH_USERNAME>/gobank:latest`
- Configure the same environment variables as in approach (1)
- Deploy

5. Initialize database schema
- Same as approach (1), apply scripts/schema.sql once.

### 3) Deploy using `render.yaml` (Blueprint from the Render Dashboard)

This repo includes a `render.yaml` Blueprint that describes:
- a web service (`gobank-service`) built from `./Dockerfile`
- a managed Postgres database (`gobank-postgres`)

1. Ensure `render.yaml` is in your default branch
- Commit and push `render.yaml` to GitHub.

2. Create a Blueprint in Render Dashboard
- In Render Dashboard, choose to create a new Blueprint (Infrastructure as Code)
- Select your GitHub repo + branch that contains `render.yaml`
- Render will parse the YAML and show the resources it will create:
  - `gobank-service`
  - `gobank-postgres`

3. Apply the Blueprint
- Confirm and apply.
- Note: Creating the Postgres database will require billing info.

4. Initialize database schema
- Apply `scripts/schema.sql` once (see below).

## Database schema initialization (required on Render)

This project expects the `accounts` table to exist.
The schema file is:
- `scripts/schema.sql`

After provisioning your Render Postgres database, run the schema once.

One common approach from your local machine:
1. Copy the database connection string from the Render Dashboard (External Database URL).
2. Run:
- `psql "<RENDER_DATABASE_URL>" -f scripts/schema.sql`

If SSL is required and not already included in the URL, add:
- `?sslmode=require`

## Notes / Troubleshooting

- Healthcheck endpoint:
  - `GET /api/health`
- If healthcheck fails on Render:
  - Confirm your service env vars match the names used by the app:
    - `PG_HOST`, `PG_DB`, `PG_USER`, `PG_PASSWORD`, `PG_CONTAINER_PORT`, `PG_SSL_MODE`
    - `APP_CONTAINER_PORT`
- If `/api/accounts` fails but `/api/health` works:
  - You likely did not run `scripts/schema.sql` yet.
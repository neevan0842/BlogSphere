# BlogSphere

BlogSphere is a full-featured blogging platform for creating, sharing, and managing posts, comments, and user profiles. It is the successor to the BlogSpot project.

---

## Production Deployment

[https://blogsphere.mine.bz/](https://blogsphere.mine.bz/)

![BlogSphere Logo](static/blogsphere.png)

---

## Tech Stack

- **Frontend:** React, TypeScript, Vite, Tailwind CSS, Zustand, Axios, React Router, React Markdown, PrismJS, Immer, Zod, Lucide React, React Hot Toast, React Infinite Scroll, Remark GFM, ESLint
- **Backend:** Go, Chi, pgx, Prometheus, Zap, MailerSend, OAuth2, SQLC, PostgreSQL, Docker, Makefile
- **Dev Tools:** Docker Compose, Air, Migrate, SQLC, ESLint

---

## Project Structure

- **backend/** – Go API server, database, migrations, internal logic, mailer, logger, utils
- **frontend/** – React app, API clients, components, pages, store, types, utils
- **static/** – Static assets (e.g., images)
- **docker-compose.yml / docker-compose.prod.yml** – Container orchestration
- **Caddyfile, prometheus.yml, grafana-dashboard.json** – Monitoring & reverse proxy configs

---

## Installation

```sh
git clone https://github.com/neevan0842/BlogSphere.git
cd BlogSphere
```

---

## Environment Setup

- **Backend:** Copy `backend/.env` and adjust values as needed.
- **Frontend:** Copy `frontend/.env.example` to `frontend/.env` and set `VITE_API_URL` (e.g., `http://localhost:8080/api/v1`).

---

## Port Forwarding

To access remote services (Grafana, Frontend, MinIO) locally, use SSH port forwarding (replace placeholders as needed):

```sh
ssh -L <local_port1>:localhost:<remote_port1> -L <local_port2>:localhost:<remote_port2> -L <local_port3>:localhost:<remote_port3> <user>@<remote_host> -N
# Example:
# ssh -L 9090:localhost:9090 -L 3000:localhost:3000 -L 9000:localhost:9000 user@your.server.com -N
```

---

## Database Seeding

To seed the database with initial data (requires Python and dependencies, adjust environment file as needed):

```sh
uv run --with psycopg2-binary --with python-dotenv --env-file <env_file> ./scripts/seed.py
# Example:
# uv run --with psycopg2-binary --with python-dotenv --env-file .env.prod ./scripts/seed.py
```

---

## Running the Application

**Step 1 – Install dependencies**

```sh
cd backend && go mod download
cd ../frontend && npm install
```

**Step 2 – Setup environment variables**

```sh
# Backend
cp backend/.env.example backend/.env
# Frontend
cp frontend/.env.example frontend/.env
```

**Step 3 – Start backend**

```sh
cd backend
make run
```

_or with Docker Compose:_

```sh
docker compose -f docker-compose.yml --env-file ./backend/.env up -d --build
```

**Step 4 – Start frontend**

```sh
cd frontend
npm run dev
```

---

## Available Scripts

### Frontend (`frontend/package.json`)

- `npm run dev` – Start development server
- `npm run build` – Build for production
- `npm run preview` – Preview production build
- `npm run lint` – Lint code

### Backend (`backend/Makefile`)

- `make run` – Start backend (with Air)
- `make build` – Build backend binary
- `make migrate-up` – Run DB migrations
- `make migrate-down` – Rollback migrations
- `make migrate-create` – Create migration
- `make migrate-version` – Show migration version
- `make sqlc-generate` – Generate SQLC code
- `make test` – Run tests

---

## API Base URL

Set in `frontend/.env` as `VITE_API_URL` (e.g., `http://localhost:8080/api/v1`).

---

## License

MIT License

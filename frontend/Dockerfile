# --- Build frontend ---
FROM node:22-alpine AS frontend-build
WORKDIR /app

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ .
RUN npm run build

# --- Final Caddy image ---
FROM caddy:2.11

# Copy static build from frontend
COPY --from=frontend-build /app/dist /srv

# Copy your Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile
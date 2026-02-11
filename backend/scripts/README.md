# Database Seeding Script

Populates the PostgreSQL database with demo data.

## Usage

```bash
uv run --with psycopg2-binary --with python-dotenv seed.py
```

## What Gets Seeded

- 20 users
- 8 categories (Technology, Design, Business, Travel, Lifestyle, Food, Health, Education)
- 100 posts with random authors
- Post-category relationships (1-3 categories per post)
- 300 comments (70% top-level, 30% replies)
- Random post likes (0-15 per post)
- Random comment likes (0-8 per comment)
- Random user follows (0-8 per user)

## Requirements

Ensure PostgreSQL is running and `.env` file is configured in `backend/` directory with database credentials.

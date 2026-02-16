# Database Seeding

Populates the BlogSphere database with realistic demo data.

## What It Does

- Ensures exactly 8 categories exist (Technology, Design, Business, Travel, Lifestyle, Food, Health, Education)
- Generates 30-50 users with realistic profiles
- Creates 100-150 blog posts with **authentic Markdown content** (code blocks, headings, lists)
- Adds ~300-450 contextual comments
- Creates realistic engagement (likes, follows) with skewed distribution
- Timestamps spread over 2 years

## Usage

```bash
# From backend directory
uv run --with psycopg2-binary --with python-dotenv --env-file .env.prod scripts/seed.py
```

## Environment Variables

Configure in `.env` or `.env.prod`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=blogsphere
```

## Requirements

- PostgreSQL running with migrations applied
- `uv` installed ([astral.sh/uv](https://astral.sh/uv))

Install `uv`:

```bash
# Unix/macOS/Linux
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows (PowerShell)
irm https://astral.sh/uv/install.ps1 | iex
```

## Re-seeding

To start fresh:

```bash
# Truncate all tables
psql -U postgres -d blogsphere -c "TRUNCATE users, categories, posts, post_categories, comments, post_likes, comment_likes, user_follows CASCADE;"

# Re-run seeding
uv run --with psycopg2-binary --with python-dotenv scripts/seed.py
```

## Output

Generates ~3,000-7,000 records with production-quality content. Sample output:

```
✓ Connected to database: blogsphere at localhost:5432
✓ Ensured 8 categories exist
✓ Inserted 42 users
✓ Inserted 127 posts
✓ Inserted 243 post-category relationships
✓ Inserted 381 comments
✓ Inserted 3,247 post likes
✓ Inserted 1,128 comment likes
✓ Inserted 387 user follows
✓ All data committed successfully
```

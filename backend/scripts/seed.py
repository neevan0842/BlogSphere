"""
BlogSphere Database Seeding Script
Generates realistic demo data with Markdown content.

Usage:
    uv run --with psycopg2-binary --with python-dotenv seed.py

Environment:
    Configure .env or .env.prod with: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
"""

import os
import sys
import uuid
import random
from datetime import datetime, timedelta
from typing import List, Tuple, Dict
import psycopg2
from psycopg2.extras import execute_values
from dotenv import load_dotenv

load_dotenv()  # loads .env by default, override with .env.prod if needed


# ==============================================================================
# CONFIGURATION
# ==============================================================================

# Exact categories that must exist (no more, no less)
REQUIRED_CATEGORIES = [
    {
        "name": "Technology",
        "slug": "technology",
        "description": "Latest trends in software, AI, web development, and digital innovation.",
        "icon": "💻",
    },
    {
        "name": "Design",
        "slug": "design",
        "description": "UI/UX design, creative inspiration, design systems, and visual art.",
        "icon": "🎨",
    },
    {
        "name": "Business",
        "slug": "business",
        "description": "Entrepreneurship, startups, productivity tips, and business strategies.",
        "icon": "📈",
    },
    {
        "name": "Travel",
        "slug": "travel",
        "description": "Travel guides, destination reviews, and adventure stories from around the world.",
        "icon": "✈️",
    },
    {
        "name": "Lifestyle",
        "slug": "lifestyle",
        "description": "Personal development, wellness, habits, and living your best life.",
        "icon": "🌟",
    },
    {
        "name": "Food",
        "slug": "food",
        "description": "Recipes, food reviews, cooking tips, and culinary adventures.",
        "icon": "🍽️",
    },
    {
        "name": "Health",
        "slug": "health",
        "description": "Health and wellness tips and information.",
        "icon": "💪",
    },
    {
        "name": "Education",
        "slug": "education",
        "description": "Learning resources and educational content.",
        "icon": "📚",
    },
]

USER_COUNT = random.randint(30, 50)
POST_COUNT = random.randint(100, 150)
COMMENT_BASE_COUNT = POST_COUNT * 3  # ~3 comments per post average


# ==============================================================================
# DATABASE CONNECTION
# ==============================================================================


def get_db_connection():
    """
    Establish database connection using environment variables.
    Supports both DB_* and POSTGRES_* naming conventions.
    """
    # host = os.getenv("DB_HOST") or os.getenv("POSTGRES_HOST", "localhost")
    host = "localhost"
    port = os.getenv("DB_PORT") or os.getenv("POSTGRES_PORT", "5432")
    user = os.getenv("DB_USER") or os.getenv("POSTGRES_USER")
    password = os.getenv("DB_PASSWORD") or os.getenv("POSTGRES_PASSWORD")
    database = os.getenv("DB_NAME") or os.getenv("POSTGRES_DB")

    if not all([user, password, database]):
        raise ValueError("Missing required database environment variables")

    try:
        conn = psycopg2.connect(
            host=host, port=port, user=user, password=password, database=database
        )
        print(f"✓ Connected to database: {database} at {host}:{port}")
        return conn
    except psycopg2.Error as e:
        print(f"✗ Failed to connect to database: {e}")
        sys.exit(1)


# ==============================================================================
# UTILITY FUNCTIONS
# ==============================================================================


def generate_slug(text: str) -> str:
    """Generate URL-friendly slug from text."""
    slug = text.lower().replace(" ", "-").replace(":", "").replace("'", "")
    slug = "".join(c for c in slug if c.isalnum() or c == "-")
    slug = "-".join(filter(None, slug.split("-")))  # Remove consecutive dashes
    return f"{slug}-{uuid.uuid4().hex[:8]}"


def random_timestamp(start_days_ago: int = 730, end_days_ago: int = 0) -> datetime:
    """Generate random timestamp within a date range (default: last 2 years)."""
    start = datetime.now() - timedelta(days=start_days_ago)
    end = datetime.now() - timedelta(days=end_days_ago)
    delta = end - start
    return start + timedelta(seconds=random.randint(0, int(delta.total_seconds())))


def weighted_random_choice(items: List, weights: List[float]) -> any:
    """Select random item with weighted probability."""
    return random.choices(items, weights=weights, k=1)[0]


def bulk_insert(cursor, table: str, columns: List[str], data: List[Tuple]) -> int:
    """Perform bulk insert using execute_values."""
    if not data:
        return 0

    query = f"""
        INSERT INTO {table} ({', '.join(columns)})
        VALUES %s
    """
    execute_values(cursor, query, data)
    return len(data)


# ==============================================================================
# MARKDOWN CONTENT GENERATORS
# ==============================================================================

POST_TITLES = [
    # Technology & Programming
    "Implementing Rate Limiting in Go with Chi Middleware",
    "Optimizing PostgreSQL Queries for High-Traffic APIs",
    "Building a Real-Time WebSocket Server with Go",
    "Deploying Go Applications with Docker and Caddy",
    "Understanding Context Cancellation in Go",
    "Designing Clean Architecture for Maintainable APIs",
    "Scaling a Blogging Platform to 100k Users",
    "Database Indexing Strategies for Better Performance",
    "Implementing JWT Authentication in Go",
    "Monitoring Go Applications with Prometheus and Grafana",
    "Building RESTful APIs with Chi Router",
    "Understanding SOLID Principles in Go",
    "Microservices Communication Patterns",
    "Implementing CQRS in Go Applications",
    "Best Practices for Error Handling in Go",
    "Using SQLC for Type-Safe Database Queries",
    "Container Orchestration with Docker Compose",
    "Implementing Graceful Shutdown in Go Services",
    "API Versioning Strategies that Actually Work",
    "Building a CLI Tool with Cobra in Go",
    # Design & UX
    "Creating Accessible Web Applications in 2026",
    "Design Systems: Building Reusable Component Libraries",
    "Dark Mode Implementation Best Practices",
    "Responsive Design Patterns for Modern Web Apps",
    "Tailwind CSS: Utility-First Styling in Practice",
    "Typography Tips for Better Readability",
    "Color Theory for Developers",
    "Prototyping with Figma: A Developer's Guide",
    # Business & Entrepreneurship
    "Lessons Learned Building a SaaS Product from Scratch",
    "Pricing Strategies for Developer Tools",
    "Finding Your First 100 Customers",
    "Building in Public: My Journey So Far",
    "Time Management Tips for Solo Founders",
    "MVPs That Actually Validate Your Idea",
    "Bootstrapping vs Venture Capital: My Experience",
    # DevOps & Infrastructure
    "Setting Up a CI/CD Pipeline with GitHub Actions",
    "Zero-Downtime Deployments with Blue-Green Strategy",
    "Database Backup Strategies for Production Systems",
    "Securing Your API: A Comprehensive Guide",
    "Load Balancing Strategies for High Availability",
    "Infrastructure as Code with Terraform Basics",
    # Career & Learning
    "How I Landed My First Tech Job Without a Degree",
    "The Art of Code Review: Giving and Receiving Feedback",
    "Staying Up to Date in a Fast-Moving Industry",
    "Building a Strong Developer Portfolio",
    "Technical Writing for Software Engineers",
    "Mentoring Junior Developers: What I've Learned",
    # Testing & Quality
    "Testing Strategies for Go Applications",
    "Integration Testing with Docker Containers",
    "Mocking External Dependencies in Unit Tests",
    "Test-Driven Development in Practice",
    # Frontend Development
    "React Performance Optimization Techniques",
    "State Management Patterns in Modern React",
    "TypeScript Tips for Better Type Safety",
    "Building Progressive Web Apps in 2026",
    "Vite vs Webpack: A Performance Comparison",
]


def generate_post_body(title: str, topic_category: str) -> str:
    """Generate realistic Markdown content for a blog post."""

    # Go connection pool snippet for Markdown generation
    go_pool_snippet = """```go
pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
if err != nil {
    // log.Fatal(err)
    fmt.Println("Failed to connect to database")
}
defer pool.Close()
```"""

    # Code snippets for different contexts
    code_snippets = {
        "go": [
            """```go
func RateLimitMiddleware(requests int, duration time.Duration) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Every(duration/time.Duration(requests)), requests)
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```""",
            """```go
type Server struct {
    db     *sql.DB
    router *chi.Mux
    logger *slog.Logger
}

func NewServer(db *sql.DB) *Server {
    s := &Server{
        db:     db,
        router: chi.NewRouter(),
        logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
    }
    s.setupRoutes()
    return s
}
```""",
            """```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

rows, err := db.QueryContext(ctx, "SELECT * FROM posts WHERE is_published = $1", true)
if err != nil {
    return fmt.Errorf("failed to query posts: %w", err)
}
defer rows.Close()
```""",
        ],
        "sql": [
            """```sql
-- Add composite index for better query performance
CREATE INDEX idx_posts_author_published 
ON posts(author_id, is_published, created_at DESC);

-- This significantly speeds up author post listings
EXPLAIN ANALYZE SELECT * FROM posts 
WHERE author_id = $1 AND is_published = true 
ORDER BY created_at DESC;
```""",
            """```sql
-- Using a CTE for complex queries
WITH popular_posts AS (
    SELECT p.id, p.title, COUNT(pl.id) as like_count
    FROM posts p
    LEFT JOIN post_likes pl ON p.id = pl.post_id
    GROUP BY p.id
    HAVING COUNT(pl.id) > 10
)
SELECT * FROM popular_posts ORDER BY like_count DESC;
```""",
        ],
        "bash": [
            """```bash
# Build and deploy with Docker
docker build -t blogsphere:latest .
docker compose up -d

# Check logs
docker compose logs -f api
```""",
            """```bash
# Database migrations
migrate -path database/migrations -database $DATABASE_URL up

# Create a new migration
migrate create -ext sql -dir database/migrations -seq add_indexes
```""",
        ],
        "yaml": [
            """```yaml
# docker-compose.yml
services:
  api:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: ${DATABASE_URL}
    depends_on:
      - postgres
  
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: blogsphere
      POSTGRES_PASSWORD: ${DB_PASSWORD}
```""",
        ],
    }

    # Realistic content templates based on title keywords
    if "rate limit" in title.lower() or "middleware" in title.lower():
        return f"""

Rate limiting is crucial for protecting your API from abuse and ensuring fair resource distribution. In this post, I'll walk through implementing a production-ready rate limiter using Go's `golang.org/x/time/rate` package.

## Why Rate Limiting Matters

I learned this the hard way when our API got hammered by a misconfigured client making 10,000 requests per second. Our database couldn't keep up, and legitimate users suffered. That's when I knew we needed proper rate limiting.

## Implementation

Here's the middleware I built:

{random.choice(code_snippets["go"])}

The key here is using the token bucket algorithm. Each request consumes a token, and tokens regenerate at a fixed rate. It's simple but effective.

## Usage in Your Router

Using this with Chi router is straightforward:

```go
r := chi.NewRouter()
r.Use(RateLimitMiddleware(100, time.Minute)) // 100 requests per minute
r.Get("/api/posts", ListPosts)
```

## Performance Considerations

- **Memory**: Each limiter instance uses minimal memory (~200 bytes)
- **Concurrency**: The `rate.Limiter` is thread-safe out of the box
- **Accuracy**: More accurate than simple counter-based approaches

## Testing Your Rate Limiter

```bash
# Use apache bench to test
ab -n 1000 -c 10 http://localhost:8080/api/posts
```

You should see 429 responses after hitting your limit.

## What's Next

In production, you'll want per-user or per-IP limiting, which requires a more sophisticated approach using Redis. But for many applications, this simple middleware is all you need.

Have you implemented rate limiting in your APIs? What approach did you use?"""

    elif (
        "postgresql" in title.lower()
        or "database" in title.lower()
        or "query" in title.lower()
    ):
        return f"""

After scaling our blogging platform to thousands of users, we started seeing slow query times. Here's what I learned about PostgreSQL optimization.

## The Problem

Our homepage was taking 2+ seconds to load. The culprit? An unoptimized query fetching posts with their categories and authors:

```sql
SELECT p.*, u.username, array_agg(c.name) as categories
FROM posts p
JOIN users u ON p.author_id = u.id
JOIN post_categories pc ON p.id = pc.post_id
JOIN categories c ON pc.category_id = c.id
WHERE p.is_published = true
GROUP BY p.id, u.username
ORDER BY p.created_at DESC;
```

## Step 1: Add Indexes

{random.choice(code_snippets["sql"])}

This single index reduced our query time from 2000ms to 150ms. Huge win!

## Step 2: Use EXPLAIN ANALYZE

Always use `EXPLAIN ANALYZE` to understand what PostgreSQL is actually doing:

```sql
EXPLAIN ANALYZE SELECT * FROM posts WHERE author_id = 'some-uuid';
```

Look for:
- **Sequential Scans** - usually bad on large tables
- **Index Scans** - much better
- **Actual time** - not just the estimate

## Step 3: Connection Pooling

We use `pgx` pool for connection management:

{go_pool_snippet}

This alone cut our connection overhead by 80%.

## Results

- Homepage load time: **2000ms → 180ms**
- Database CPU usage: **85% → 35%**
- Concurrent users supported: **50 → 500+**

## Key Takeaways

1. Index your foreign keys
2. Use composite indexes for multi-column WHERE clauses
3. Monitor with `pg_stat_statements`
4. Connection pooling is non-negotiable

What database optimization techniques have worked for you?"""

    elif "docker" in title.lower() or "deploy" in title.lower():
        return f"""

Deploying Go applications doesn't have to be complicated. Here's my production setup using Docker and Caddy as a reverse proxy.

## The Stack

- **Go API** - compiled binary
- **PostgreSQL** - database
- **Caddy** - reverse proxy with automatic HTTPS
- **Docker Compose** - orchestration

## Multi-Stage Dockerfile

{random.choice(code_snippets["bash"])}

This produces a ~20MB image compared to ~1GB with the full Go image.

## Docker Compose Setup

{random.choice(code_snippets["yaml"])}

## Caddy Configuration

Caddy makes HTTPS trivial:

```caddyfile
yourdomain.com {{
    reverse_proxy api:8080
    encode gzip
    log
}}
```

No certificate management needed! Caddy handles Let's Encrypt automatically.

## Deployment Workflow

1. Push to main branch
2. GitHub Actions runs tests
3. Build Docker image
4. Push to registry
5. SSH to server and pull new image
6. Zero-downtime restart with `docker compose up -d`

## Monitoring

I use Prometheus + Grafana for monitoring:

```go
// Expose metrics endpoint
r.Get("/metrics", promhttp.Handler().ServeHTTP)
```

## Lessons Learned

- Always use health checks in your containers
- Volume mount your PostgreSQL data (learned this the hard way)
- Use environment variables for all configuration
- Keep your images small for faster deployments

This setup has been running in production for 8 months with zero issues. Simple, reliable, and easy to maintain.

What's your deployment strategy?"""

    elif "testing" in title.lower() or "tdd" in title.lower():
        return f"""

I used to skip tests. "They'll slow me down," I thought. Then I spent 3 hours debugging an issue that a simple unit test would have caught in 30 seconds.

## Testing Philosophy

My approach to testing Go applications:

1. **Unit tests** - test business logic in isolation
2. **Integration tests** - test database interactions
3. **E2E tests** - test critical user flows

Not everything needs 100% coverage. Focus on:
- Business logic
- Edge cases
- Public APIs

## Table-Driven Tests

Go's testing conventions make table-driven tests elegant:

```go
func TestGenerateSlug(t *testing.T) {{
    tests := []struct {{
        name  string
        input string
        want  string
    }}{{
        {{"simple", "Hello World", "hello-world"}},
        {{"special chars", "Go & PostgreSQL!", "go-postgresql"}},
        {{"unicode", "Café", "caf"}},
    }}
    
    for _, tt := range tests {{
        t.Run(tt.name, func(t *testing.T) {{
            got := GenerateSlug(tt.input)
            if !strings.Contains(got, tt.want) {{
                t.Errorf("got %v, want %v", got, tt.want)
            }}
        }})
    }}
}}
```

## Mocking Database Calls

I use `sqlc` which generates interfaces perfect for mocking:

```go
type MockQuerier struct {{
    posts []*Post
}}

func (m *MockQuerier) ListPosts(ctx context.Context) ([]*Post, error) {{
    return m.posts, nil
}}
```

## Integration Tests with Docker

Spin up a real PostgreSQL instance for integration tests:

{random.choice(code_snippets["bash"])}

## What I Test (and Don't)

**Do test:**
- Business logic functions
- Request validation
- Error handling
- Database queries

**Don't test:**
- Third-party libraries
- Generated code (like sqlc output)
- Trivial getters/setters

## Results

Since adopting TDD:
- Bugs in production: **Down 70%**
- Refactoring confidence: **Way up**
- Development speed: **Actually faster**

The upfront time investment pays off immediately.

How do you approach testing? What's your workflow?"""

    elif (
        "react" in title.lower()
        or "frontend" in title.lower()
        or "typescript" in title.lower()
    ):
        return f"""

Frontend development has evolved rapidly. Here's what I've learned building modern web applications with React and TypeScript.

## The Setup

Our frontend stack:
- **React 18** with hooks
- **TypeScript** for type safety
- **Vite** for blazing fast builds
- **Tailwind CSS** for styling

## Type-Safe API Calls

TypeScript shines when integrated with your backend types:

```typescript
interface Post {{
  id: string;
  title: string;
  body: string;
  author: User;
  createdAt: string;
}}

async function fetchPosts(): Promise<Post[]> {{
  const response = await fetch('/api/posts');
  if (!response.ok) throw new Error('Failed to fetch posts');
  return response.json();
}}
```

## Custom Hooks for Data Fetching

I created a reusable hook for API calls:

```typescript
function useApi<T>(url: string) {{
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  
  useEffect(() => {{
    fetch(url)
      .then(res => res.json())
      .then(setData)
      .catch(setError)
      .finally(() => setLoading(false));
  }}, [url]);
  
  return {{ data, loading, error }};
}}
```

Usage is clean:

```typescript
function PostList() {{
  const {{ data: posts, loading, error }} = useApi<Post[]>('/api/posts');
  
  if (loading) return <Spinner />;
  if (error) return <Error message={{error.message}} />;
  
  return <div>{{posts?.map(post => <PostCard key={{post.id}} post={{post}} />)}}</div>;
}}
```

## Performance Optimization

Key techniques that improved our app performance:

1. **Code splitting** - lazy load routes
2. **Memoization** - use `React.memo` for expensive components
3. **Virtual scrolling** - for long lists
4. **Optimistic updates** - instant UI feedback

## State Management

For global state, I keep it simple with Zustand:

```typescript
const useUserStore = create<UserState>((set) => ({{
  user: null,
  setUser: (user) => set({{ user }}),
  logout: () => set({{ user: null }}),
}}));
```

No boilerplate, just clean state management.

## Lessons Learned

- TypeScript catches bugs at compile time, not runtime
- Simple state management wins over complex Redux setups
- Tailwind CSS makes styling enjoyable
- Vite's hot reload is a game changer

What's your frontend setup? Any tools you can't live without?"""

    elif (
        "saas" in title.lower()
        or "startup" in title.lower()
        or "building" in title.lower()
    ):
        return f"""

I launched my SaaS product 6 months ago. Here's what worked, what didn't, and what I'd do differently.

## The Idea

I noticed developers struggling with API monitoring. Existing solutions were either too complex or too expensive. So I built a simple, affordable alternative.

## Month 1: Building the MVP

I gave myself 4 weeks to ship. The MVP had:
- API endpoint monitoring
- Email alerts
- Simple dashboard

That's it. No fancy features. **Ship small, iterate fast.**

## The Tech Stack

I chose boring technology:
- **Go** for the backend (reliable, fast)
- **PostgreSQL** (proven, stable)
- **React** for the frontend (I already knew it)

No microservices, no Kubernetes, no complexity. Just a monolith that works.

## Month 2: Finding Customers

I did three things:
1. Posted on Hacker News (200 signups!)
2. Wrote content on dev.to (slow but steady)
3. Direct outreach to relevant communities

The Hacker News post got me on the map. Then I had to deliver.

## Pricing Strategy

I struggled with pricing. Too low and people don't value it. Too high and nobody buys.

Final tiers:
- **Free** - 10 API endpoints
- **Pro** - $19/month, 100 endpoints
- **Business** - $49/month, unlimited

Most customers choose Pro. The free tier drives signups.

## What Actually Matters

Forget vanity metrics. Focus on:
- **Paying customers** (I have 43!)
- **Churn rate** (mine is 8%)
- **Customer feedback** (I talk to users weekly)

Revenue: **$1,247 MRR**. Not life-changing, but growing.

## Mistakes I Made

1. **Over-engineering** - I rewrote the dashboard 3 times
2. **Ignoring marketing** - I should have started earlier
3. **No email sequence** - Lost potential customers

## What Worked

- **Building in public** - Daily updates on Twitter
- **Customer interviews** - Every feature request goes through this
- **Simple pricing** - No complicated tiers

## Next 6 Months

Goals:
- Reach $5k MRR
- Add Slack integration
- Publish case studies

The journey continues. Building a SaaS is a marathon, not a sprint.

Are you building something? I'd love to hear about it!"""

    elif (
        "clean architecture" in title.lower()
        or "solid" in title.lower()
        or "design" in title.lower()
    ):
        return f"""

Clean code isn't just about formatting. It's about designing systems that are easy to understand, test, and maintain.

## The Dependency Rule

The fundamental rule of clean architecture: **dependencies point inward**. 

```
Presentation Layer → Business Logic → Data Layer
```

Your business logic should never depend on database details.

## Go Example Structure

```
backend/
  cmd/api/          # Entry point
  internal/
    domain/         # Business entities
    usecase/        # Business logic
    repository/     # Data access interfaces
    handler/        # HTTP handlers
  database/         # Implementation details
```

## Dependency Injection

{random.choice(code_snippets["go"])}

Benefits:
- Easy to test (inject mocks)
- Easy to swap implementations
- Clear dependencies

## Interface Segregation

Keep interfaces small and focused:

```go
// Good: Small, focused interfaces
type PostReader interface {{
    GetPost(id string) (*Post, error)
}}

type PostWriter interface {{
    CreatePost(post *Post) error
}}

// Bad: Too many methods
type PostRepository interface {{
    GetPost(id string) (*Post, error)
    ListPosts() ([]*Post, error)
    CreatePost(post *Post) error
    UpdatePost(post *Post) error
    DeletePost(id string) error
    SearchPosts(query string) ([]*Post, error)
}}
```

## Repository Pattern

Abstract database access behind interfaces:

```go
type PostRepository interface {{
    FindByID(ctx context.Context, id string) (*Post, error)
    Save(ctx context.Context, post *Post) error
}}

// Implementation can use any database
type PostgresPostRepository struct {{
    db *sql.DB
}}

func (r *PostgresPostRepository) FindByID(ctx context.Context, id string) (*Post, error) {{
    // PostgreSQL-specific code here
}}
```

## Testing Clean Architecture

```go
func TestCreatePost(t *testing.T) {{
    mockRepo := &MockPostRepository{{}}
    service := NewPostService(mockRepo)
    
    post := &Post{{Title: "Test Post"}}
    err := service.CreatePost(context.Background(), post)
    
    if err != nil {{
        t.Errorf("unexpected error: %v", err)
    }}
}}
```

No database needed for this test!

## When to Use This

**Use clean architecture when:**
- Your application will grow in complexity
- You have multiple data sources
- You need comprehensive testing

**Skip it when:**
- Building a simple CRUD app
- Prototyping rapidly
- Working alone on a small project

## Results

In our codebase:
- **Test coverage**: 85%
- **Time to add features**: Faster
- **Bugs introduced**: Fewer

The upfront design cost pays dividends.

What architecture patterns do you use? Let me know in the comments!"""

    else:
        # Generic template for other topics
        intros = [
            f"I've been working with {random.choice(['Go', 'PostgreSQL', 'Docker', 'React', 'TypeScript'])} for a while now, and here's what I've learned about {title.split(':')[0].lower()}.",
            f"After building several production applications, I wanted to share my experience with {title.split(':')[0].lower()}.",
            f"This topic came up in a code review recently, so I thought I'd write about it.",
        ]

        sections = [
            f"## The Problem\n\nWhen I first encountered this challenge, I tried several approaches. None of them felt quite right until I discovered this pattern.",
            f"## My Approach\n\nHere's the solution I implemented:\n\n{random.choice(code_snippets[random.choice(['go', 'sql', 'bash'])])}",
            f"## Key Insights\n\n- **Performance matters** - Always benchmark your code\n- **Simplicity wins** - Don't over-engineer\n- **Documentation helps** - Future you will thank you",
            f"## Common Pitfalls\n\nMistakes I made so you don't have to:\n\n1. Not considering edge cases\n2. Ignoring error handling\n3. Premature optimization",
        ]

        conclusions = [
            "\n## Wrapping Up\n\nThis approach has worked well in production for several months. Your mileage may vary depending on your specific use case.\n\nWhat's your experience with this? Let me know in the comments!",
            "\n## What's Next\n\nI'm planning to explore this topic further. Stay tuned for a follow-up post with advanced techniques.\n\nHave questions? Drop a comment below!",
        ]

        return (
            f"\n\n{random.choice(intros)}\n\n"
            + "\n\n".join(random.sample(sections, k=3))
            + random.choice(conclusions)
        )


# ==============================================================================
# SEEDING FUNCTIONS
# ==============================================================================


def seed_categories(cursor) -> Dict[str, str]:
    """
    Ensure exactly 8 required categories exist using UPSERT.
    Returns mapping of slug -> category_id.
    """
    print(f"\nSeeding {len(REQUIRED_CATEGORIES)} required categories...")

    category_map = {}

    for category in REQUIRED_CATEGORIES:
        category_id = str(uuid.uuid4())
        created_at = random_timestamp(900, 400)  # Categories created earlier

        cursor.execute(
            """
            INSERT INTO categories (id, name, slug, description, icon, created_at)
            VALUES (%s, %s, %s, %s, %s, %s)
            ON CONFLICT (slug) DO UPDATE SET
                name = EXCLUDED.name,
                description = EXCLUDED.description,
                icon = EXCLUDED.icon
            RETURNING id
            """,
            (
                category_id,
                category["name"],
                category["slug"],
                category["description"],
                category["icon"],
                created_at,
            ),
        )

        result = cursor.fetchone()
        category_map[category["slug"]] = result[0] if result else category_id

    print(f"✓ Ensured {len(REQUIRED_CATEGORIES)} categories exist")
    return category_map


def seed_users(cursor, count: int) -> List[Dict]:
    """Seed users and return list of user dicts with metadata."""
    print(f"\nSeeding {count} users...")

    first_names = [
        "Alex",
        "Jordan",
        "Taylor",
        "Morgan",
        "Casey",
        "Riley",
        "Avery",
        "Quinn",
        "Sam",
        "Charlie",
        "Dakota",
        "Finley",
        "Rowan",
        "Sage",
        "River",
        "Sky",
        "Parker",
        "Reese",
        "Hayden",
        "Emerson",
        "Kai",
        "Phoenix",
        "Ash",
        "Blake",
        "Cameron",
        "Drew",
        "Ellis",
        "Harper",
        "Indigo",
        "Jules",
    ]

    last_names = [
        "Smith",
        "Johnson",
        "Williams",
        "Brown",
        "Jones",
        "Garcia",
        "Miller",
        "Davis",
        "Rodriguez",
        "Martinez",
        "Chen",
        "Lee",
        "Kim",
        "Patel",
        "Singh",
        "Nguyen",
        "Cohen",
        "Murphy",
        "O'Brien",
        "Schmidt",
        "Weber",
        "Anderson",
        "Taylor",
        "Thomas",
        "Moore",
        "Jackson",
        "White",
        "Harris",
    ]

    users_data = []
    users = []

    for i in range(count):
        user_id = str(uuid.uuid4())

        first = random.choice(first_names)
        last = random.choice(last_names)
        username = f"{first.lower()}.{last.lower()}{random.randint(1, 999)}"
        email = f"{username}@example.com"
        google_id = f"google_{uuid.uuid4().hex}"
        avatar_url = f"https://i.pravatar.cc/150?u={user_id}"

        descriptions = [
            "Full-stack developer building cool stuff with Go and React.",
            "Backend engineer passionate about scalable systems.",
            "Software architect, previously at Google. Sharing what I learn.",
            "Indie hacker building SaaS products. Documenting the journey.",
            "Tech lead focused on clean code and developer experience.",
            "Developer advocate helping others learn to code.",
            "Open source contributor and occasional blogger.",
            "Engineering manager turned IC. Love mentoring devs.",
        ]
        description = random.choice(descriptions)

        # Earlier users created earlier (realistic)
        days_ago = 730 - (i * 10)  # Spread over 2 years
        created_at = random_timestamp(min(days_ago, 730), max(days_ago - 30, 1))

        # Some users are more "popular" (will get more engagement)
        popularity = weighted_random_choice(
            ["low", "medium", "high", "viral"], [0.5, 0.3, 0.15, 0.05]
        )

        users_data.append(
            (
                user_id,
                google_id,
                username,
                email,
                description,
                avatar_url,
                created_at,
                created_at,
            )
        )

        users.append(
            {
                "id": user_id,
                "username": username,
                "created_at": created_at,
                "popularity": popularity,
            }
        )

    bulk_insert(
        cursor,
        "users",
        [
            "id",
            "google_id",
            "username",
            "email",
            "description",
            "avatar_url",
            "created_at",
            "updated_at",
        ],
        users_data,
    )

    print(f"✓ Inserted {len(users_data)} users")
    return users


def seed_posts(cursor, users: List[Dict], count: int) -> List[Dict]:
    """Seed posts with realistic Markdown content."""
    print(f"\nSeeding {count} posts...")

    posts_data = []
    posts = []

    # Shuffle titles and repeat if needed
    titles = POST_TITLES * ((count // len(POST_TITLES)) + 1)
    random.shuffle(titles)

    for i in range(count):
        post_id = str(uuid.uuid4())

        # Select author (popular users post more)
        author_weights = {
            "low": 1,
            "medium": 2,
            "high": 4,
            "viral": 6,
        }
        author = weighted_random_choice(
            users, [author_weights[u["popularity"]] for u in users]
        )

        title = titles[i]
        slug = generate_slug(title)

        # Determine topic/category from title
        topic_category = "technology"
        if any(word in title.lower() for word in ["design", "ui", "ux", "tailwind"]):
            topic_category = "design"
        elif any(
            word in title.lower() for word in ["business", "saas", "startup", "price"]
        ):
            topic_category = "business"

        # Generate realistic Markdown body
        body = generate_post_body(title, topic_category)

        # 90% published
        is_published = random.random() > 0.1

        # Posts created after author, but not too recent
        author_created = author["created_at"]
        post_days_ago = (datetime.now() - author_created).days - random.randint(1, 30)
        created_at = random_timestamp(post_days_ago, 1)
        updated_at = created_at

        # Quality score affects engagement
        quality = weighted_random_choice(
            ["average", "good", "viral"], [0.7, 0.25, 0.05]
        )

        posts_data.append(
            (
                post_id,
                author["id"],
                title,
                slug,
                body,
                is_published,
                created_at,
                updated_at,
            )
        )

        posts.append(
            {
                "id": post_id,
                "author_id": author["id"],
                "title": title,
                "created_at": created_at,
                "quality": quality,
            }
        )

    bulk_insert(
        cursor,
        "posts",
        [
            "id",
            "author_id",
            "title",
            "slug",
            "body",
            "is_published",
            "created_at",
            "updated_at",
        ],
        posts_data,
    )

    print(f"✓ Inserted {len(posts_data)} posts")
    return posts


def seed_post_categories(
    cursor, posts: List[Dict], category_map: Dict[str, str]
) -> int:
    """Assign 1-3 categories to each post based on title/content."""
    print(f"\nSeeding post-category relationships...")

    post_categories_data = []
    category_ids = list(category_map.values())

    # Mapping keywords to categories
    category_keywords = {
        category_map["technology"]: [
            "go",
            "programming",
            "api",
            "database",
            "docker",
            "deploy",
            "postgresql",
            "rate limit",
            "websocket",
            "jwt",
            "monitoring",
            "sqlc",
            "microservices",
            "testing",
        ],
        category_map["design"]: [
            "design",
            "ui",
            "ux",
            "css",
            "tailwind",
            "typography",
            "color",
            "figma",
            "accessible",
        ],
        category_map["business"]: [
            "saas",
            "startup",
            "pricing",
            "building",
            "customer",
            "mvp",
            "founder",
            "business",
        ],
        category_map["education"]: [
            "learn",
            "guide",
            "tutorial",
            "tips",
            "understanding",
        ],
        category_map["lifestyle"]: ["journey", "experience", "lessons learned"],
    }

    for post in posts:
        title_lower = post["title"].lower()

        # Determine relevant categories
        matched_categories = []
        for cat_id, keywords in category_keywords.items():
            if any(keyword in title_lower for keyword in keywords):
                matched_categories.append(cat_id)

        # If no match, pick 1-2 random categories
        if not matched_categories:
            matched_categories = random.sample(category_ids, k=random.randint(1, 2))

        # Add 1-3 categories
        num_categories = min(random.randint(1, 3), len(matched_categories))
        selected = random.sample(matched_categories, k=num_categories)

        for category_id in selected:
            post_categories_data.append((post["id"], category_id))

    count = bulk_insert(
        cursor, "post_categories", ["post_id", "category_id"], post_categories_data
    )

    print(f"✓ Inserted {count} post-category relationships")
    return count


def seed_comments(
    cursor, posts: List[Dict], users: List[Dict], count: int
) -> List[Dict]:
    """Seed realistic comments that reference post content."""
    print(f"\nSeeding {count} comments...")

    comment_templates = [
        "Great writeup! The section on {} really helped clarify things for me.",
        "Thanks for sharing this. I've been struggling with {} and this gives me a good starting point.",
        "Interesting approach! Have you considered {} as an alternative?",
        "Could you elaborate more on the {} part? I'm curious about edge cases.",
        "This is exactly what I needed for my current project. The {} example was perfect.",
        "Excellent article! Bookmarked for future reference. The {} technique is brilliant.",
        "I have a different approach to {} that might be worth exploring...",
        "One thing to watch out for: {} can cause issues at scale.",
        "How does this perform with large datasets? I'm concerned about {}.",
        "Just implemented this in production. Works great! Minor note on {}:",
        "Have you benchmarked this? Curious about {} performance.",
        "Great point about {}. I made that mistake in my last project.",
        "This deserves more attention. The {} insight alone is worth it.",
        "Thanks! This helped me fix a bug related to {}.",
        "Solid advice. I'd add that {} is also important to consider.",
    ]

    reply_templates = [
        "Thanks for reading! Good question about {}.",
        "Great point! I'll add a section about that.",
        "Exactly! {} is crucial here.",
        "Thanks for the feedback! You're right about {}.",
        "I should have mentioned that. {} definitely matters.",
        "Good catch! Updated the post to clarify {}.",
    ]

    topics = [
        "error handling",
        "performance",
        "testing",
        "the implementation",
        "database design",
        "API design",
        "concurrency",
        "scalability",
        "caching",
        "monitoring",
        "deployment",
        "the architecture",
        "type safety",
        "migrations",
        "indexing",
        "rate limiting",
    ]

    comments_data = []
    comments = []

    for i in range(count):
        comment_id = str(uuid.uuid4())

        # Select post (popular posts get more comments)
        quality_weights = {"average": 1, "good": 3, "viral": 8}
        post = weighted_random_choice(
            posts, [quality_weights[p["quality"]] for p in posts]
        )

        # Select commenter (different from author usually)
        author_user_ids = [post["author_id"]]
        potential_commenters = [u for u in users if u["id"] not in author_user_ids]
        if not potential_commenters:
            potential_commenters = users

        user = random.choice(potential_commenters)

        # Author replies sometimes
        is_author_reply = random.random() < 0.15 and user["id"] != post["author_id"]
        if is_author_reply:
            # Find author
            post_author = next(u for u in users if u["id"] == post["author_id"])
            user = post_author
            body = random.choice(reply_templates).format(random.choice(topics))
        else:
            body = random.choice(comment_templates).format(random.choice(topics))

        # Comments created after post
        post_created = post["created_at"]
        days_after_post = random.randint(
            0, min(30, (datetime.now() - post_created).days)
        )
        created_at = post_created + timedelta(
            days=days_after_post, hours=random.randint(0, 23)
        )

        comments_data.append(
            (comment_id, post["id"], user["id"], body, created_at, created_at)
        )

        comments.append(
            {
                "id": comment_id,
                "post_id": post["id"],
                "created_at": created_at,
            }
        )

    bulk_insert(
        cursor,
        "comments",
        ["id", "post_id", "user_id", "body", "created_at", "updated_at"],
        comments_data,
    )

    print(f"✓ Inserted {len(comments_data)} comments")
    return comments


def seed_post_likes(cursor, posts: List[Dict], users: List[Dict]) -> int:
    """Seed post likes with realistic skewed distribution."""
    print(f"\nSeeding post likes...")

    post_likes_data = []
    used_pairs = set()

    # Viral posts get tons of likes, average posts get few
    like_ranges = {
        "average": (0, 25),
        "good": (15, 75),
        "viral": (50, 150),
    }

    for post in posts:
        min_likes, max_likes = like_ranges[post["quality"]]
        # Ensure max_likes does not exceed number of users
        max_likes = min(max_likes, len(users))
        # Ensure min_likes does not exceed max_likes
        min_likes = min(min_likes, max_likes)
        if min_likes > max_likes:
            min_likes = max_likes
        num_likes = random.randint(min_likes, max_likes)
        likers = random.sample(users, num_likes)
        for user in likers:
            pair = (post["id"], user["id"])
            if pair not in used_pairs:
                like_id = str(uuid.uuid4())
                post_likes_data.append((like_id, post["id"], user["id"]))
                used_pairs.add(pair)
    count = bulk_insert(
        cursor, "post_likes", ["id", "post_id", "user_id"], post_likes_data
    )
    print(f"✓ Inserted {count} post likes")
    return count


def seed_comment_likes(cursor, comments: List[Dict], users: List[Dict]) -> int:
    """Seed comment likes."""
    print(f"\nSeeding comment likes...")

    comment_likes_data = []
    used_pairs = set()

    for comment in comments:
        # 0-12 likes per comment
        num_likes = random.randint(0, min(12, len(users)))
        likers = random.sample(users, num_likes)

        for user in likers:
            pair = (comment["id"], user["id"])
            if pair not in used_pairs:
                like_id = str(uuid.uuid4())
                comment_likes_data.append((like_id, comment["id"], user["id"]))
                used_pairs.add(pair)

    count = bulk_insert(
        cursor, "comment_likes", ["id", "comment_id", "user_id"], comment_likes_data
    )

    print(f"✓ Inserted {count} comment likes")
    return count


def seed_user_follows(cursor, users: List[Dict]) -> int:
    """Seed user follows (popular users get more followers)."""
    print(f"\nSeeding user follows...")

    follows_data = []
    used_pairs = set()

    follow_ranges = {
        "low": (0, 5),
        "medium": (3, 15),
        "high": (10, 30),
        "viral": (20, 50),
    }

    for follower in users:
        # Each user follows some others
        num_follows = random.randint(0, min(12, len(users) - 1))

        # Tend to follow popular users
        popularity_weights = {"low": 1, "medium": 2, "high": 3, "viral": 5}
        potential_followees = [u for u in users if u["id"] != follower["id"]]

        if potential_followees:
            followees = random.choices(
                potential_followees,
                weights=[
                    popularity_weights[u["popularity"]] for u in potential_followees
                ],
                k=min(num_follows, len(potential_followees)),
            )

            for followee in followees:
                pair = (follower["id"], followee["id"])
                if pair not in used_pairs:
                    follows_data.append((follower["id"], followee["id"]))
                    used_pairs.add(pair)

    count = bulk_insert(
        cursor, "user_follows", ["follower_id", "followee_id"], follows_data
    )

    print(f"✓ Inserted {count} user follows")
    return count


def print_summary(cursor):
    """Print summary of seeded data."""
    print("\n" + "=" * 70)
    print("DATABASE SEEDING SUMMARY")
    print("=" * 70)

    tables = [
        ("users", "Users"),
        ("categories", "Categories"),
        ("posts", "Posts"),
        ("post_categories", "Post-Category Links"),
        ("comments", "Comments"),
        ("post_likes", "Post Likes"),
        ("comment_likes", "Comment Likes"),
        ("user_follows", "User Follows"),
    ]

    for table, label in tables:
        cursor.execute(f"SELECT COUNT(*) FROM {table}")
        count = cursor.fetchone()[0]
        print(f"{label:.<40} {count:>6}")

    print("=" * 70)

    # Additional stats
    cursor.execute(
        """
        SELECT 
            COUNT(DISTINCT p.id) as total_posts,
            COUNT(DISTINCT CASE WHEN p.is_published THEN p.id END) as published_posts,
            ROUND(AVG(like_count), 1) as avg_likes
        FROM posts p
        LEFT JOIN (
            SELECT post_id, COUNT(*) as like_count
            FROM post_likes
            GROUP BY post_id
        ) pl ON p.id = pl.post_id
    """
    )

    stats = cursor.fetchone()
    print(f"\nPublished Posts: {stats[1]}/{stats[0]}")
    print(f"Average Likes per Post: {stats[2]}")

    cursor.execute(
        """
        SELECT u.username, COUNT(p.id) as post_count
        FROM users u
        LEFT JOIN posts p ON u.id = p.author_id
        GROUP BY u.id, u.username
        ORDER BY post_count DESC
        LIMIT 5
    """
    )

    print("\nTop 5 Authors by Post Count:")
    for username, count in cursor.fetchall():
        print(f"  {username}: {count} posts")

    print("=" * 70)


# ==============================================================================
# MAIN
# ==============================================================================


def main():
    """Main seeding function."""
    print("=" * 70)
    print("BlogSphere Database Seeding Script")
    print("=" * 70)
    print(f"Target: {USER_COUNT} users, {POST_COUNT} posts")

    conn = None
    cursor = None

    try:
        conn = get_db_connection()
        cursor = conn.cursor()

        # Seed in correct dependency order
        category_map = seed_categories(cursor)
        users = seed_users(cursor, USER_COUNT)
        posts = seed_posts(cursor, users, POST_COUNT)
        seed_post_categories(cursor, posts, category_map)
        comments = seed_comments(cursor, posts, users, COMMENT_BASE_COUNT)
        seed_post_likes(cursor, posts, users)
        seed_comment_likes(cursor, comments, users)
        seed_user_follows(cursor, users)

        # Commit all changes
        conn.commit()
        print("\n✓ All data committed successfully")

        # Print summary
        print_summary(cursor)

    except Exception as e:
        print(f"\n✗ Error during seeding: {e}")
        if conn:
            conn.rollback()
            print("✓ Changes rolled back")
        sys.exit(1)

    finally:
        if cursor:
            cursor.close()
        if conn:
            conn.close()
            print("\n✓ Database connection closed")


if __name__ == "__main__":
    main()

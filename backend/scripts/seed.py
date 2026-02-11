"""
PostgreSQL Database Seeding Script
Populates the BlogSphere database with realistic demo data.
"""

import os
import sys
import uuid
import random
from datetime import datetime, timedelta
from typing import List, Tuple
import psycopg2
from psycopg2.extras import execute_values
from dotenv import load_dotenv

# Load environment variables
load_dotenv()


def get_db_connection():
    """
    Establish database connection using environment variables.
    Supports both DB_* and POSTGRES_* naming conventions.
    """
    host = os.getenv("DB_HOST") or os.getenv("POSTGRES_HOST", "localhost")
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


def generate_slug(text: str) -> str:
    """Generate URL-friendly slug from text."""
    slug = text.lower().replace(" ", "-")
    # Remove special characters
    slug = "".join(c for c in slug if c.isalnum() or c == "-")
    # Add random suffix to ensure uniqueness
    return f"{slug}-{uuid.uuid4().hex[:8]}"


def random_timestamp(start_days_ago: int = 365, end_days_ago: int = 0) -> datetime:
    """Generate random timestamp within a date range."""
    start = datetime.now() - timedelta(days=start_days_ago)
    end = datetime.now() - timedelta(days=end_days_ago)
    delta = end - start
    return start + timedelta(seconds=random.randint(0, int(delta.total_seconds())))


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


def seed_users(cursor, count: int = 20) -> List[str]:
    """Seed users table and return list of user IDs."""
    print(f"\nSeeding {count} users...")

    first_names = [
        "Emma",
        "Liam",
        "Olivia",
        "Noah",
        "Ava",
        "Ethan",
        "Sophia",
        "Mason",
        "Isabella",
        "William",
        "Mia",
        "James",
        "Charlotte",
        "Benjamin",
        "Amelia",
        "Lucas",
        "Harper",
        "Henry",
        "Evelyn",
        "Alexander",
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
        "Hernandez",
        "Lopez",
        "Gonzalez",
        "Wilson",
        "Anderson",
        "Thomas",
        "Taylor",
        "Moore",
        "Jackson",
        "Martin",
    ]

    users_data = []
    user_ids = []

    for i in range(count):
        user_id = str(uuid.uuid4())
        user_ids.append(user_id)

        first = random.choice(first_names)
        last = random.choice(last_names)
        username = f"{first.lower()}.{last.lower()}{random.randint(1, 999)}"
        email = f"{username}@example.com"
        google_id = f"google_{uuid.uuid4().hex}"
        avatar_url = f"https://i.pravatar.cc/150?u={user_id}"
        description = f"Tech enthusiast and {random.choice(['writer', 'blogger', 'developer', 'designer', 'creator'])}. Love sharing ideas!"
        created_at = random_timestamp(180, 1)

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

    count_inserted = bulk_insert(
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

    print(f"✓ Inserted {count_inserted} users")
    return user_ids


def seed_categories(cursor, count: int = 10) -> List[str]:
    """Seed categories table and return list of category IDs."""
    print(f"\nSeeding {count} categories...")

    category_names = [
        "Technology",
        "Design",
        "Business",
        "Travel",
        "Lifestyle",
        "Food",
        "Health",
        "Education",
    ]

    # Use all available categories (max 8)
    categories_data = []
    category_ids = []

    for name in category_names:
        category_id = str(uuid.uuid4())
        category_ids.append(category_id)
        slug = generate_slug(name)
        created_at = random_timestamp(200, 30)

        categories_data.append((category_id, name, slug, created_at))

    count_inserted = bulk_insert(
        cursor, "categories", ["id", "name", "slug", "created_at"], categories_data
    )

    print(f"✓ Inserted {count_inserted} categories")
    return category_ids


def seed_posts(cursor, user_ids: List[str], count: int = 100) -> List[str]:
    """Seed posts table and return list of post IDs."""
    print(f"\nSeeding {count} posts...")

    title_templates = [
        "Getting Started with {}",
        "10 Tips for Better {}",
        "The Ultimate Guide to {}",
        "Why {} Matters in 2026",
        "How to Master {}",
        "Understanding {} in Depth",
        "Best Practices for {}",
        "Common Mistakes in {}",
        "Advanced {} Techniques",
        "{}: A Comprehensive Overview",
    ]

    topics = [
        "React",
        "TypeScript",
        "Python",
        "Go",
        "Docker",
        "Kubernetes",
        "PostgreSQL",
        "Redis",
        "GraphQL",
        "REST APIs",
        "Microservices",
        "CI/CD",
        "Testing",
        "Performance Optimization",
        "Security",
        "Cloud Architecture",
        "System Design",
        "Algorithms",
        "Data Structures",
    ]

    posts_data = []
    post_ids = []

    for i in range(count):
        post_id = str(uuid.uuid4())
        post_ids.append(post_id)

        author_id = random.choice(user_ids)
        topic = random.choice(topics)
        title = random.choice(title_templates).format(topic)
        slug = generate_slug(title)

        # Generate realistic post body
        paragraphs = random.randint(3, 8)
        body = "\n\n".join(
            [
                f"Lorem ipsum dolor sit amet, consectetur adipiscing elit. "
                f"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. "
                f"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris. "
                f"Paragraph {j+1} of the post about {topic}."
                for j in range(paragraphs)
            ]
        )

        is_published = random.random() > 0.2  # 80% published
        created_at = random_timestamp(120, 1)
        updated_at = created_at

        posts_data.append(
            (
                post_id,
                author_id,
                title,
                slug,
                body,
                is_published,
                created_at,
                updated_at,
            )
        )

    count_inserted = bulk_insert(
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

    print(f"✓ Inserted {count_inserted} posts")
    return post_ids


def seed_post_categories(cursor, post_ids: List[str], category_ids: List[str]) -> int:
    """Seed post_categories junction table."""
    print(f"\nSeeding post-category relationships...")

    post_categories_data = []
    used_pairs = set()

    for post_id in post_ids:
        # Each post gets 1-3 categories
        num_categories = random.randint(1, min(3, len(category_ids)))
        selected_categories = random.sample(category_ids, num_categories)

        for category_id in selected_categories:
            pair = (post_id, category_id)
            if pair not in used_pairs:
                post_categories_data.append(pair)
                used_pairs.add(pair)

    count_inserted = bulk_insert(
        cursor, "post_categories", ["post_id", "category_id"], post_categories_data
    )

    print(f"✓ Inserted {count_inserted} post-category relationships")
    return count_inserted


def seed_comments(
    cursor, post_ids: List[str], user_ids: List[str], count: int = 300
) -> List[str]:
    """Seed comments table with threaded comments and return list of comment IDs."""
    print(f"\nSeeding {count} comments...")

    comment_templates = [
        "Great article! This really helped me understand {}.",
        "Thanks for sharing. I've been looking for information on {}.",
        "Interesting perspective on {}. Have you considered...?",
        "Could you elaborate more on the {} part?",
        "This is exactly what I needed for my {} project!",
        "Well written! The section about {} was particularly useful.",
        "I disagree with the point about {}. In my experience...",
        "Excellent tutorial on {}! Bookmarked for future reference.",
    ]

    comments_data = []
    comment_ids = []

    # First, create top-level comments (70% of total)
    top_level_count = int(count * 0.7)

    for i in range(top_level_count):
        comment_id = str(uuid.uuid4())
        comment_ids.append(comment_id)

        post_id = random.choice(post_ids)
        user_id = random.choice(user_ids)
        body = random.choice(comment_templates).format(
            random.choice(
                ["this topic", "the concept", "the approach", "these techniques"]
            )
        )
        created_at = random_timestamp(90, 0)

        comments_data.append(
            (comment_id, post_id, user_id, None, body, created_at, created_at)
        )

    # Insert top-level comments first
    if comments_data:
        bulk_insert(
            cursor,
            "comments",
            [
                "id",
                "post_id",
                "user_id",
                "parent_comment_id",
                "body",
                "created_at",
                "updated_at",
            ],
            comments_data,
        )

    # Now create threaded replies (30% of total)
    reply_count = count - top_level_count
    replies_data = []

    for i in range(reply_count):
        reply_id = str(uuid.uuid4())
        comment_ids.append(reply_id)

        parent_comment_id = random.choice(
            comment_ids[: len(comment_ids) // 2]
        )  # Reply to earlier comments

        # Get the post_id from parent comment
        cursor.execute(
            "SELECT post_id FROM comments WHERE id = %s", (parent_comment_id,)
        )
        result = cursor.fetchone()
        if not result:
            continue

        post_id = result[0]
        user_id = random.choice(user_ids)
        body = random.choice(
            [
                "Thanks for the comment! I appreciate your feedback.",
                "That's a great point. You're absolutely right.",
                "I see what you mean. Let me clarify...",
                "Good question! Here's my take on that...",
                "Agreed! I had the same experience.",
            ]
        )
        created_at = random_timestamp(60, 0)

        replies_data.append(
            (
                reply_id,
                post_id,
                user_id,
                parent_comment_id,
                body,
                created_at,
                created_at,
            )
        )

    # Insert replies
    if replies_data:
        bulk_insert(
            cursor,
            "comments",
            [
                "id",
                "post_id",
                "user_id",
                "parent_comment_id",
                "body",
                "created_at",
                "updated_at",
            ],
            replies_data,
        )

    total_inserted = len(comments_data) + len(replies_data)
    print(
        f"✓ Inserted {total_inserted} comments ({len(comments_data)} top-level, {len(replies_data)} replies)"
    )
    return comment_ids


def seed_post_likes(cursor, post_ids: List[str], user_ids: List[str]) -> int:
    """Seed post_likes table."""
    print(f"\nSeeding post likes...")

    post_likes_data = []
    used_pairs = set()

    # Each post gets 0-15 likes
    for post_id in post_ids:
        num_likes = random.randint(0, min(15, len(user_ids)))
        likers = random.sample(user_ids, num_likes)

        for user_id in likers:
            pair = (post_id, user_id)
            if pair not in used_pairs:
                like_id = str(uuid.uuid4())
                post_likes_data.append((like_id, post_id, user_id))
                used_pairs.add(pair)

    count_inserted = bulk_insert(
        cursor, "post_likes", ["id", "post_id", "user_id"], post_likes_data
    )

    print(f"✓ Inserted {count_inserted} post likes")
    return count_inserted


def seed_comment_likes(cursor, comment_ids: List[str], user_ids: List[str]) -> int:
    """Seed comment_likes table."""
    print(f"\nSeeding comment likes...")

    comment_likes_data = []
    used_pairs = set()

    # Each comment gets 0-8 likes
    for comment_id in comment_ids:
        num_likes = random.randint(0, min(8, len(user_ids)))
        likers = random.sample(user_ids, num_likes)

        for user_id in likers:
            pair = (comment_id, user_id)
            if pair not in used_pairs:
                like_id = str(uuid.uuid4())
                comment_likes_data.append((like_id, comment_id, user_id))
                used_pairs.add(pair)

    count_inserted = bulk_insert(
        cursor, "comment_likes", ["id", "comment_id", "user_id"], comment_likes_data
    )

    print(f"✓ Inserted {count_inserted} comment likes")
    return count_inserted


def seed_user_follows(cursor, user_ids: List[str]) -> int:
    """Seed user_follows table."""
    print(f"\nSeeding user follows...")

    follows_data = []
    used_pairs = set()

    # Each user follows 0-8 other users
    for follower_id in user_ids:
        num_follows = random.randint(0, min(8, len(user_ids) - 1))

        # Get potential followees (excluding self)
        potential_followees = [uid for uid in user_ids if uid != follower_id]
        followees = random.sample(potential_followees, num_follows)

        for followee_id in followees:
            pair = (follower_id, followee_id)
            if pair not in used_pairs:
                follows_data.append((follower_id, followee_id))
                used_pairs.add(pair)

    count_inserted = bulk_insert(
        cursor, "user_follows", ["follower_id", "followee_id"], follows_data
    )

    print(f"✓ Inserted {count_inserted} user follows")
    return count_inserted


def print_summary(cursor):
    """Print summary of seeded data."""
    print("\n" + "=" * 60)
    print("DATABASE SEEDING SUMMARY")
    print("=" * 60)

    tables = [
        "users",
        "categories",
        "posts",
        "post_categories",
        "comments",
        "post_likes",
        "comment_likes",
        "user_follows",
    ]

    for table in tables:
        cursor.execute(f"SELECT COUNT(*) FROM {table}")
        count = cursor.fetchone()[0]
        print(f"{table.ljust(20)}: {count:>6} rows")

    print("=" * 60)


def main():
    """Main seeding function."""
    print("=" * 60)
    print("BlogSphere Database Seeding Script")
    print("=" * 60)

    conn = None
    cursor = None

    try:
        # Connect to database
        conn = get_db_connection()
        cursor = conn.cursor()

        # Begin transaction
        print("\n▶ Starting transaction...")

        # Seed data in order of dependencies
        user_ids = seed_users(cursor, count=20)
        category_ids = seed_categories(cursor, count=8)
        post_ids = seed_posts(cursor, user_ids, count=100)
        seed_post_categories(cursor, post_ids, category_ids)
        comment_ids = seed_comments(cursor, post_ids, user_ids, count=300)
        seed_post_likes(cursor, post_ids, user_ids)
        seed_comment_likes(cursor, comment_ids, user_ids)
        seed_user_follows(cursor, user_ids)

        # Commit transaction
        conn.commit()
        print("\n✓ Transaction committed successfully")

        # Print summary
        print_summary(cursor)

        print("\n✓ Database seeding completed successfully!")

    except Exception as e:
        if conn:
            conn.rollback()
            print(f"\n✗ Error occurred. Transaction rolled back.")
        print(f"✗ Error: {e}")
        import traceback

        traceback.print_exc()
        sys.exit(1)

    finally:
        if cursor:
            cursor.close()
        if conn:
            conn.close()
            print("\n✓ Database connection closed")


if __name__ == "__main__":
    main()

// ============================================================================
// Database Schema Types - Matching Backend Tables
// ============================================================================

// Users table
export interface User {
  id: string; // UUID
  google_id: string;
  username: string | null;
  email: string;
  avatar_url: string | null;
  created_at: string; // TIMESTAMPTZ
  updated_at: string; // TIMESTAMPTZ
}

// Categories table
export interface Category {
  id: string; // UUID
  name: string;
  slug: string;
  created_at: string; // TIMESTAMPTZ
}

// Posts table
export interface Post {
  id: string; // UUID
  author_id: string; // UUID reference to users
  title: string;
  slug: string;
  body: string;
  is_published: boolean;
  published_at: string | null; // TIMESTAMPTZ
  created_at: string; // TIMESTAMPTZ
  updated_at: string; // TIMESTAMPTZ
}

// Comments table
export interface Comment {
  id: string; // UUID
  post_id: string; // UUID reference to posts
  user_id: string; // UUID reference to users
  parent_comment_id: string | null; // UUID reference to comments (for threading)
  body: string;
  created_at: string; // TIMESTAMPTZ
  updated_at: string; // TIMESTAMPTZ
}

// Post likes table
export interface PostLike {
  id: string; // UUID
  post_id: string; // UUID
  user_id: string; // UUID
}

// Comment likes table
export interface CommentLike {
  id: string; // UUID
  comment_id: string; // UUID
  user_id: string; // UUID
}

// User follows table
export interface UserFollow {
  follower_id: string; // UUID
  followee_id: string; // UUID
}

// ============================================================================
// Extended Types - For API Responses with Joined Data
// ============================================================================

// Post with author information
export interface PostWithAuthor extends Post {
  author: User;
}

// Post with full details for display
export interface PostWithDetails extends Post {
  author: User;
  categories: Category[];
  like_count: number;
  comment_count: number;
  user_has_liked?: boolean; // Whether current user has liked this post
}

// Comment with author information
export interface CommentWithAuthor extends Comment {
  author: User;
  like_count: number;
  user_has_liked?: boolean; // Whether current user has liked this comment
  is_own?: boolean; // Whether comment belongs to current user
}

// ============================================================================
// UI-Specific Types (for backward compatibility and UI features)
// ============================================================================

// For category display with UI metadata
export interface CategoryDisplay extends Category {
  description: string;
  color: string;
  icon: string;
}

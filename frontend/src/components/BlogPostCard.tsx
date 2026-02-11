import { Link } from "react-router";
import { Heart, MessageCircle } from "lucide-react";
import { formatRelativeDate } from "../utils/date.utils";

interface BlogPostCardProps {
  title: string;
  body: string;
  author: {
    username: string;
    avatar_url: string | null;
  };
  created_at: string;
  categories: { name: string }[];
  slug: string;
  like_count: number;
  comment_count: number;
  user_has_liked?: boolean; // Optional prop to indicate if the current user has liked the post
}

const BlogPostCard = ({
  title,
  body,
  author,
  created_at,
  categories,
  slug,
  like_count,
  comment_count,
  user_has_liked = false,
}: BlogPostCardProps) => {
  // Extract excerpt from body (first 150 characters)
  const excerpt = body.substring(0, 150) + (body.length > 150 ? "..." : "");

  return (
    <Link to={`/post/${slug}`}>
      <article className="group border border-border rounded-lg overflow-hidden hover:border-primary/50 hover:shadow-md transition-all duration-200 bg-card h-full flex flex-col">
        {/* Content */}
        <div className="p-6 space-y-4 flex flex-col h-full">
          {/* Category Tag */}
          <div>
            <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-primary/10 text-primary">
              {categories[0]?.name || "Uncategorized"}
            </span>
          </div>

          {/* Title & Excerpt */}
          <div className="space-y-2 flex-1">
            <h2 className="text-xl font-bold text-foreground group-hover:text-primary transition-colors line-clamp-2">
              {title}
            </h2>
            <p className="text-muted-foreground text-sm line-clamp-3">
              {excerpt}
            </p>
          </div>

          {/* Author & Date */}
          <div className="flex items-center gap-3 pt-2">
            <img
              src={author.avatar_url || "/placeholder.svg"}
              alt={author.username || "User"}
              className="h-8 w-8 rounded-full object-cover"
            />
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-foreground truncate">
                {author.username || "Anonymous"}
              </p>
              <p className="text-xs text-muted-foreground">
                {formatRelativeDate(created_at) || "\u00A0"}
              </p>
            </div>
          </div>

          {/* Interactions */}
          <div className="flex items-center gap-4 pt-4 border-t border-border">
            <div className="flex items-center gap-2 text-muted-foreground hover:text-primary transition-colors">
              <Heart
                className={`h-4 w-4 transition-all ${
                  user_has_liked ? "fill-primary text-primary" : ""
                }`}
              />
              <span className="text-xs font-medium">{like_count}</span>
            </div>
            <div className="flex items-center gap-2 text-muted-foreground">
              <MessageCircle className="h-4 w-4" />
              <span className="text-xs font-medium">{comment_count}</span>
            </div>
          </div>
        </div>
      </article>
    </Link>
  );
};

export default BlogPostCard;

import { Link } from "react-router";
import { Heart, MessageCircle } from "lucide-react";
import { formatRelativeDate } from "../utils/date.utils";
import useUserStore from "../store/userStore";
import toast from "react-hot-toast";
import { togglePostLikeByPostID } from "../api/postApi";
import { useState } from "react";
import type { PostType } from "../types/types";

const BlogPostCard = ({ post }: { post: PostType }) => {
  const {
    id,
    title,
    body,
    author,
    created_at,
    categories,
    slug,
    like_count,
    comment_count,
    user_has_liked,
  } = post;
  // Extract excerpt from body (first 150 characters)
  const excerpt = body.substring(0, 150) + (body.length > 150 ? "..." : "");
  const { isAuthenticated } = useUserStore();
  const [isLiked, setIsLiked] = useState(user_has_liked);
  const [likeCount, setLikeCount] = useState(like_count);

  const handleLike = async (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    if (!isAuthenticated) {
      toast.error("You must be logged in to like posts");
      return;
    }
    const result = await togglePostLikeByPostID(id);
    if (!result) {
      toast.error("Failed to toggle like.");
      return;
    }
    setIsLiked(result.liked);
    setLikeCount((prev) => prev + (result.liked ? 1 : -1));
  };

  return (
    <article className="group border border-border rounded-lg overflow-hidden hover:border-primary/50 hover:shadow-md transition-all duration-200 bg-card h-full flex flex-col">
      {/* Content */}
      <div className="p-6 flex flex-col h-full">
        <Link to={`/post/${slug}`} className="group/link block mb-4">
          {/* Category Tags */}
          {categories.length > 0 && (
            <div className="flex items-center gap-2 flex-wrap mb-4">
              {categories.map((category) => {
                return (
                  <span
                    key={category.id}
                    className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-primary/10 text-primary"
                  >
                    {category.name}
                  </span>
                );
              })}
            </div>
          )}

          {/* Title & Excerpt */}
          <div className="space-y-2">
            <h2 className="text-xl font-bold text-foreground group-hover/link:text-primary transition-colors line-clamp-2 min-h-14">
              {title}
            </h2>
            <p className="text-muted-foreground text-sm line-clamp-3">
              {excerpt}
            </p>
          </div>
        </Link>

        {/* Author & Date */}
        <Link
          to={`/u/${author.username || "unknown"}`}
          className="group/link block mb-4"
        >
          <div className="flex items-center gap-3">
            <img
              src={author.avatar_url || "/placeholder.svg"}
              alt={author.username || "User"}
              className="h-8 w-8 rounded-full object-cover"
            />
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-foreground group-hover/link:text-primary truncate">
                {author.username || "Anonymous"}
              </p>
              <p className="text-xs text-muted-foreground">
                {formatRelativeDate(created_at) || "\u00A0"}
              </p>
            </div>
          </div>
        </Link>

        {/* Interactions */}
        <div className="flex items-center gap-4 pt-4 border-t border-border mt-auto">
          <button
            onClick={handleLike}
            className="flex items-center gap-2 text-muted-foreground hover:text-primary transition-colors"
          >
            <Heart
              className={`h-4 w-4 transition-all ${
                isLiked ? "fill-primary text-primary" : ""
              }`}
            />
            <span className="text-xs font-medium">{likeCount}</span>
          </button>
          <div className="flex items-center gap-2 text-muted-foreground">
            <MessageCircle className="h-4 w-4" />
            <span className="text-xs font-medium">{comment_count}</span>
          </div>
        </div>
      </div>
    </article>
  );
};

export default BlogPostCard;

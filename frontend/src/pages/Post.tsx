import { ArrowLeft, Heart, MessageCircle, Share2 } from "lucide-react";
import PageLayout from "../components/PageLayout";
import { Link, useNavigate, useParams } from "react-router";
import useUserStore from "../store/userStore";
import React, { useEffect, useState } from "react";
import Comment from "../components/Comment";
import type { CommentWithAuthor, PostType } from "../types/types";
import {
  getCommentsByPostSlug,
  getPostBySlug,
  togglePostLikeByPostID,
} from "../api/postApi";
import toast from "react-hot-toast";
import { formatRelativeDate } from "../utils/date.utils";

const Post = () => {
  const { isAuthenticated } = useUserStore();
  const params = useParams();
  const navigate = useNavigate();
  const { slug } = params;
  const [post, setPost] = useState<PostType | null>(null);
  const [isLiked, setIsLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(0);
  const [newComment, setNewComment] = useState("");
  const [comments, setComments] = useState<CommentWithAuthor[]>([]);

  const handleLike = async (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    if (!isAuthenticated) {
      toast.error("You must be logged in to like posts");
      return;
    }
    const result = await togglePostLikeByPostID(post?.id || "");
    if (!result) {
      toast.error("Failed to toggle like. Please try again.");
      return;
    }
    setIsLiked(result.liked);
    setLikeCount((prev) => prev + (result.liked ? 1 : -1));
  };

  function handleAddComment(event: React.MouseEvent<HTMLButtonElement>): void {
    event.preventDefault();
    if (!isAuthenticated) {
      toast.error("You must be logged in to comment");
      return;
    }
  }

  useEffect(() => {
    const fetchPost = async () => {
      const [postData, commentsData] = await Promise.all([
        getPostBySlug(slug || ""),
        getCommentsByPostSlug(slug || ""),
      ]);
      if (!postData) {
        setPost(null);
        toast.error("Post not found");
        navigate("/");
        return;
      }

      if (commentsData) {
        setComments(commentsData);
      }
      setPost(postData);
      setIsLiked(postData.user_has_liked);
      setLikeCount(postData.like_count);
    };
    fetchPost();
  }, [slug, navigate]);

  if (!post) {
    return (
      <PageLayout>
        <div className="flex-1 flex items-center justify-center">
          <p className="text-foreground">Loading...</p>
        </div>
      </PageLayout>
    );
  }

  return (
    <PageLayout>
      <article className="container mx-auto px-4 py-8 md:py-12 max-w-3xl">
        {/* Back Link */}
        <Link
          to="/"
          className="inline-flex items-center gap-2 text-primary hover:text-primary/80 transition-colors mb-8"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to articles
        </Link>

        {/* Article Header */}
        <header className="mb-8 space-y-4">
          <div className="flex items-center gap-2 flex-wrap">
            {post.categories.length > 0 &&
              post.categories.map((category) => {
                return (
                  <span
                    key={category.id}
                    className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-primary/10 text-primary mb-4"
                  >
                    {category.name}
                  </span>
                );
              })}
          </div>

          <h1 className="text-4xl md:text-5xl font-bold text-balance">
            {post.title}
          </h1>

          {/* Author Info */}
          <div className="flex items-center gap-4 py-6 border-t border-b border-border">
            <img
              src={post.author.avatar_url || "/placeholder.svg"}
              alt={post.author.username || "User"}
              className="h-12 w-12 rounded-full object-cover"
            />
            <div className="flex-1">
              <p className="font-semibold text-foreground">
                {post.author.username || "Anonymous"}
              </p>
              <p className="text-sm text-muted-foreground">
                {post.created_at
                  ? formatRelativeDate(post.created_at)
                  : "Draft"}
              </p>
            </div>

            {isAuthenticated && (
              <button className="px-4 py-2 rounded-lg border border-primary text-primary hover:bg-primary/10 transition-colors font-medium">
                Follow
              </button>
            )}
          </div>
        </header>

        {/* Article Content */}
        <div className="prose prose-lg dark:prose-invert max-w-none mb-8 text-foreground">
          {post.body.split("\n\n").map((paragraph: string, idx: number) => {
            if (paragraph.startsWith("#")) {
              const level = paragraph.match(/^#+/)?.[0].length || 2;
              const text = paragraph.replace(/^#+\s*/, "");
              const HeadingTag = `h${level}` as any;
              return (
                <HeadingTag
                  key={idx}
                  className={`font-bold my-4 ${
                    level === 2 ? "text-2xl mt-6" : "text-xl mt-5"
                  }`}
                >
                  {text}
                </HeadingTag>
              );
            }
            return (
              <p
                key={idx}
                className="text-base leading-relaxed text-foreground mb-4"
              >
                {paragraph}
              </p>
            );
          })}
        </div>

        {/* Interactions */}
        <div className="border-y border-border py-6 mb-8 flex items-center gap-4">
          <button
            onClick={handleLike}
            className="flex items-center gap-2 text-muted-foreground hover:text-primary transition-colors"
          >
            <Heart
              className={`h-5 w-5 transition-all ${
                isLiked ? "fill-primary text-primary" : ""
              }`}
            />
            <span className="font-medium">{likeCount}</span>
          </button>

          <div className="flex items-center gap-2 text-muted-foreground">
            <MessageCircle className="h-5 w-5" />
            <span className="font-medium">{comments.length}</span>
          </div>

          <button className="flex items-center gap-2 text-muted-foreground hover:text-primary transition-colors">
            <Share2 className="h-5 w-5" />
            <span className="font-medium">Share</span>
          </button>
        </div>

        {/* Comments Section */}
        <section className="space-y-8">
          <h2 className="text-2xl font-bold text-foreground">Comments</h2>

          {/* Comment Form */}
          {isAuthenticated && (
            <div className="space-y-2 p-4 rounded-lg border border-border bg-card">
              <h3 className="text-sm font-semibold text-foreground">
                Add a comment
              </h3>
              <textarea
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                placeholder="Share your thoughts..."
                className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 resize-none text-sm"
                rows={2}
              />
              <div className="flex justify-end">
                <button
                  onClick={handleAddComment}
                  disabled={!newComment.trim()}
                  className="px-4 py-1.5 text-sm rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium disabled:opacity-50"
                >
                  Post
                </button>
              </div>
            </div>
          )}

          {/* Comments List */}
          {comments.length > 0 ? (
            <div className="space-y-6">
              {comments.map((comment) => (
                <Comment
                  key={comment.id}
                  body={comment.body}
                  author={comment.author}
                  created_at={comment.created_at}
                />
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <MessageCircle className="h-12 w-12 text-muted-foreground mx-auto mb-4 opacity-50" />
              <p className="text-muted-foreground">
                No comments yet. Be the first to share your thoughts!
              </p>
            </div>
          )}
        </section>
      </article>
    </PageLayout>
  );
};

export default Post;

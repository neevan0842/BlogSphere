import { Edit, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import useUserStore from "../store/userStore";
import { formatRelativeDate } from "../utils/date.utils";
import type { CommentWithAuthor } from "../types/types";
import { deleteCommentByID, updateCommentByID } from "../api/commentApi";
import toast from "react-hot-toast";
import { Link } from "react-router";

interface CommentProps {
  comment: CommentWithAuthor;
  onUpdateComment: (id: string, body: string) => void;
  onDeleteComment: (id: string) => void;
}

const Comment = ({
  comment,
  onUpdateComment,
  onDeleteComment,
}: CommentProps) => {
  const { id, body, author, created_at } = comment;
  const [isEditing, setIsEditing] = useState(false);
  const [editedContent, setEditedContent] = useState(body || "");
  const [isOwner, setIsOwner] = useState(false);
  const { isAuthenticated, user } = useUserStore();

  const handleUpdateComment = async () => {
    const data = await updateCommentByID({
      commentId: id,
      body: editedContent,
    });
    if (!data) {
      toast.error("Failed to update comment.");
      return;
    }
    onUpdateComment(data.id, data.body);
    setIsEditing(false);
    setEditedContent(data.body);
  };

  const handleDeleteComment = async () => {
    const result = await deleteCommentByID({ commentId: id });
    if (!result) {
      toast.error("Failed to delete comment.");
      return;
    }
    onDeleteComment(id);
  };

  useEffect(() => {
    setIsOwner(isAuthenticated && user?.id === author?.id);
  }, [author, user, isAuthenticated]);

  return (
    <div className="flex gap-4">
      <Link to={`/u/${author?.username || "unknown"}`} className="shrink-0">
        <img
          src={author?.avatar_url || "/placeholder.svg"}
          alt={author?.username || "User"}
          className="h-10 w-10 rounded-full object-cover shrink-0"
        />
      </Link>
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between gap-2 mb-1">
          <Link to={`/u/${author?.username || "unknown"}`} className="shrink-0">
            <h4 className="font-semibold text-foreground">
              {author?.username || "Anonymous"}
            </h4>
          </Link>
          <span className="text-xs text-muted-foreground">
            {formatRelativeDate(created_at)}
          </span>
        </div>

        {isEditing ? (
          <div className="space-y-2">
            <textarea
              value={editedContent}
              onChange={(e) => setEditedContent(e.target.value)}
              className="w-full px-3 py-2 rounded-lg border border-border bg-card text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 resize-none"
              rows={3}
            />
            <div className="flex gap-2">
              <button
                onClick={handleUpdateComment}
                className="px-3 py-1 rounded text-sm font-medium bg-primary text-primary-foreground hover:bg-primary/90"
              >
                Save
              </button>
              <button
                onClick={() => setIsEditing(false)}
                className="px-3 py-1 rounded text-sm font-medium border border-border hover:bg-secondary"
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <>
            <p className="text-foreground text-sm leading-relaxed mb-2">
              {body}
            </p>

            {isOwner && (
              <div className="flex gap-3">
                <button
                  onClick={() => setIsEditing(true)}
                  className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-primary transition-colors"
                >
                  <Edit className="h-3 w-3" />
                  Edit
                </button>
                <button
                  onClick={handleDeleteComment}
                  className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-destructive transition-colors"
                >
                  <Trash2 className="h-3 w-3" />
                  Delete
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default Comment;

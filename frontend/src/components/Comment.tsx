import { Edit, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import useUserStore from "../store/userStore";
import { formatRelativeDate } from "../utils/date.utils";
import type { User } from "../types/types";

interface CommentProps {
  body: string;
  author: User;
  created_at: string;
}

const Comment = ({ body, author, created_at }: CommentProps) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editedContent, setEditedContent] = useState(body);
  const [isOwner, setIsOwner] = useState(false);
  const { isAuthenticated, user } = useUserStore();

  // TODO handle comment update and delete

  useEffect(() => {
    setIsOwner(isAuthenticated && user?.id === author.id);
  }, [author, user, isAuthenticated]);

  return (
    <div className="flex gap-4">
      <img
        src={author.avatar_url || "/placeholder.svg"}
        alt={author.username || "User"}
        className="h-10 w-10 rounded-full object-cover shrink-0"
      />

      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between gap-2 mb-1">
          <h4 className="font-semibold text-foreground">
            {author.username || "Anonymous"}
          </h4>
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
                onClick={() => setIsEditing(false)}
                className="px-3 py-1 rounded text-sm font-medium bg-primary text-primary-foreground hover:bg-primary/90"
              >
                Save
              </button>
              <button
                onClick={() => {
                  setIsEditing(false);
                  setEditedContent(body);
                }}
                className="px-3 py-1 rounded text-sm font-medium border border-border hover:bg-secondary"
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <>
            <p className="text-foreground text-sm leading-relaxed mb-2">
              {editedContent}
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
                <button className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-destructive transition-colors">
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

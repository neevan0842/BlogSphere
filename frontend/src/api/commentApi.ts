import type { CommentWithAuthor } from "../types/types";
import { api } from "./api";

export const addCommentToPost = async ({
  postId,
  body,
}: {
  postId: string;
  body: string;
}): Promise<CommentWithAuthor | null> => {
  try {
    if (!postId.trim() || !body.trim()) {
      return null;
    }
    const response = await api.post("/comments", {
      post_id: postId,
      body,
    });
    return response.data as CommentWithAuthor;
  } catch (error) {
    return null;
  }
};

export const updateCommentByID = async ({
  commentId,
  body,
}: {
  commentId: string;
  body: string;
}): Promise<CommentWithAuthor | null> => {
  try {
    if (!commentId.trim() || !body.trim()) {
      return null;
    }
    const response = await api.patch(`/comments/${commentId}`, {
      body,
    });
    return response.data as CommentWithAuthor;
  } catch (error) {
    return null;
  }
};

export const deleteCommentByID = async ({
  commentId,
}: {
  commentId: string;
}): Promise<boolean> => {
  try {
    const response = await api.delete(`/comments/${commentId}`);
    return response.status === 204;
  } catch (error) {
    return false;
  }
};

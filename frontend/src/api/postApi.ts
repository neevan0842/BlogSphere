import type {
  CommentWithAuthor,
  PaginatedPostsResponse,
  PostLikeResponse,
  PostType,
} from "../types/types";
import { api } from "./api";

export const getPostsPaginated = async (
  search: string,
  categorySlug: string,
  page: number,
  limit: number,
): Promise<PaginatedPostsResponse> => {
  try {
    const response = await api.get("/posts", {
      params: { search, category: categorySlug, page, limit },
    });
    return response.data as PaginatedPostsResponse;
  } catch (error) {
    return {
      posts: [],
      page,
      limit,
      hasMore: false,
    };
  }
};

export const getPostBySlug = async (slug: string): Promise<PostType | null> => {
  try {
    if (!slug.trim()) {
      return null;
    }
    const response = await api.get(`/posts/${slug.trim()}`);
    return response.data as PostType;
  } catch (error) {
    return null;
  }
};

export const getCommentsByPostSlug = async (
  slug: string,
): Promise<CommentWithAuthor[] | null> => {
  try {
    if (!slug.trim()) {
      return null;
    }
    const response = await api.get(`/posts/${slug.trim()}/comments`);
    return response.data as CommentWithAuthor[];
  } catch (error) {
    return null;
  }
};

export const togglePostLikeByPostID = async (
  postID: string,
): Promise<PostLikeResponse | null> => {
  try {
    if (!postID.trim()) {
      return null;
    }
    const response = await api.post(`/posts/${postID.trim()}/likes`);
    return response.data as PostLikeResponse;
  } catch (error) {
    return null;
  }
};

export const createPost = async ({
  title = "",
  body = "",
  categoryIDs = [],
}: {
  title: string;
  body: string;
  categoryIDs: string[];
}): Promise<PostType | null> => {
  try {
    const response = await api.post("/posts", {
      title: title.trim(),
      body: body.trim(),
      category_ids: categoryIDs,
    });
    if (response.status !== 201) {
      return null;
    }
    return response.data as PostType;
  } catch (error) {
    return null;
  }
};

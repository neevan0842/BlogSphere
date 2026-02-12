import type { PaginatedPostsResponse } from "../types/types";
import { api } from "./api";

const getPostsPaginated = async (
  search: string,
  page: number,
  limit: number,
): Promise<PaginatedPostsResponse> => {
  try {
    const response = await api.get("/posts", {
      params: { search, page, limit },
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

export { getPostsPaginated };

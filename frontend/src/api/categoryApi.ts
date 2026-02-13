import type { CategoryDisplay } from "../types/types";
import { api } from "./api";

export const getCategories = async (): Promise<CategoryDisplay[] | null> => {
  try {
    const response = await api.get("/categories");
    return response.data as CategoryDisplay[];
  } catch (error) {
    return null;
  }
};

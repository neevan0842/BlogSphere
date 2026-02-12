import { api } from "./api";
import { getUserIDFromToken } from "../utils/auth.utils";
import useUserStore from "../store/userStore";
import type { PostType, User } from "../types/types";

export const getGoogleAuthURL = async (): Promise<string> => {
  try {
    const response = await api.get("/auth/google");
    return response.data.url || "";
  } catch (error) {
    return "";
  }
};

export const exchangeGoogleCodeForToken = async (
  code: string,
  state: string,
): Promise<boolean> => {
  try {
    const response = await api.get("/auth/google/callback", {
      params: { code, state },
    });

    // Store tokens in localStorage
    localStorage.setItem("access-token", response.data.access_token);
    localStorage.setItem("refresh-token", response.data.refresh_token);

    // Decode user ID from access token
    const userID = getUserIDFromToken(response.data.access_token);
    if (!userID) {
      return false;
    }

    //fetch user details
    const user = await getUserDetailsFromUserID(userID);
    if (!user) {
      return false;
    }

    // Update user store with authenticated user ID
    const { setUser } = useUserStore.getState();
    setUser(user);
    return true;
  } catch (error) {
    return false;
  }
};

export const refreshAccessToken = async (): Promise<boolean> => {
  try {
    const response = await api.post("/auth/refresh", {
      refresh_token: localStorage.getItem("refresh-token"),
    });

    // Decode user ID from new access token
    const userID = getUserIDFromToken(response.data.access_token);
    if (!userID) {
      return false;
    }

    //fetch user details
    const user = await getUserDetailsFromUserID(userID);
    if (!user) {
      return false;
    }

    // Update user store with authenticated user ID and new access token
    const { setUser } = useUserStore.getState();
    setUser(user);
    localStorage.setItem("access-token", response.data.access_token);
    return true;
  } catch (error) {
    return false;
  }
};

export const getUserDetailsFromUserID = async (
  userID: string,
): Promise<User | null> => {
  try {
    const response = await api.get(`/users/${userID}`);
    return response.data as User;
  } catch (error) {
    return null;
  }
};

export const getUserDetailsFromUsername = async (
  username: string,
): Promise<User | null> => {
  try {
    const response = await api.get(`/users/u/${username}`);
    return response.data as User;
  } catch (error) {
    return null;
  }
};

export const getUserPosts = async (username: string): Promise<PostType[]> => {
  try {
    const response = await api.get(`/users/u/${username}/posts`);
    return response.data as PostType[];
  } catch (error) {
    return [];
  }
};

export const getUserLikedPosts = async (
  username: string,
): Promise<PostType[]> => {
  try {
    const response = await api.get(`/users/u/${username}/liked-posts`);
    return response.data as PostType[];
  } catch (error) {
    return [];
  }
};

export const logout = () => {
  localStorage.removeItem("access-token");
  localStorage.removeItem("refresh-token");
  const { clearUser } = useUserStore.getState();
  clearUser();
};

export const updateUserDescription = async (
  userID: string,
  description: string,
): Promise<User | null> => {
  try {
    const response = await api.patch(`/users/${userID}`, { description });
    const updatedUser = response.data as User;
    if (updatedUser) {
      const { setUser } = useUserStore.getState();
      setUser(updatedUser);
      return updatedUser;
    }
    return null;
  } catch (error) {
    return null;
  }
};

export const deleteUserAccount = async (userID: string): Promise<boolean> => {
  try {
    const response = await api.delete(`/users/${userID}`);
    if (response.status === 204) {
      return true;
    }
    return false;
  } catch (error) {
    return false;
  }
};

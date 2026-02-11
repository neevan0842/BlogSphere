import toast from "react-hot-toast";
import { api } from "./api";
import { getUserIDFromToken } from "../utils/auth.utils";
import useUserStore from "../store/userStore";
import type { LikedPost, User, UserPost } from "../types/types";

const getGoogleAuthURL = async (): Promise<string> => {
  try {
    const response = await api.get("/auth/google");
    return response.data.url || "";
  } catch (error) {
    toast.error("Failed to get Google authentication URL. Please try again.");
    return "";
  }
};

const exchangeGoogleCodeForToken = async (
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
      toast.error(
        "Failed to decode user information from token. Please try again.",
      );
      return false;
    }

    //getch user details
    const user = await getUserDetailsFromUserID(userID);
    if (!user) {
      toast.error("Failed to fetch user details. Please try again.");
      return false;
    }

    // Update user store with authenticated user ID
    const { setUser } = useUserStore.getState();
    setUser(user);
    toast.success("Successfully signed in with Google!");
    return true;
  } catch (error) {
    toast.error("Google authentication failed. Please try again.");
    return false;
  }
};

const refreshAccessToken = async (): Promise<boolean> => {
  try {
    const response = await api.post("/auth/refresh", {
      refresh_token: localStorage.getItem("refresh-token"),
    });

    // Decode user ID from new access token
    const userID = getUserIDFromToken(response.data.access_token);
    if (!userID) {
      toast.error(
        "Failed to decode user information from refreshed token. Please sign in again.",
      );
      return false;
    }

    //fetch user details
    const user = await getUserDetailsFromUserID(userID);
    if (!user) {
      toast.error("Failed to fetch user details. Please sign in again.");
      return false;
    }

    // Update user store with authenticated user ID and new access token
    const { setUser } = useUserStore.getState();
    setUser(user);
    localStorage.setItem("access-token", response.data.access_token);
    return true;
  } catch (error) {
    toast.error("Failed to refresh access token. Please sign in again.");
    return false;
  }
};

const getUserDetailsFromUserID = async (
  userID: string,
): Promise<User | null> => {
  try {
    const response = await api.get(`/users/${userID}`);
    return response.data as User;
  } catch (error) {
    toast.error("Failed to fetch user details. Please try again.");
    return null;
  }
};

const getUserDetailsFromUsername = async (
  username: string,
): Promise<User | null> => {
  try {
    const response = await api.get(`/users/u/${username}`);
    return response.data as User;
  } catch (error) {
    toast.error("Failed to fetch user details. Please try again.");
    return null;
  }
};

const getUserPosts = async (username: string): Promise<UserPost[]> => {
  try {
    const response = await api.get(`/users/u/${username}/posts`);
    return response.data as UserPost[];
  } catch (error) {
    toast.error("Failed to fetch user's posts. Please try again.");
    return [];
  }
};

const getUserLikedPosts = async (username: string): Promise<LikedPost[]> => {
  try {
    const response = await api.get(`/users/u/${username}/liked-posts`);
    return response.data as LikedPost[];
  } catch (error) {
    toast.error("Failed to fetch user's liked posts. Please try again.");
    return [];
  }
};

const logout = () => {
  localStorage.removeItem("access-token");
  localStorage.removeItem("refresh-token");
  const { clearUser } = useUserStore.getState();
  clearUser();
  toast.success("Successfully logged out.");
};

const updateUserDescription = async (
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
    toast.error("Failed to update user description. Please try again.");
    return null;
  }
};

export {
  getGoogleAuthURL,
  exchangeGoogleCodeForToken,
  refreshAccessToken,
  getUserDetailsFromUserID,
  getUserDetailsFromUsername,
  getUserPosts,
  getUserLikedPosts,
  logout,
  updateUserDescription,
};

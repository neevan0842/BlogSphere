import toast from "react-hot-toast";
import { apiUnauthenticated } from "./api";
import { getUserIDFromToken } from "../utils/auth.utils";
import useUserStore from "../store/userStore";

const getGoogleAuthURL = async (): Promise<string> => {
  try {
    const response = await apiUnauthenticated.get("/auth/google");
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
    const response = await apiUnauthenticated.get("/auth/google/callback", {
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

    // Update user store with authenticated user ID
    const { setUserID } = useUserStore();
    setUserID(userID);
    toast.success("Successfully signed in with Google!");
    return true;
  } catch (error) {
    toast.error("Google authentication failed. Please try again.");
    return false;
  }
};

const refreshAccessToken = async (): Promise<boolean> => {
  try {
    const response = await apiUnauthenticated.post("/auth/refresh", {
      refresh_token: localStorage.getItem("refresh-token"),
    });
    if (!response.data.access_token) {
      toast.error("Failed to refresh access token. Please sign in again.");
      return false;
    }

    // Decode user ID from new access token
    const userID = getUserIDFromToken(response.data.access_token);
    if (!userID) {
      toast.error(
        "Failed to decode user information from refreshed token. Please sign in again.",
      );
      return false;
    }

    // Update user store with authenticated user ID and new access token
    const { setUserID } = useUserStore();
    setUserID(userID);
    localStorage.setItem("access-token", response.data.access_token);
    return true;
  } catch (error) {
    toast.error("Failed to refresh access token. Please sign in again.");
    return false;
  }
};

export { getGoogleAuthURL, exchangeGoogleCodeForToken, refreshAccessToken };

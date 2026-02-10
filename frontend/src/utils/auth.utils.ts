import { jwtDecode, type JwtPayload } from "jwt-decode";

interface CustomJwtPayload extends JwtPayload {
  userID?: string;
}
const getUserIDFromToken = (token: string): string | null => {
  try {
    const decoded = jwtDecode<CustomJwtPayload>(token);
    if (decoded.userID) {
      return decoded.userID || null;
    }
    return null;
  } catch (error) {
    return null;
  }
};

const isTokenExpired = (token: string): boolean => {
  try {
    const decoded = jwtDecode<CustomJwtPayload>(token);
    const tokenExpiration = decoded.exp;
    const now = Date.now() / 1000;
    return tokenExpiration! < now;
  } catch (error) {
    console.error("Token decoding error:", error);
    return true; // Treat decoding errors as if the token is near expiration
  }
};

export { getUserIDFromToken, isTokenExpired };
export type { CustomJwtPayload };

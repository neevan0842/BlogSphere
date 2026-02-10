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

export { getUserIDFromToken };

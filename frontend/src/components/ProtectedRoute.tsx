import { Navigate, Outlet } from "react-router";
import { useEffect, useState } from "react";
import { isTokenExpired } from "../utils/auth.utils";
import useUserStore from "../store/userStore";
import { refreshAccessToken } from "../api/userApi";

function ProtectedRoute() {
  const { clearUser } = useUserStore();
  const [isAuthorized, setIsAuthorized] = useState<boolean | null>(null);

  useEffect(() => {
    const auth = async () => {
      const accessToken = localStorage.getItem("access-token");
      const refreshToken = localStorage.getItem("refresh-token");

      // Check if tokens exist and refresh token is valid
      if (!accessToken || !refreshToken || isTokenExpired(refreshToken)) {
        setIsAuthorized(false);
        clearUser();
        localStorage.removeItem("access-token");
        localStorage.removeItem("refresh-token");
        return;
      }

      const isExpired = isTokenExpired(accessToken);
      if (isExpired) {
        // Clear user state and tokens if the access token is expired
        clearUser();
        localStorage.removeItem("access-token");

        // Attempt to refresh the access token using the refresh token
        const success = await refreshAccessToken();
        if (!success) {
          setIsAuthorized(false);
          clearUser();
          return;
        }
        setIsAuthorized(true);
      } else {
        setIsAuthorized(true);
      }
    };

    auth();
  }, [clearUser]);

  if (isAuthorized === null) {
    return (
      <div>
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
            <p className="mt-4 text-muted-foreground">
              Checking authentication...
            </p>
          </div>
        </div>
      </div>
    );
  }

  return isAuthorized ? <Outlet /> : <Navigate to="/signin" />;
}

export default ProtectedRoute;

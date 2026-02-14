import { Route, Routes } from "react-router";
import { useEffect } from "react";
import Home from "./pages/Home";
import Profile from "./pages/Profile";
import Signin from "./pages/Signin";
import CreateOrUpdate from "./pages/CreateOrUpdate";
import GoogleCallback from "./components/GoogleCallback";
import NotFound from "./pages/NotFound";
import Categories from "./pages/Categories";
import Post from "./pages/Post";
import { Toaster } from "react-hot-toast";
import useThemeStore from "./store/themeStore";
import useUserStore from "./store/userStore";
import ProtectedRoute from "./components/ProtectedRoute";
import { isTokenExpired } from "./utils/auth.utils";

const App = () => {
  const { theme } = useThemeStore();
  const { isAuthenticated, clearUser } = useUserStore();

  // Validate auth state on mount
  useEffect(() => {
    if (isAuthenticated) {
      const refreshToken = localStorage.getItem("refresh-token");
      if (!refreshToken || isTokenExpired(refreshToken)) {
        clearUser();
        localStorage.removeItem("access-token");
        localStorage.removeItem("refresh-token");
      }
    }
  }, [isAuthenticated, clearUser]);

  useEffect(() => {
    document.documentElement.classList.remove("light", "dark");
    document.documentElement.classList.add(theme);
  }, [theme]);

  return (
    <div>
      <Toaster position="top-center" reverseOrder={false} />
      <Routes>
        <Route index element={<Home />} />
        <Route path="/signin" element={<Signin />} />
        <Route path="/categories" element={<Categories />} />
        <Route path="/categories/:slug" element={<Categories />} />
        <Route path="/post/:slug" element={<Post />} />
        <Route path="/u/:username" element={<Profile />} />

        {/* Protected routes */}
        <Route element={<ProtectedRoute />}>
          <Route
            path="/create"
            element={<CreateOrUpdate isUpdatePage={false} />}
          />
          <Route
            path="/post/:postId/edit"
            element={<CreateOrUpdate isUpdatePage={true} />}
          />
        </Route>

        {/* Google OAuth callback route */}
        <Route path="/auth/google/callback" element={<GoogleCallback />} />

        {/* Fallback route */}
        <Route path="/not-found" element={<NotFound />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </div>
  );
};

export default App;

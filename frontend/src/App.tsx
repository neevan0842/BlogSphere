import { Route, Routes } from "react-router";
import { useEffect } from "react";
import Home from "./pages/Home";
import Profile from "./pages/Profile";
import Signin from "./pages/Signin";
import Create from "./pages/Create";
import GoogleCallback from "./components/GoogleCallback";
import NotFound from "./pages/NotFound";
import Categories from "./pages/Categories";
import Post from "./pages/Post";
import { Toaster } from "react-hot-toast";
import useThemeStore from "./store/themeStore";
import ProtectedRoute from "./components/ProtectedRoute";

const App = () => {
  const { theme } = useThemeStore();

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
          <Route path="/create" element={<Create />} />
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

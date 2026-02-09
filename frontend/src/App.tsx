import { Route, Routes } from "react-router";
import Home from "./pages/Home";
import Profile from "./pages/Profile";
import Signin from "./pages/Signin";
import Create from "./pages/Create";
import GoogleCallback from "./components/GoogleCallback";
import NotFound from "./pages/NotFound";
import Categories from "./pages/Categories";
import Post from "./pages/Post";

const App = () => {
  return (
    <div>
      <Routes>
        <Route index element={<Home />} />
        <Route path="/signin" element={<Signin />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/create" element={<Create />} />
        <Route path="/categories" element={<Categories />} />
        <Route path="/post/:slug" element={<Post />} />

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

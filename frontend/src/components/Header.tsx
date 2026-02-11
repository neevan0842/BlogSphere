import { useState } from "react";
import { Link } from "react-router";
import useUserStore from "../store/userStore";
import useThemeStore from "../store/themeStore";
import { Menu, Moon, Sun, X } from "lucide-react";

const Header = () => {
  const { isAuthenticated, user } = useUserStore();
  const { theme, toggleTheme } = useThemeStore();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="container mx-auto flex items-center justify-between px-4 py-4">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-white font-bold text-lg">
            B
          </div>
          <span className="text-xl font-bold text-primary hidden sm:inline">
            Blogsphere
          </span>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex items-center gap-8">
          <Link
            to="/"
            className="text-foreground hover:text-primary transition-colors font-medium"
          >
            Home
          </Link>
          <Link
            to="/categories"
            className="text-foreground hover:text-primary transition-colors font-medium"
          >
            Categories
          </Link>
          {isAuthenticated && (
            <Link
              to={`/u/${user?.username || ""}`}
              className="text-foreground hover:text-primary transition-colors font-medium"
            >
              Profile
            </Link>
          )}
        </nav>

        {/* Right Actions */}
        <div className="flex items-center gap-4">
          {/* Theme Toggle */}
          <button
            onClick={toggleTheme}
            className="inline-flex items-center justify-center rounded-lg border border-border hover:bg-secondary transition-colors p-2"
            aria-label="Toggle theme"
          >
            {theme === "dark" ? (
              <Sun className="h-5 w-5 text-primary" />
            ) : (
              <Moon className="h-5 w-5 text-primary" />
            )}
          </button>

          {/* Login Button / Profile */}
          {!isAuthenticated ? (
            <Link
              to="/signin"
              className="hidden sm:inline-flex items-center justify-center px-4 py-2 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium"
            >
              Sign In
            </Link>
          ) : (
            <Link
              to={`/u/${user?.username || ""}`}
              className="sm:block w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-sm font-bold"
            >
              <img
                src={user?.avatar_url || "/placeholder.svg"}
                alt={user?.username || "User"}
                className="w-8 h-8 rounded-full object-cover"
                referrerPolicy="no-referrer"
              />
            </Link>
          )}

          {/* Mobile Menu Toggle */}
          <button
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            className="md:hidden p-2"
            aria-label="Toggle mobile menu"
          >
            {isMobileMenuOpen ? (
              <X className="h-5 w-5 text-foreground" />
            ) : (
              <Menu className="h-5 w-5 text-foreground" />
            )}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      {isMobileMenuOpen && (
        <div className="border-t border-border bg-card md:hidden">
          <nav className="container mx-auto flex flex-col gap-2 px-4 py-4">
            <Link
              to="/"
              className="px-4 py-2 rounded-lg hover:bg-secondary transition-colors text-foreground font-medium"
            >
              Home
            </Link>
            <Link
              to="/categories"
              className="px-4 py-2 rounded-lg hover:bg-secondary transition-colors text-foreground font-medium"
            >
              Categories
            </Link>
            {isAuthenticated && (
              <Link
                to={`/u/${user?.username || ""}`}
                className="px-4 py-2 rounded-lg hover:bg-secondary transition-colors text-foreground font-medium"
              >
                Profile
              </Link>
            )}
            {!isAuthenticated && (
              <Link
                to="/signin"
                className="px-4 py-2 rounded-lg bg-primary text-primary-foreground text-center font-medium"
              >
                Sign In
              </Link>
            )}
          </nav>
        </div>
      )}
    </header>
  );
};

export default Header;

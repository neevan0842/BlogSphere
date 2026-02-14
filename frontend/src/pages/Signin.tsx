import { Link, useNavigate } from "react-router";
import { getGoogleAuthURL } from "../api/userApi";
import toast from "react-hot-toast";
import { useEffect } from "react";
import useUserStore from "../store/userStore";

const Login = () => {
  const { isAuthenticated } = useUserStore();
  const navigate = useNavigate();

  const handleGoogleSignIn = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    const url = await getGoogleAuthURL();
    if (url !== "") {
      window.location.href = url;
    } else {
      toast.error("Failed to get Google authentication URL.");
    }
  };

  useEffect(() => {
    if (isAuthenticated) {
      navigate("/");
    }
  }, [navigate, isAuthenticated]);

  return (
    <div>
      <div className="min-h-screen flex overflow-hidden">
        {/* Left Panel - Sign In Form */}
        <div className="w-full lg:w-1/2 flex items-center justify-center px-4 py-8 bg-background">
          <div className="w-full max-w-sm">
            {/* Logo */}
            <Link to="/" className="inline-flex items-center gap-2 mb-8">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary text-white font-bold text-lg">
                B
              </div>
              <span className="text-2xl font-bold text-primary">
                Blogsphere
              </span>
            </Link>

            <div className="space-y-8">
              <div>
                <h1 className="text-3xl font-bold text-foreground mb-2">
                  Welcome Back
                </h1>
                <p className="text-muted-foreground">
                  Sign in to your account to continue sharing your stories
                </p>
              </div>

              {/* Google Sign In Button */}
              <button
                onClick={handleGoogleSignIn}
                className="w-full flex items-center justify-center gap-3 px-6 py-3 rounded-lg border border-border hover:bg-secondary transition-colors font-medium text-foreground disabled:opacity-50"
              >
                <svg
                  className="h-5 w-5"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                    fill="#4285F4"
                  />
                  <path
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                    fill="#34A853"
                  />
                  <path
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                    fill="#FBBC05"
                  />
                  <path
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                    fill="#EA4335"
                  />
                </svg>
                Sign in with Google
              </button>

              {/* Divider */}
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-border"></div>
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-2 bg-background text-muted-foreground">
                    Or continue as guest
                  </span>
                </div>
              </div>

              {/* Guest Sign In */}
              <Link
                to="/"
                className="w-full flex items-center justify-center px-6 py-3 rounded-lg bg-secondary hover:bg-secondary/80 transition-colors font-medium text-foreground"
              >
                Browse as Guest
              </Link>
            </div>

            {/* Footer */}
            <p className="text-center text-sm text-muted-foreground mt-8">
              By signing in, you agree to our
              <Link to="#" className="text-primary hover:underline">
                {` Terms of Service `}
              </Link>
              and
              <Link to="#" className="text-primary hover:underline">
                {` Privacy Policy `}
              </Link>
            </p>
          </div>
        </div>

        {/* Right Panel - Welcome Message */}
        <div className="hidden lg:flex w-1/2 bg-linear-to-br from-primary/20 via-primary/10 to-accent/10 items-center justify-center p-8">
          <div className="max-w-md text-center space-y-6">
            <div className="h-16 w-16 rounded-full bg-primary/20 flex items-center justify-center mx-auto">
              <svg
                className="h-8 w-8 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 6.253v13m0-13C6.5 6.253 2 10.88 2 17s4.5 10.747 10 10.747c5.5 0 10-4.88 10-10.747 0-6.12-4.5-10.747-10-10.747z"
                />
              </svg>
            </div>

            <h2 className="text-3xl font-bold text-foreground">
              Share Your Unique Voice
            </h2>

            <p className="text-lg text-muted-foreground leading-relaxed">
              Blogsphere is a welcoming community where writers and readers
              connect. Share your thoughts, stories, and expertise with an
              audience that cares.
            </p>

            <div className="space-y-3 pt-4">
              <div className="flex gap-3 items-start text-left">
                <div className="h-6 w-6 rounded-full bg-primary/30 flex items-center justify-center shrink-0 text-primary text-sm font-bold">
                  ✓
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    Easy to Start
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Sign in with Google and start writing in minutes
                  </p>
                </div>
              </div>

              <div className="flex gap-3 items-start text-left">
                <div className="h-6 w-6 rounded-full bg-primary/30 flex items-center justify-center shrink-0 text-primary text-sm font-bold">
                  ✓
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    Build Community
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Engage with readers through comments and likes
                  </p>
                </div>
              </div>

              <div className="flex gap-3 items-start text-left">
                <div className="h-6 w-6 rounded-full bg-primary/30 flex items-center justify-center shrink-0 text-primary text-sm font-bold">
                  ✓
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    Express Yourself
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Use our rich editor to tell your stories beautifully
                  </p>
                </div>
              </div>
            </div>

            <blockquote className="border-l-4 border-primary pl-4 py-4 text-left mt-8">
              <p className="text-sm italic text-muted-foreground mb-2">
                "Every story matters. Every voice deserves to be heard."
              </p>
              <footer className="text-xs text-muted-foreground">
                — Blogsphere Community
              </footer>
            </blockquote>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;

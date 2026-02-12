import { useEffect, useRef } from "react";
import toast from "react-hot-toast";
import { useNavigate, useSearchParams } from "react-router";
import { exchangeGoogleCodeForToken } from "../api/userAuth";

const GoogleCallback = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const hasProcessed = useRef(false);

  useEffect(() => {
    // Prevent duplicate processing (React StrictMode runs effects twice)
    if (hasProcessed.current) return;

    const processGoogleCallback = async () => {
      hasProcessed.current = true;

      const code = searchParams.get("code");
      const state = searchParams.get("state");

      if (!code || !state) {
        toast.error("Google authentication failed. Please try again.");
        navigate("/signin");
        return;
      }

      const success = await exchangeGoogleCodeForToken(code, state);
      if (success) {
        toast.success("Signed in with Google successfully!");
        navigate("/");
      } else {
        toast.error("Failed to sign in with Google. Please try again.");
        navigate("/signin");
      }
    };

    processGoogleCallback();
  }, [searchParams, navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
        <p className="mt-4 text-muted-foreground">Completing sign in...</p>
      </div>
    </div>
  );
};

export default GoogleCallback;

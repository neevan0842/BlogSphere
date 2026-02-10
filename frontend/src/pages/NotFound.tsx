import { ArrowLeft, Home } from "lucide-react";
import { Link } from "react-router";
import PageLayout from "../components/PageLayout";
import useUserStore from "../store/userStore";

const NotFound = () => {
  const { isAuthenticated } = useUserStore();

  return (
    <PageLayout>
      <section className="container mx-auto px-4 py-16 md:py-24">
        <div className="max-w-2xl mx-auto text-center space-y-8">
          {/* 404 Icon/Number */}
          <div className="space-y-4">
            <h1 className="text-8xl md:text-9xl font-bold text-primary opacity-20">
              404
            </h1>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground">
              Page Not Found
            </h2>
            <p className="text-lg text-muted-foreground max-w-md mx-auto">
              Sorry, we couldn't find the page you're looking for. It might have
              been moved or deleted.
            </p>
          </div>

          {/* Action Buttons */}
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4 pt-4">
            <Link
              to="/"
              className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium"
            >
              <Home className="h-4 w-4" />
              Go to Home
            </Link>
            <button
              onClick={() => window.history.back()}
              className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-lg border border-border text-foreground hover:bg-secondary transition-colors font-medium"
            >
              <ArrowLeft className="h-4 w-4" />
              Go Back
            </button>
          </div>

          {/* Helpful Links */}
          <div className="pt-8 border-t border-border">
            <p className="text-sm text-muted-foreground mb-4">
              You might be interested in:
            </p>
            <div className="flex flex-wrap items-center justify-center gap-4 text-sm">
              <Link to="/" className="text-primary hover:underline font-medium">
                Home
              </Link>
              <span className="text-muted-foreground">•</span>
              <Link
                to="/categories"
                className="text-primary hover:underline font-medium"
              >
                Categories
              </Link>
              {isAuthenticated && (
                <>
                  <span className="text-muted-foreground">•</span>
                  <Link
                    to="/create"
                    className="text-primary hover:underline font-medium"
                  >
                    Write a Post
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </section>
    </PageLayout>
  );
};

export default NotFound;

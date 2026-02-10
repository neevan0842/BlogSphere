import { Link } from "react-router";

const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="border-t border-border bg-secondary/50 mt-12">
      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="text-center md:text-left">
            <h3 className="text-lg font-bold text-primary mb-1">Blogsphere</h3>
            <p className="text-muted-foreground text-sm">
              Share your stories, connect with readers, inspire the world.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-6 text-sm text-muted-foreground">
            <Link
              to="https://github.com/neevan0842"
              target="_blank"
              className="hover:text-primary transition-colors"
            >
              About
            </Link>
            <Link to="#" className="hover:text-primary transition-colors">
              Privacy
            </Link>
            <Link to="#" className="hover:text-primary transition-colors">
              Terms
            </Link>
            <Link
              to="https://github.com/neevan0842"
              target="_blank"
              className="hover:text-primary transition-colors"
            >
              Contact
            </Link>
          </div>
        </div>

        <div className="border-t border-border mt-6 pt-6 text-center text-sm text-muted-foreground">
          <p>&copy; {currentYear} Blogsphere. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;

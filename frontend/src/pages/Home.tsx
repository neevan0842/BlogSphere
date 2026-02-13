import PageLayout from "../components/PageLayout";
import SearchDisplayPost from "../components/SearchDisplayPost";
import useUserStore from "../store/userStore";

const Home = () => {
  const { isAuthenticated } = useUserStore();

  return (
    <PageLayout>
      {/* Hero Section */}
      <section className="border-b border-border bg-secondary/30">
        <div className="container mx-auto px-4 py-12 md:py-16">
          <div className="max-w-2xl mx-auto text-center space-y-6">
            <h1 className="text-4xl md:text-5xl font-bold text-balance text-foreground">
              Discover Stories That Matter
            </h1>
            <p className="text-lg text-muted-foreground text-balance">
              Read insightful articles, share your thoughts, and connect with a
              community of writers and readers passionate about ideas.
            </p>

            {isAuthenticated && (
              <a
                href="/create"
                className="inline-flex items-center justify-center px-6 py-3 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium"
              >
                Write Your Story
              </a>
            )}
          </div>
        </div>
      </section>

      {/* Search and Posts Section */}
      <SearchDisplayPost />
    </PageLayout>
  );
};

export default Home;

import debounce from "lodash.debounce";
import { useEffect, useRef, useState } from "react";
import type { PostType } from "../types/types";
import { getPostsPaginated } from "../api/postApi";
import { ArrowLeft, Search } from "lucide-react";
import BlogPostCard from "./BlogPostCard";
import InfiniteScroll from "react-infinite-scroll-component";
import { Link } from "react-router";

interface SearchDisplayPostProps {
  categorySlug?: string;
  showBackButton?: boolean;
}

const SearchDisplayPost = ({
  categorySlug = "",
  showBackButton = false,
}: SearchDisplayPostProps) => {
  const LIMIT = 20;
  const [posts, setPosts] = useState<PostType[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [inputValue, setInputValue] = useState("");
  const [query, setQuery] = useState("");

  const fetchPosts = async (reset = false, searchQuery?: string) => {
    const currentQuery = searchQuery ?? query;
    const pageToFetch = reset ? 1 : page;
    const result = await getPostsPaginated(
      currentQuery,
      categorySlug,
      pageToFetch,
      LIMIT,
    );
    const { posts: data, hasMore: more } = result;
    setPosts((prev) => (reset ? data : [...prev, ...data]));
    setHasMore(more);
    setPage(pageToFetch + 1);
  };

  //debounce search input
  const debouncedSearch = useRef(
    debounce((value: string) => {
      setPage(1);
      setPosts([]);
      setQuery(value);
      setHasMore(true);
      fetchPosts(true, value);
    }, 500),
  ).current;

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
    debouncedSearch(e.target.value.trim());
  };

  // Fetch initial posts on mount
  useEffect(() => {
    fetchPosts(true);
  }, [categorySlug]);

  useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  return (
    <div>
      {/* Search Bar */}
      <section className="sticky top-16 z-30 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/80 border-b border-border py-4">
        <div className="container mx-auto px-4">
          <div className="max-w-2xl mx-auto flex flex-col md:flex-row md:items-center gap-4">
            {showBackButton && (
              <Link
                to="/categories"
                className="inline-flex items-center gap-2 text-muted-foreground hover:text-primary transition-colors text-sm whitespace-nowrap"
              >
                <ArrowLeft className="h-4 w-4" />
                Back to categories
              </Link>
            )}
            <div className="relative flex-1">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
              <input
                type="text"
                placeholder="Search articles..."
                value={inputValue}
                onChange={handleSearchChange}
                className="w-full pl-12 pr-4 py-2.5 rounded-lg border border-border bg-card text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
              />
            </div>
          </div>
        </div>
      </section>

      {/* Blog Posts Grid */}
      <section className="container mx-auto px-4 py-12">
        {posts.length > 0 ? (
          <InfiniteScroll
            dataLength={posts.length}
            next={fetchPosts}
            hasMore={hasMore}
            loader={
              <p className="text-center text-muted-foreground py-4">
                Loading more articles...
              </p>
            }
            endMessage={
              <p className="text-center text-muted-foreground py-4">
                No more articles to load.
              </p>
            }
          >
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {posts.map((post) => (
                <BlogPostCard key={post.id} post={post} />
              ))}
            </div>
          </InfiniteScroll>
        ) : (
          <div className="text-center py-12">
            <div className="inline-flex h-16 w-16 items-center justify-center rounded-full bg-muted mb-4">
              <Search className="h-8 w-8 text-muted-foreground" />
            </div>
            <h2 className="text-xl font-bold text-foreground mb-2">
              No articles found
            </h2>
            <p className="text-muted-foreground">
              Try adjusting your search query or browse all articles.
            </p>
          </div>
        )}
      </section>
    </div>
  );
};

export default SearchDisplayPost;

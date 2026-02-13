import PageLayout from "../components/PageLayout";
import { useEffect, useState } from "react";
import type { CategoryDisplay } from "../types/types";
import { getCategories } from "../api/categoryApi";
import toast from "react-hot-toast";
import SearchDisplayPost from "../components/SearchDisplayPost";
import { capitalize } from "../utils/utils";
import { Link, useNavigate, useParams } from "react-router";

const Categories = () => {
  const { slug } = useParams();
  const navigate = useNavigate();
  const [categories, setCategories] = useState<CategoryDisplay[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchCategories = async () => {
      setLoading(true);
      const data = await getCategories();
      if (!data) {
        toast.error("Failed to load categories. Please try again later.");
        setLoading(false);
        return;
      }
      setCategories(data);
      setLoading(false);
    };
    fetchCategories();
  }, []);

  useEffect(() => {
    if (!loading && categories.length > 0 && slug) {
      const categoryExists = categories.some((cat) => cat.slug === slug);
      if (!categoryExists) {
        toast.error("Category not found.");
        navigate("/categories");
      }
    }
  }, [slug, categories, navigate, loading]);

  return (
    <PageLayout>
      {/* Hero Section */}
      <section className="border-b border-border bg-secondary/30 py-12">
        <div className="container mx-auto px-4">
          <div className="max-w-2xl mx-auto text-center">
            <h1 className="text-4xl md:text-5xl font-bold text-balance text-foreground mb-4">
              {!slug ? "Explore by Category" : `${capitalize(slug)} Articles`}
            </h1>
            <p className="text-lg text-muted-foreground text-balance">
              Browse articles by topic and discover content that interests you.
            </p>
          </div>
        </div>
      </section>

      {/* Categories Grid */}
      <section className="container mx-auto px-4 py-12">
        {loading ? (
          <div className="text-center py-12">
            <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-primary border-r-transparent"></div>
            <p className="mt-4 text-muted-foreground">Loading categories...</p>
          </div>
        ) : !!slug ? (
          <SearchDisplayPost categorySlug={slug} showBackButton={true} />
        ) : categories.length > 0 ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {categories.map((category) => (
              <Link
                key={category.id}
                to={`/categories/${category.slug}`}
                className="p-6 rounded-lg border border-border hover:border-primary/50 hover:shadow-md transition-all duration-200 bg-card text-left group"
              >
                <div className="text-4xl mb-4 opacity-80 group-hover:opacity-100 transition-opacity">
                  {category.icon}
                </div>
                <h2 className="text-2xl font-bold text-foreground group-hover:text-primary transition-colors mb-2">
                  {category.name}
                </h2>
                <p className="text-muted-foreground mb-4">
                  {category.description}
                </p>
                <div className="inline-flex items-center text-primary font-medium group-hover:gap-2 transition-all">
                  View Articles
                  <span>→</span>
                </div>
              </Link>
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-lg text-muted-foreground">
              No categories available yet.
            </p>
          </div>
        )}
      </section>
    </PageLayout>
  );
};

export default Categories;

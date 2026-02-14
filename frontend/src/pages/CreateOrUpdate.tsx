import { ArrowLeft, Bold, Heading2, Italic, List, Save } from "lucide-react";
import PageLayout from "../components/PageLayout";
import { Link, useNavigate, useParams } from "react-router";
import useUserStore from "../store/userStore";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { getCategories } from "../api/categoryApi";
import type { CategoryDisplay } from "../types/types";
import { z } from "zod";
import { createPost, getPostByID, updatePostByID } from "../api/postApi";
import MarkdownRenderer from "../components/MarkdownRenderer";

interface CreateOrUpdateProps {
  isUpdatePage: boolean;
}

const CreateOrUpdate = ({ isUpdatePage = false }: CreateOrUpdateProps) => {
  const { postId } = useParams(); // postId is used only for update page
  const navigate = useNavigate();
  const { isAuthenticated } = useUserStore();
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [categories, setCategories] = useState<CategoryDisplay[]>([]);
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const [isPreview, setIsPreview] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const postSchema = z.object({
    title: z
      .string()
      .min(3, "Title must be at least 3 characters")
      .max(200, "Title is too long"),
    body: z.string().min(10, "content must be at least 10 characters"),
    selectedCategories: z
      .array(z.string())
      .min(1, "At least one category must be selected")
      .max(3, "You can select a maximum of 3 categories"),
  });

  const toggleCategory = (categoryId: string) => {
    setSelectedCategories((prev) => {
      if (prev.includes(categoryId)) {
        return prev.filter((c) => c !== categoryId);
      } else {
        if (prev.length >= 3) {
          toast.error("You can select a maximum of 3 categories.");
          return prev;
        }
        return [...prev, categoryId];
      }
    });
  };

  const handlePublishOrUpdate = async () => {
    const validation = postSchema.safeParse({
      title: title.trim(),
      body: body.trim(),
      selectedCategories,
    });
    if (!validation.success) {
      toast.error(validation.error.issues[0].message);
      return;
    }

    setIsSaving(true);
    let response = null;
    if (isUpdatePage) {
      // Update existing post
      response = await updatePostByID({
        postID: postId || "",
        title: title.trim(),
        body: body.trim(),
        categoryIDs: selectedCategories,
      });
    } else {
      // Create new post
      response = await createPost({
        title: title.trim(),
        body: body.trim(),
        categoryIDs: selectedCategories,
      });
    }
    if (!response) {
      toast.error(`Failed to ${isUpdatePage ? "update" : "create"} post.`);
      setIsSaving(false);
      return;
    }
    toast.success(`Post ${isUpdatePage ? "updated" : "created"} successfully!`);
    navigate(`/post/${response.slug}`);
    setIsSaving(false);
  };

  const handleBodyChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setBody(e.target.value);
    e.target.style.height = "auto";
    e.target.style.height = `${Math.max(360, e.target.scrollHeight)}px`;
  };

  const insertMarkdown = (before: string, after: string = "") => {
    const textarea = document.getElementById("content") as HTMLTextAreaElement;
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const selected = body.substring(start, end);
    const newContent =
      body.substring(0, start) +
      before +
      selected +
      after +
      body.substring(end);
    setBody(newContent);
  };

  useEffect(() => {
    if (!isAuthenticated) {
      toast.error("You must be signed in to create a post.");
      navigate("/signin");
    }
  }, [isAuthenticated, navigate]);

  useEffect(() => {
    const fetchCategories = async () => {
      const data = await getCategories();
      if (!data) {
        toast.error("Failed to load categories.");
        return;
      }
      setCategories(data);
    };
    fetchCategories();
  }, []);

  useEffect(() => {
    const fetchPostDetails = async () => {
      if (!isUpdatePage) {
        return;
      }

      // For update page, we will fetch the post details and populate the form
      const post = await getPostByID(postId || "");
      if (!post) {
        toast.error("Failed to load post details");
        navigate("/");
        return;
      }
      setTitle(post.title);
      setBody(post.body);
      setSelectedCategories(post.categories.map((c) => c.id));
    };
    fetchPostDetails();
  }, [isUpdatePage, postId, navigate]);

  return (
    <PageLayout>
      <div className="container mx-auto px-4 py-8 md:py-12 max-w-4xl">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <Link
            to="/"
            className="inline-flex items-center gap-2 text-primary hover:text-primary/80 transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back
          </Link>
          <h1 className="text-2xl md:text-3xl font-bold">
            {isUpdatePage ? "Update Post" : "Create New Post"}
          </h1>
          <div className="w-20" />
        </div>

        <form
          className="space-y-8"
          onSubmit={(e) => {
            e.preventDefault();
            handlePublishOrUpdate();
          }}
        >
          {/* Title */}
          <div className="space-y-2">
            <label className="block text-sm font-semibold text-foreground">
              Post Title *
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Enter your post title..."
              className="w-full px-4 py-3 rounded-lg border border-border bg-card text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 text-lg"
            />
          </div>

          {/* Content Editor */}
          <div className="space-y-2">
            <label className="block text-sm font-semibold text-foreground">
              Post Content *
            </label>

            {/* Editor Toolbar */}
            <div className="flex items-center justify-between bg-card border border-border rounded-t-lg px-3 py-2">
              <div className="flex flex-wrap gap-2">
                <button
                  type="button"
                  onClick={() => insertMarkdown("**", "**")}
                  className="p-2 rounded hover:bg-secondary transition-colors text-foreground"
                  title="Bold"
                >
                  <Bold className="h-4 w-4" />
                </button>
                <button
                  type="button"
                  onClick={() => insertMarkdown("*", "*")}
                  className="p-2 rounded hover:bg-secondary transition-colors text-foreground"
                  title="Italic"
                >
                  <Italic className="h-4 w-4" />
                </button>
                <button
                  type="button"
                  onClick={() => insertMarkdown("## ")}
                  className="p-2 rounded hover:bg-secondary transition-colors text-foreground"
                  title="Heading"
                >
                  <Heading2 className="h-4 w-4" />
                </button>
                <button
                  type="button"
                  onClick={() => insertMarkdown("- ")}
                  className="p-2 rounded hover:bg-secondary transition-colors text-foreground"
                  title="List"
                >
                  <List className="h-4 w-4" />
                </button>
              </div>
              <button
                type="button"
                onClick={() => setIsPreview((prev) => !prev)}
                className="text-sm font-medium text-primary hover:text-primary/80 transition-colors"
              >
                {isPreview ? "Edit" : "Preview"}
              </button>
            </div>

            {/* Content Textarea */}
            {isPreview ? (
              <MarkdownRenderer
                content={body}
                className="prose prose-sm sm:prose-base lg:prose-lg dark:prose-invert max-w-none w-full px-4 py-3 rounded-b-lg border border-t-0 border-border bg-card min-h-90"
              />
            ) : (
              <textarea
                id="content"
                value={body}
                onChange={handleBodyChange}
                placeholder="Write your story here... You can use markdown formatting."
                className="w-full px-4 py-3 rounded-b-lg border border-t-0 border-border bg-card text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 resize-none overflow-hidden"
                style={{ minHeight: "360px" }}
              />
            )}

            <p className="text-xs text-muted-foreground">
              Tip: Use ** for bold, * for italic, ## for headings, and - for
              lists
            </p>
          </div>

          {/* Categories */}
          <div className="space-y-3">
            <label className="block text-sm font-semibold text-foreground">
              Categories * (select 1-3 categories)
            </label>
            <p className="text-xs text-muted-foreground">
              {selectedCategories.length}/3 categories selected
            </p>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
              {categories.map((category) => {
                const isSelected = selectedCategories.includes(category.id);
                const isDisabled =
                  !isSelected && selectedCategories.length >= 3;

                return (
                  <button
                    key={category.id}
                    type="button"
                    onClick={() => toggleCategory(category.id)}
                    disabled={isDisabled}
                    className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                      isSelected
                        ? "bg-primary text-primary-foreground"
                        : isDisabled
                          ? "border border-border text-muted-foreground opacity-50 cursor-not-allowed"
                          : "border border-border text-foreground hover:bg-secondary"
                    }`}
                  >
                    {category.name}
                  </button>
                );
              })}
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-4 pt-6 border-t border-border">
            <Link
              to="/"
              className="px-6 py-3 rounded-lg border border-border text-foreground hover:bg-secondary transition-colors font-medium"
            >
              Cancel
            </Link>
            <button
              type="submit"
              disabled={isSaving}
              className="px-6 py-3 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium flex items-center gap-2 disabled:opacity-50"
            >
              <Save className="h-4 w-4" />
              {isUpdatePage ? "Update Post" : "Publish Post"}
            </button>
          </div>
        </form>
      </div>
    </PageLayout>
  );
};

export default CreateOrUpdate;

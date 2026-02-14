import { Link, useNavigate, useParams } from "react-router";
import PageLayout from "../components/PageLayout";
import { useEffect, useState } from "react";
import type { PostType, User } from "../types/types";
import { Edit2, LogOut, Trash2 } from "lucide-react";
import {
  deleteUserAccount,
  getUserDetailsFromUsername,
  getUserLikedPosts,
  getUserPosts,
  logout,
  updateUserDescription,
} from "../api/userApi";
import BlogPostCard from "../components/BlogPostCard";
import useUserStore from "../store/userStore";
import toast from "react-hot-toast";
import { deletePostByID } from "../api/postApi";
import ConfirmModal from "../components/ConfirmModal";

const Profile = () => {
  const navigate = useNavigate();
  const { username } = useParams();
  const { user: authenticatedUser, isAuthenticated } = useUserStore();
  const [user, setUser] = useState<User | null>(null);
  const [isOwner, setIsOwner] = useState<boolean>(false);
  const [userPosts, setUserPosts] = useState<PostType[]>([]);
  const [likedPosts, setLikedPosts] = useState<PostType[]>([]);
  const [activeTab, setActiveTab] = useState<"posts" | "liked">("posts");
  const [isEditing, setIsEditing] = useState(false);
  const [editDescription, setEditDescription] = useState("");
  const [showDeleteAccountModal, setShowDeleteAccountModal] = useState(false);
  const [showDeletePostModal, setShowDeletePostModal] = useState(false);
  const [postToDelete, setPostToDelete] = useState<string | null>(null);

  const handleLogout = () => {
    logout();
    navigate("/");
    toast.success("Logged out successfully.");
  };

  const handleDeleteAccountClick = () => {
    if (!isOwner) {
      toast.error("You do not have permission to delete this account.");
      return;
    }
    setShowDeleteAccountModal(true);
  };

  const handleDeleteAccountConfirm = async () => {
    const result = await deleteUserAccount(user?.id || "");
    if (result) {
      logout();
      navigate("/", { replace: true });
      toast.success("Account deleted successfully.");
    } else {
      toast.error("Failed to delete account.");
    }
  };

  const handlePostDeleteClick = (postID: string) => {
    if (!isOwner) {
      toast.error("You do not have permission to delete this post.");
      return;
    }
    setPostToDelete(postID);
    setShowDeletePostModal(true);
  };

  const handlePostDeleteConfirm = async () => {
    if (!postToDelete) return;

    const result = await deletePostByID(postToDelete);
    if (result) {
      toast.success("Post deleted successfully.");
      setUserPosts((prevPosts) =>
        prevPosts.filter((post) => post.id !== postToDelete),
      );
    } else {
      toast.error("Failed to delete post.");
    }
    setPostToDelete(null);
  };

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    const updatedUser = await updateUserDescription(
      user?.id || "",
      editDescription,
    );
    if (!updatedUser) {
      return;
    }
    setUser(updatedUser);
    setIsEditing(false);
    toast.success("Profile updated successfully.");
  };

  useEffect(() => {
    let isMounted = true;

    const fetchUserProfileDetails = async () => {
      const [userData, userPostsData, likedPostsData] = await Promise.all([
        getUserDetailsFromUsername(username || ""),
        getUserPosts(username || ""),
        getUserLikedPosts(username || ""),
      ]);

      // Only update state and show errors if component is still mounted
      if (isMounted) {
        if (userData) {
          setUser(userData);
          setUserPosts(userPostsData || []);
          setLikedPosts(likedPostsData || []);
          setIsOwner(isAuthenticated && authenticatedUser?.id === userData.id);
        } else {
          toast.error("User not found.");
          navigate("/not-found");
        }
      }
    };

    fetchUserProfileDetails();

    return () => {
      isMounted = false;
    };
  }, [username]);

  // Fetch user details on component mount
  if (!user) {
    return (
      <>
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
            <p className="mt-4 text-muted-foreground">Loading profile...</p>
          </div>
        </div>
      </>
    );
  }

  return (
    <PageLayout>
      {/* Profile Header */}
      <section className="border-b border-border bg-secondary/30 py-12">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto">
            <div className="flex flex-col md:flex-row items-start md:items-center gap-8">
              {/* Avatar */}
              <img
                src={user.avatar_url || "/placeholder.svg"}
                alt={user.username || "User"}
                className="h-24 w-24 md:h-32 md:w-32 rounded-full object-cover border-4 border-primary/20"
              />

              {/* Profile Info */}
              <div className="flex-1">
                {!isEditing ? (
                  <>
                    <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-2">
                      {user.username || "Anonymous"}
                    </h1>
                    <p className="text-muted-foreground mb-4">{user.email}</p>
                    <p className="text-foreground mb-6 max-w-xl">
                      {user.description || "No description available."}
                    </p>

                    <div className="space-y-4">
                      {/* Action Buttons */}
                      {isOwner && (
                        <div>
                          <button
                            onClick={() => setIsEditing(true)}
                            className="inline-flex items-center gap-2 px-4 py-2 mr-2 rounded-lg border border-border text-foreground hover:bg-secondary transition-colors font-medium"
                          >
                            <Edit2 className="h-4 w-4" />
                            Edit Profile
                          </button>
                          <button
                            onClick={handleLogout}
                            className="inline-flex items-center gap-2 px-4 py-2 rounded-lg border border-border text-foreground hover:bg-secondary transition-colors font-medium"
                          >
                            <LogOut className="h-4 w-4" />
                            Logout
                          </button>
                        </div>
                      )}
                    </div>
                  </>
                ) : (
                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-semibold text-foreground mb-1">
                        Description
                      </label>
                      <input
                        type="text"
                        value={editDescription}
                        autoFocus
                        onChange={(e) => setEditDescription(e.target.value)}
                        className="w-full px-3 py-2 rounded-lg border border-border bg-card text-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
                      />
                    </div>

                    <div className="flex gap-2 pt-4">
                      <button
                        onClick={handleSaveProfile}
                        className="px-4 py-2 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium"
                      >
                        Save Changes
                      </button>
                      <button
                        onClick={() => setIsEditing(false)}
                        className="px-4 py-2 rounded-lg border border-border text-foreground hover:bg-secondary transition-colors font-medium"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Tabs */}
      <section className="border-b border-border sticky top-16 z-20 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/80">
        <div className="container mx-auto px-4">
          <div className="flex gap-8 max-w-3xl mx-auto">
            <button
              onClick={() => setActiveTab("posts")}
              className={`px-4 py-4 font-medium border-b-2 transition-colors ${
                activeTab === "posts"
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground"
              }`}
            >
              Your Posts ({userPosts.length})
            </button>
            <button
              onClick={() => setActiveTab("liked")}
              className={`px-4 py-4 font-medium border-b-2 transition-colors ${
                activeTab === "liked"
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground"
              }`}
            >
              Liked Posts ({likedPosts.length})
            </button>
          </div>
        </div>
      </section>

      {/* Content */}
      <section className="container mx-auto px-4 py-12">
        <div className="max-w-3xl mx-auto">
          {activeTab === "posts" && (
            <div>
              {userPosts.length > 0 ? (
                <div className="grid gap-6 md:grid-cols-2">
                  {userPosts.map((post) => (
                    <div key={post.id} className="relative">
                      <BlogPostCard post={post} />
                      {isOwner && (
                        <div className="absolute top-4 right-4 flex gap-2">
                          <Link
                            to={`/post/${post.id}/edit`}
                            className="p-2 rounded-lg bg-background/80 text-foreground hover:bg-background transition-colors"
                          >
                            <Edit2 className="h-4 w-4" />
                          </Link>
                          <button
                            onClick={() => handlePostDeleteClick(post.id)}
                            className="p-2 rounded-lg bg-background/80 text-destructive hover:bg-background transition-colors"
                          >
                            <Trash2 className="h-4 w-4" />
                          </button>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-12">
                  <p className="text-muted-foreground mb-4">
                    You haven't published any posts yet.
                  </p>
                  <Link
                    to="/create"
                    className="inline-flex items-center justify-center px-6 py-3 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors font-medium"
                  >
                    Write Your First Post
                  </Link>
                </div>
              )}
            </div>
          )}

          {activeTab === "liked" && (
            <div>
              {likedPosts.length > 0 ? (
                <div className="grid gap-6 md:grid-cols-2">
                  {likedPosts.map((post) => (
                    <BlogPostCard key={post.id} post={post} />
                  ))}
                </div>
              ) : (
                <div className="text-center py-12">
                  <p className="text-muted-foreground">
                    You haven't liked any posts yet. Explore articles and show
                    some love!
                  </p>
                </div>
              )}
            </div>
          )}
        </div>
      </section>

      {/* Danger Zone */}
      {!isEditing && isOwner && (
        <section className="container mx-auto px-4 py-12">
          <div className="max-w-3xl mx-auto">
            <div className="border-2 border-destructive rounded-lg p-6 bg-destructive/15">
              <h3 className="text-lg font-bold text-destructive mb-2">
                Danger Zone
              </h3>
              <p className="text-destructive/80 mb-4">
                Once you delete your account, there is no going back. Please be
                certain.
              </p>
              <button
                onClick={handleDeleteAccountClick}
                className="px-4 py-2 rounded-lg bg-destructive text-destructive-foreground hover:bg-destructive/90 transition-colors font-medium"
              >
                Delete Account
              </button>
            </div>
          </div>
        </section>
      )}

      {/* Delete Account Modal */}
      <ConfirmModal
        isOpen={showDeleteAccountModal}
        onClose={() => setShowDeleteAccountModal(false)}
        onConfirm={handleDeleteAccountConfirm}
        title="Delete Account"
        message="Are you sure you want to delete your account? This action cannot be undone and all your posts will be permanently deleted."
        confirmText="Delete Account"
        cancelText="Cancel"
        variant="danger"
      />

      {/* Delete Post Modal */}
      <ConfirmModal
        isOpen={showDeletePostModal}
        onClose={() => {
          setShowDeletePostModal(false);
          setPostToDelete(null);
        }}
        onConfirm={handlePostDeleteConfirm}
        title="Delete Post"
        message="Are you sure you want to delete this post? This action cannot be undone."
        confirmText="Delete Post"
        cancelText="Cancel"
        variant="danger"
      />
    </PageLayout>
  );
};

export default Profile;

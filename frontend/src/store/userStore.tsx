import { create } from "zustand";
import { persist } from "zustand/middleware";
import { immer } from "zustand/middleware/immer";

type UserStoreState = {
  userID: string | null;
  isAuthenticated: boolean;
  setUserID: (userID: string) => void;
  clearUser: () => void;
};

const useUserStore = create<UserStoreState>()(
  persist(
    immer((set) => ({
      userID: null,
      isAuthenticated: false,
      setUserID: (userID: string) =>
        set((state) => {
          state.userID = userID;
          state.isAuthenticated = true;
        }),
      clearUser: () =>
        set((state) => {
          state.userID = null;
          state.isAuthenticated = false;
        }),
    })),
    { name: "user-storage" },
  ),
);

export default useUserStore;

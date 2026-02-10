import { create } from "zustand";
import { persist } from "zustand/middleware";
import { immer } from "zustand/middleware/immer";
import type { User } from "../types/types";

type UserStoreState = {
  user: User | null;
  isAuthenticated: boolean;
  setUser: (user: User) => void;
  clearUser: () => void;
};

const useUserStore = create<UserStoreState>()(
  persist(
    immer((set) => ({
      user: null,
      isAuthenticated: false,
      setUser: (user: User) =>
        set((state) => {
          state.user = user;
          state.isAuthenticated = true;
        }),
      clearUser: () =>
        set((state) => {
          state.user = null;
          state.isAuthenticated = false;
        }),
    })),
    { name: "user-storage" },
  ),
);

export default useUserStore;

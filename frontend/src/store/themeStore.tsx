import { createStore } from "zustand/vanilla";
import { persist } from "zustand/middleware";
import { immer } from "zustand/middleware/immer";

type ThemeStoreState = {
  theme: "light" | "dark";
  toggleTheme: () => void;
};

const useThemeStore = createStore<ThemeStoreState>()(
  persist(
    immer((set) => ({
      theme: window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light",
      toggleTheme: () =>
        set((state) => {
          const newTheme = state.theme === "light" ? "dark" : "light";
          state.theme = newTheme;
        }),
    })),
    { name: "theme-storage" },
  ),
);

export default useThemeStore;

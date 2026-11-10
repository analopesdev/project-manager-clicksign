import { ref } from "vue";

const MAX_SEARCH_HISTORY = 5;
const STORAGE_KEY = "project-search-history";

function loadHistory(): string[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    const parsed: unknown = raw ? JSON.parse(raw) : [];
    return Array.isArray(parsed)
      ? parsed
          .filter((item): item is string => typeof item === "string")
          .slice(0, MAX_SEARCH_HISTORY)
      : [];
  } catch {
    return [];
  }
}

function saveHistory(items: string[]) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
  } catch {
    return;
  }
}

const isOpen = ref(false);
const hasProjects = ref(false);
const term = ref("");
const history = ref<string[]>(loadHistory());

function useProjectSearch() {
  function toggle() {
    if (isOpen.value) close();
    else isOpen.value = true;
  }

  function close() {
    remember(term.value);
    isOpen.value = false;
    term.value = "";
  }

  function remember(value: string) {
    const clean = value.trim();
    if (!clean) return;
    history.value = [
      clean,
      ...history.value.filter(item => item !== clean)
    ].slice(0, MAX_SEARCH_HISTORY);
    saveHistory(history.value);
  }

  function forget(value: string) {
    history.value = history.value.filter(item => item !== value);
    saveHistory(history.value);
  }

  return {
    isOpen,
    hasProjects,
    term,
    history,
    toggle,
    close,
    remember,
    forget
  };
}

export { useProjectSearch };

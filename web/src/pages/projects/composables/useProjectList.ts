import axios from "axios";
import { onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useProjectSearch } from "@/composables/useProjectSearch";
import {
  DEFAULT_SORT,
  SORT_OPTIONS,
  deleteProject,
  listProjects,
  setFavorite,
  type Project,
  type ProjectSort
} from "../service";

function parseSort(value: unknown): ProjectSort {
  return SORT_OPTIONS.find(option => option === value) ?? DEFAULT_SORT;
}

export function useProjectList() {
  const route = useRoute();
  const router = useRouter();
  const search = useProjectSearch();

  const favoriteOnly = ref(route.query.favorite === "1");
  const sort = ref<ProjectSort>(parseSort(route.query.sort));

  const initialSearch =
    typeof route.query.q === "string" ? route.query.q.trim() : "";
  if (initialSearch) {
    search.isOpen.value = true;
    search.term.value = initialSearch;
  }

  const projects = ref<Project[]>([]);
  const total = ref(0);
  const loaded = ref(false);
  const loading = ref(false);
  const error = ref<string | null>(null);

  let controller: AbortController | null = null;

  function isSearching() {
    return search.term.value.trim() !== "";
  }

  async function load() {
    controller?.abort();
    const current = new AbortController();
    controller = current;

    loading.value = true;
    error.value = null;
    try {
      const result = await listProjects(
        {
          search: search.term.value,
          favorite: !isSearching() && favoriteOnly.value,
          sort: sort.value
        },
        current.signal
      );
      projects.value = result.projects;
      total.value = result.total;
      loaded.value = true;
      search.hasProjects.value = result.total > 0;
      if (result.total === 0) {
        search.isOpen.value = false;
        search.term.value = "";
      }
    } catch (err) {
      if (axios.isCancel(err)) return;
      error.value = "Não foi possível carregar os projetos. Tente novamente.";
    } finally {
      if (controller === current) loading.value = false;
    }
  }

  async function toggleFavorite(project: Project) {
    const previous = project.favorite;
    project.favorite = !previous;
    try {
      await setFavorite(project.id, project.favorite);
      if (favoriteOnly.value) await load();
    } catch {
      project.favorite = previous;
    }
  }

  const removeTarget = ref<Project | null>(null);
  const removeOpen = ref(false);
  const removing = ref(false);
  const removeError = ref<string | null>(null);

  function askRemove(project: Project) {
    removeTarget.value = project;
    removeError.value = null;
    removeOpen.value = true;
  }

  async function confirmRemove() {
    if (!removeTarget.value || removing.value) return;
    removing.value = true;
    removeError.value = null;
    try {
      await deleteProject(removeTarget.value.id);
      removeOpen.value = false;
      await load();
    } catch {
      removeError.value =
        "Não foi possível remover o projeto. Tente novamente.";
    } finally {
      removing.value = false;
    }
  }

  function syncUrl() {
    const query: Record<string, string> = {};
    const term = search.term.value.trim();
    if (term) query.q = term;
    if (favoriteOnly.value) query.favorite = "1";
    if (sort.value !== DEFAULT_SORT) query.sort = sort.value;
    void router.replace({ query });
  }

  watch([favoriteOnly, sort, search.term], () => {
    syncUrl();
    void load();
  });

  watch(
    () => route.query,
    query => {
      favoriteOnly.value = query.favorite === "1";
      sort.value = parseSort(query.sort);
      const term = typeof query.q === "string" ? query.q.trim() : "";
      if (term === search.term.value.trim()) return;
      search.isOpen.value = term !== "";
      search.term.value = term;
    }
  );

  onMounted(load);
  onUnmounted(() => {
    controller?.abort();
    search.isOpen.value = false;
    search.term.value = "";
    search.hasProjects.value = false;
  });

  return {
    projects,
    total,
    loaded,
    loading,
    error,
    favoriteOnly,
    sort,
    searchTerm: search.term,
    searchOpen: search.isOpen,
    closeSearch: search.close,
    load,
    toggleFavorite,
    remove: {
      target: removeTarget,
      open: removeOpen,
      removing,
      error: removeError,
      ask: askRemove,
      confirm: confirmRemove
    }
  };
}

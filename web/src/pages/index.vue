<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import DeleteProjectDialog from "./projects/components/DeleteProjectDialog.vue";
import PageHeader from "./projects/components/PageHeader.vue";
import ProjectCard from "./projects/components/ProjectCard.vue";
import ProjectsEmpty from "./projects/components/ProjectsEmpty.vue";
import ProjectsFilter from "./projects/components/ProjectsFilter.vue";
import ProjectsSearch from "./projects/components/ProjectsSearch.vue";
import { useProjectList } from "./projects/composables/useProjectList";
import type { Project } from "./projects/service";

const router = useRouter();
const {
  projects,
  total,
  loaded,
  loading,
  error,
  favoriteOnly,
  sort,
  searchTerm,
  searchOpen,
  closeSearch,
  load,
  toggleFavorite,
  remove
} = useProjectList();

const isEmpty = computed(
  () => loaded.value && !error.value && total.value === 0
);
const initialLoading = computed(() => !loaded.value && !error.value);
const isSearching = computed(
  () => searchOpen.value && searchTerm.value.trim() !== ""
);

function goToCreate() {
  void router.push("/projects/new");
}

function goToEdit(project: Project) {
  void router.push(`/projects/${project.id}/edit`);
}
</script>

<template>
  <q-page class="projects-page" :class="{ 'projects-page--empty': isEmpty }">
    <div
      v-if="initialLoading"
      class="state"
      role="status"
      aria-label="Carregando projetos"
    >
      <q-spinner color="primary" size="40px" />
    </div>

    <ProjectsEmpty v-else-if="isEmpty" @create="goToCreate" />

    <template v-else>
      <ProjectsSearch
        :open="searchOpen"
        v-model="searchTerm"
        @close="closeSearch"
      />

      <PageHeader
        v-if="isSearching"
        class="results-header"
        title="Resultado da busca"
        @back="closeSearch"
      />

      <ProjectsFilter
        v-else
        v-model:favorite-only="favoriteOnly"
        v-model:sort="sort"
        :total="loaded ? total : null"
        @create="goToCreate"
      />

      <div v-if="error" class="state" role="alert">
        <p>{{ error }}</p>
        <q-btn color="primary" label="Tentar novamente" no-caps @click="load" />
      </div>

      <p v-else-if="projects.length === 0 && !loading" class="state">
        Nenhum projeto encontrado.
      </p>

      <ul
        v-else
        class="grid"
        :class="{ 'grid--loading': loading }"
        :aria-busy="loading"
      >
        <li v-for="project in projects" :key="project.id" class="grid-item">
          <ProjectCard
            :project="project"
            :search-term="searchTerm"
            @remove="remove.ask"
            @edit="goToEdit"
            @toggle-favorite="toggleFavorite"
          />
        </li>
      </ul>
    </template>

    <DeleteProjectDialog
      v-model="remove.open.value"
      :project="remove.target.value"
      :loading="remove.removing.value"
      :error="remove.error.value"
      @confirm="remove.confirm"
    />
  </q-page>
</template>

<style scoped lang="scss">
.projects-page {
  padding: $space-6 $space-4;

  @media (min-width: $breakpoint-sm-min) {
    padding: $page-padding-top $page-gutter $space-8;
  }
}

.projects-page--empty {
  padding: $page-padding-top $page-gutter;
}

.results-header {
  margin-bottom: $space-5;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: $space-5;
  margin: 0;
  padding: 0;
  list-style: none;
}

.grid-item {
  display: contents;
}

.grid--loading {
  opacity: 0.6;
  transition: opacity $duration-base $ease-out;
}

.state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $space-4;
  padding: 64px 0;
  text-align: center;
}

.sr-only {
  @include sr-only;
}
</style>

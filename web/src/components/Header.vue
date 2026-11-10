<script setup lang="ts">
import logoUrl from "@/assets/logo.png";
import searchUrl from "@/assets/icons/search.svg";
import { useProjectSearch } from "@/composables/useProjectSearch";
import { computed } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const { isOpen, hasProjects, toggle } = useProjectSearch();

const showSearch = computed(() => route.path === "/" && hasProjects.value);
</script>

<template>
  <q-header class="app-header" height-hint="80">
    <q-toolbar class="toolbar">
      <router-link to="/" class="brand">
        <q-img
          fit="contain"
          :src="logoUrl"
          alt=""
          class="brand-mark"
          width="62"
          height="62"
        />
        <span class="brand-name">Gerenciador<br />de Projetos</span>
      </router-link>

      <div v-if="showSearch" class="actions">
        <q-btn
          flat
          round
          dense
          :icon="`img:${searchUrl}`"
          class="search-button"
          aria-label="Buscar projetos"
          :aria-expanded="isOpen"
          @click="toggle"
        />
      </div>
    </q-toolbar>
  </q-header>
</template>

<style scoped lang="scss">
.app-header {
  background: $header-bg;
  box-shadow: $shadow-sm;
}

.toolbar {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  height: $header-height;
  padding-inline: $space-4;

  @media (min-width: $breakpoint-sm-min) {
    padding-inline: $space-8;
  }
}

.brand {
  display: flex;
  grid-column: 2;
  align-items: center;
  gap: $space-3;
  text-decoration: none;
  color: inherit;

  &:focus-visible {
    border-radius: $radius-sm;
    outline: 2px solid currentcolor;
    outline-offset: 4px;
  }
}

.brand-mark {
  width: auto;
  height: 62px;
}

.brand-name {
  font-size: $text-md;
  line-height: 22px;
}

.search-button {
  margin-right: 13px;
  padding: 0;

  :deep(.q-icon) {
    font-size: $icon-md;
  }
}

.actions {
  grid-column: 3;
  justify-self: end;
}
</style>

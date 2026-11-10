<script setup lang="ts">
import plusCircleUrl from "@/assets/icons/plus-circle.svg";
import { ref } from "vue";
import type { ProjectSort } from "../service";

const favoriteOnly = defineModel<boolean>("favoriteOnly", { required: true });
const sort = defineModel<ProjectSort>("sort", { required: true });

defineProps<{
  total: number;
}>();

const emit = defineEmits<{
  create: [];
}>();

const sortOpen = ref(false);

const sortOptions: { label: string; value: ProjectSort }[] = [
  { label: "Ordem alfabética", value: "name" },
  { label: "Iniciados mais recentes", value: "start_date" },
  { label: "Prazo mais próximo", value: "end_date" }
];
</script>

<template>
  <div class="filter">
    <h1 class="heading">
      Projetos
      <span class="count" :aria-label="`${total} projetos`">({{ total }})</span>
    </h1>

    <div class="controls">
      <q-toggle
        v-model="favoriteOnly"
        label="Apenas Favoritos"
        class="favorites-toggle"
      />

      <q-select
        v-model="sort"
        :options="sortOptions"
        emit-value
        map-options
        outlined
        dense
        options-dense
        dropdown-icon="expand_more"
        aria-label="Ordenar projetos"
        class="sort-select"
        :class="{ 'sort-select--open': sortOpen }"
        popup-content-class="sort-menu"
        menu-anchor="bottom left"
        menu-self="top left"
        :menu-offset="[0, 0]"
        @popup-show="sortOpen = true"
        @popup-hide="sortOpen = false"
      />

      <q-btn
        unelevated
        rounded
        no-caps
        color="primary"
        :icon="`img:${plusCircleUrl}`"
        label="Novo projeto"
        class="create-button"
        @click="emit('create')"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
$toggle-width: 48px;
$toggle-height: 24px;
$toggle-thumb: 12px;
$toggle-inset: 6px;

.filter {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: $space-4 $space-5;
  margin-bottom: $space-5;
}

.heading {
  @include text($text-2xl, $weight-semibold, 1);

  display: flex;
  align-items: center;
  gap: $space-2;
  margin: 0;
  color: $heading;
}

.count {
  @include text($text-count, $weight-regular, 1);

  color: $text-label;
}

.controls {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: $space-4 $space-5;
}

.favorites-toggle {
  @include text($text-sm, $weight-regular, 22px);

  color: $text-toggle;

  :deep(.q-toggle__inner) {
    width: $toggle-width;
    min-width: $toggle-width;
    height: $toggle-height;
    padding: 0;
    color: $toggle-off;
  }

  :deep(.q-toggle__inner--truthy) {
    color: $highlight;
  }

  :deep(.q-toggle__label) {
    padding-left: $space-2;
  }

  :deep(.q-toggle__track) {
    width: 100%;
    height: 100%;
    border-radius: $toggle-height * 0.5;
    opacity: 1;
  }

  :deep(.q-toggle__thumb) {
    top: $toggle-inset;
    left: $toggle-inset;
    width: $toggle-thumb;
    height: $toggle-thumb;
  }

  :deep(.q-toggle__inner--truthy .q-toggle__thumb) {
    left: $toggle-width - $toggle-inset - $toggle-thumb;
  }

  :deep(.q-toggle__thumb::after),
  :deep(.q-toggle__inner--truthy .q-toggle__thumb::after) {
    background: $surface;
    box-shadow: none;
  }
}

.sort-select {
  width: $select-width;
  max-width: 100%;
  font-size: $text-sm;

  :deep(.q-field__control) {
    height: $control-height;
    padding: 0 $space-4;
    color: $text-body;
    background: $surface;
    border-radius: $radius-md;

    &::before {
      border: 1px solid $border-input;
      border-radius: $radius-md;
    }
  }

  :deep(.q-field__native) {
    padding: 0;
    line-height: 18px;
    color: $text-body;
  }

  :deep(.q-field__append .q-icon) {
    font-size: $icon-sm;
    color: $text-muted;
  }
}

.sort-select--open {
  :deep(.q-field__control) {
    border-radius: $radius-md $radius-md 0 0;

    &::before {
      border-color: $primary;
      border-radius: $radius-md $radius-md 0 0;
    }

    &::after {
      border-color: $primary;
      border-width: 1px;
      border-radius: $radius-md $radius-md 0 0;
      opacity: 1;
      transform: none;
    }
  }

  :deep(.q-field__append .q-icon) {
    color: $text-label;
  }
}

.create-button {
  @include text($text-sm, $weight-regular, 24px);

  height: $control-height;
  min-height: $control-height;
  padding: 0 $space-5;
  color: $text-inverse;
  background: $primary;

  :deep(.q-icon) {
    font-size: $icon-lg;
  }

  :deep(.q-btn__content) {
    line-height: 24px;
  }
}

@media (max-width: $breakpoint-xs-max) {
  .filter {
    gap: $space-4;
  }

  .controls {
    flex-direction: column;
    align-items: stretch;
    gap: $space-4;
    width: 100%;
  }

  .sort-select {
    width: 100%;
  }

  .create-button {
    width: 100%;
  }
}
</style>

<style lang="scss">
.q-menu.sort-menu {
  width: $select-width;
  overflow: hidden;
  background: $surface;
  border: 1px solid $primary;
  border-top: 0;
  border-radius: 0 0 $radius-lg $radius-lg;
  box-shadow: $shadow-sm;

  .q-item {
    @include text($text-sm, $weight-regular, 16px);

    min-height: 54px;
    padding: 0 $space-4;
    color: $text-body;
    border-bottom: 1px solid $divider;

    &:last-child {
      border-bottom: 0;
    }
  }

  .q-focus-helper {
    display: none;
  }

  .q-item--active,
  .q-item.q-manual-focusable--focused {
    color: $text-body;
    background: transparent;
  }

  .q-item:hover {
    background: $hover-tint;
  }
}
</style>

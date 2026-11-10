<script setup lang="ts">
import { computed } from "vue";
import type { Project } from "../service";
import coverPlaceholderUrl from "@/assets/bg-project-cover.png";
import calendarDayUrl from "@/assets/icons/calendar-day-light.svg";
import calendarCheckUrl from "@/assets/icons/calendar-check-light.svg";
import editUrl from "@/assets/icons/edit.svg";
import trashUrl from "@/assets/icons/trash.svg";
import HighlightText from "./HighlightText.vue";

const longDateFormatter = new Intl.DateTimeFormat("pt-BR", {
  day: "numeric",
  month: "long",
  year: "numeric"
});

function formatLongDate(isoDate: string): string {
  const [year, month, day] = isoDate.split("-").map(Number);
  if (!year || !month || !day) return isoDate;
  return longDateFormatter.format(new Date(year, month - 1, day));
}

const props = defineProps<{
  project: Project;
  searchTerm?: string;
}>();

const emit = defineEmits<{
  toggleFavorite: [project: Project];
  edit: [project: Project];
  remove: [project: Project];
}>();

const favoriteLabel = computed(() =>
  props.project.favorite
    ? `Remover ${props.project.name} dos favoritos`
    : `Adicionar ${props.project.name} aos favoritos`
);
</script>

<template>
  <q-card tag="article" flat class="project-card">
    <div class="cover">
      <q-img
        v-if="project.coverUrl"
        :src="project.coverUrl"
        :alt="`Capa do projeto ${project.name}`"
        :ratio="3 / 2"
        loading="lazy"
        class="cover-image"
      />
      <div v-else class="cover-placeholder">
        <q-img
          :src="coverPlaceholderUrl"
          alt="Capa do projeto"
          class="cover-placeholder-image"
        />
      </div>

      <div class="actions">
        <q-btn
          flat
          round
          dense
          size="md"
          :icon="project.favorite ? 'star' : 'star_border'"
          class="favorite"
          :class="{ 'is-favorite': project.favorite }"
          :aria-label="favoriteLabel"
          :aria-pressed="project.favorite"
          @click="emit('toggleFavorite', project)"
        />

        <q-btn
          round
          unelevated
          dense
          size="sm"
          icon="more_horiz"
          class="menu-trigger"
          :aria-label="`Ações do projeto ${project.name}`"
        >
          <q-menu
            anchor="bottom right"
            self="top right"
            :offset="[0, 14]"
            class="project-menu"
          >
            <q-list class="menu-list">
              <q-item
                v-close-popup
                clickable
                class="menu-item"
                @click="emit('edit', project)"
              >
                <img :src="editUrl" alt="" width="24" height="24" />
                <span>Editar</span>
              </q-item>
              <q-item
                v-close-popup
                clickable
                class="menu-item"
                @click="emit('remove', project)"
              >
                <img :src="trashUrl" alt="" width="24" height="24" />
                <span>Remover</span>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>
    </div>

    <h2 class="title">
      <HighlightText :text="project.name" :term="searchTerm" />
    </h2>

    <p class="client">
      <strong>Cliente:</strong>
      <HighlightText :text="project.client" :term="searchTerm" />
    </p>

    <ul class="dates">
      <li class="date">
        <img
          :src="calendarDayUrl"
          alt=""
          width="24"
          height="24"
          class="date-icon"
        />
        <span class="sr-only">Início:</span>
        <time :datetime="project.startDate">{{
          formatLongDate(project.startDate)
        }}</time>
      </li>
      <li class="date">
        <img
          :src="calendarCheckUrl"
          alt=""
          width="24"
          height="24"
          class="date-icon"
        />
        <span class="sr-only">Término:</span>
        <time :datetime="project.endDate">{{
          formatLongDate(project.endDate)
        }}</time>
      </li>
    </ul>
  </q-card>
</template>

<style scoped lang="scss">
.project-card {
  display: grid;
  grid-template-rows: subgrid;
  grid-row: span 4;
  row-gap: 0;
  height: 100%;
  border: 1px solid $border;
  border-radius: $radius-lg;
}

.cover {
  position: relative;
  aspect-ratio: 3 / 2;
  margin: -1px -1px 0;
  overflow: hidden;
  background: $primary;
  border-bottom: 1px solid $border;
  border-radius: $radius-lg $radius-lg 0 0;
}

.cover-image {
  height: 100%;
}

.cover-placeholder {
  display: grid;
  place-items: center;
  height: 100%;
  color: $accent;
}

.actions {
  position: absolute;
  right: $space-4;
  bottom: $space-4;
  display: flex;
  align-items: center;
  gap: $space-3;
}

.favorite {
  color: $text-inverse;

  &.is-favorite :deep(.q-icon) {
    background: $gradient-favorite;
    background-clip: text;
    -webkit-text-fill-color: transparent;
  }
}

.menu-trigger {
  color: $primary;
  background: white;
}

.title,
.client,
.dates {
  margin: 0;
  padding-inline: $space-5;
}

.title {
  @include text($text-lg, $weight-bold, 25px);

  padding-top: 20px;
  overflow-wrap: anywhere;
  color: $heading;
}

.client {
  @include text($text-sm, $weight-regular, 20px);

  padding-top: $space-1;
  color: $text-muted;

  &::after {
    display: block;
    height: 1px;
    margin-top: $space-4;
    content: "";
    background: $divider;
  }

  strong {
    margin-right: $space-1;
    font-weight: $weight-bold;
  }
}

.dates {
  display: flex;
  flex-direction: column;
  gap: $space-3;
  padding-top: $space-4;
  padding-bottom: $space-5;
  list-style: none;
}

.date {
  @include text($text-sm, $weight-regular, 20px);

  display: flex;
  align-items: center;
  gap: $space-4;
  color: $text-muted;
}

.sr-only {
  @include sr-only;
}
</style>

<style lang="scss">
$menu-width: 240px;
$menu-item-height: 50px;
$menu-arrow-width: 16px;
$menu-arrow-height: 12px;
$menu-arrow-right: 7px;

.q-menu.project-menu {
  width: $menu-width;
  overflow: visible;
  background: $surface;
  border-radius: $radius-md;
  box-shadow: $shadow-sm;

  &::before {
    position: absolute;
    top: -$menu-arrow-height;
    right: $menu-arrow-right;
    width: $menu-arrow-width;
    height: $menu-arrow-height;
    content: "";
    background: $surface;
    clip-path: polygon(50% 0, 100% 100%, 0 100%);
  }

  .menu-item {
    @include text($text-sm, $weight-regular, 16px);

    align-items: center;
    gap: $space-4;
    min-height: $menu-item-height;
    padding: 0 20px;
    color: $text-label;

    & + .menu-item {
      border-top: 1px solid $divider-soft;
    }

    &:hover {
      background: $hover-tint;
    }

    .q-focus-helper {
      display: none;
    }
  }
}
</style>

<script setup lang="ts">
import searchUrl from "@/assets/icons/search-purple.svg";
import historyUrl from "@/assets/icons/history.svg";
import { useProjectSearch } from "@/composables/useProjectSearch";

const SEARCH_DEBOUNCE_MS = 300;

defineProps<{
  open: boolean;
}>();

const term = defineModel<string>({ required: true });

const emit = defineEmits<{
  close: [];
}>();

const { history, remember, forget } = useProjectSearch();

function submit() {
  remember(term.value);
}

function pick(value: string) {
  term.value = value;
  remember(value);
}
</script>

<template>
  <Transition name="fade" appear>
    <div v-if="open && !term" class="backdrop" @click="emit('close')" />
  </Transition>

  <Transition name="drop" appear>
    <div
      v-if="open"
      class="search"
      :class="{ 'search--results': term }"
      role="search"
      @keydown.esc="emit('close')"
    >
      <div class="field">
        <img :src="searchUrl" alt="" width="20" height="20" class="icon" />
        <q-input
          v-model.trim="term"
          :debounce="SEARCH_DEBOUNCE_MS"
          borderless
          dense
          autofocus
          placeholder="Digite o nome do projeto..."
          aria-label="Buscar projeto pelo nome"
          class="input"
          @keydown.enter="submit"
        />
      </div>

      <ul
        v-if="history.length && !term"
        class="history"
        aria-label="Buscas recentes"
      >
        <li v-for="item in history" :key="item" class="history-item">
          <button type="button" class="history-pick" @click="pick(item)">
            <img
              :src="historyUrl"
              alt=""
              width="16"
              height="16"
              class="history-icon"
            />
            <span>{{ item }}</span>
          </button>
          <q-btn
            flat
            round
            dense
            size="sm"
            icon="close"
            :aria-label="`Remover ${item} do histórico`"
            class="history-remove"
            @click="forget(item)"
          />
        </li>
      </ul>
    </div>
  </Transition>
</template>

<style scoped lang="scss">
.backdrop {
  position: fixed;
  inset: 0;
  z-index: 2999;
}

.search {
  position: fixed;
  top: 0;
  right: 0;
  left: 0;
  z-index: 3000;
  overflow: hidden;
  background: $surface;
  border: 2px solid $primary;
  border-radius: $radius-lg;
  box-shadow: $shadow-sm;
}

.search--results {
  border-color: transparent;
  border-radius: 0;
}

.drop-enter-active {
  transition:
    transform $duration-base $ease-out,
    opacity $duration-base $ease-out;
}

.drop-leave-active {
  transition:
    transform $duration-fast $ease-in,
    opacity $duration-fast $ease-in;

  .history {
    display: none;
  }
}

.drop-enter-from,
.drop-leave-to {
  opacity: 0;
  transform: translateY(-$space-4);
}

.fade-enter-active {
  transition: opacity $duration-base $ease-out;
}

.fade-leave-active {
  transition: opacity $duration-fast $ease-in;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (prefers-reduced-motion: reduce) {
  .drop-enter-active,
  .drop-leave-active,
  .fade-enter-active,
  .fade-leave-active {
    transition: none;
  }
}

.field {
  display: flex;
  align-items: center;
  gap: $space-4;
  height: $header-height;
  padding: 0 $space-5;
}

.icon {
  flex: none;
}

.input {
  @include text($text-sm);

  flex: 1;
  color: $text-body;

  :deep(input::placeholder) {
    color: $text-placeholder;
    opacity: 1;
  }
}

.history {
  margin: 0;
  padding: 0;
  border-top: 1px solid $divider-soft;
  list-style: none;
}

.history-item {
  display: flex;
  align-items: center;
  height: $button-height-lg;
  padding-right: 20px;
  color: $text-muted;
  border-bottom: 1px solid $divider-soft;

  &:last-child {
    border-bottom: 0;
  }
}

.history-pick {
  @include text($text-sm, $weight-regular, 16px);

  display: flex;
  flex: 1;
  align-items: center;
  gap: $space-3;
  height: 100%;
  padding: 0 0 0 $space-5;
  text-align: left;
  color: inherit;
  background: none;
  border: 0;
  cursor: pointer;

  &:hover,
  &:focus-visible {
    background: $hover-tint;
  }
}

.history-icon {
  flex: none;
}

.history-remove {
  color: $text-muted;
}
</style>

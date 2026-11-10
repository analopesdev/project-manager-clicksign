<script setup lang="ts">
defineProps<{
  title: string;
}>();

const emit = defineEmits<{
  back: [];
}>();
</script>

<template>
  <header>
    <button type="button" class="back" @click="emit('back')">
      <q-icon name="arrow_back" size="16px" aria-hidden="true" />
      Voltar
    </button>
    <h1 class="title">{{ title }}</h1>
  </header>
</template>

<style scoped lang="scss">
.back {
  @include button-reset;
  @include text($text-sm);

  display: flex;
  align-items: center;
  gap: $space-2;
  height: $control-height;
  color: $text-label;

  &:hover {
    text-decoration: underline;
  }

  &:focus-visible {
    border-radius: $radius-sm;
    outline: 2px solid $primary;
    outline-offset: $space-1;
  }
}

.title {
  @include page-title;

  margin-top: $space-1;
}
</style>

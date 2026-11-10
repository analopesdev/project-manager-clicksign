<script setup lang="ts">
import plusCircleUrl from "@/assets/icons/plus-circle.svg";

const emit = defineEmits<{
  create: [];
}>();
</script>

<template>
  <section class="empty">
    <h1 class="title">Nenhum projeto</h1>
    <p class="description">
      Clique no botão abaixo para criar o primeiro e gerenciá-lo.
    </p>
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
  </section>
</template>

<style scoped lang="scss">
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - #{$header-height} - 120px);
  padding: $space-6 $space-4;
  text-align: center;
  background: $surface;
  border-radius: $radius-sm;
}

.title {
  @include text($text-2xl, $weight-semibold, 29px);

  margin: 0;
  color: $heading;
}

.description {
  @include text($text-sm, $weight-regular, 24px);

  margin: $space-5 0 $space-6;
  color: $text-muted;
}

.create-button {
  @include button-large;
  @include button-primary;

  padding: 0 36px;

  :deep(.q-icon) {
    font-size: $icon-lg;
  }
}
</style>

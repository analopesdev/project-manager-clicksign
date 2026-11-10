<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
  text: string;
  term?: string | undefined;
}>();

const parts = computed(() => {
  const term = props.term?.trim();
  if (!term) return [props.text];
  const escaped = term.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  return props.text.split(new RegExp(`(${escaped})`, "i"));
});
</script>

<template>
  <template v-for="(part, index) in parts" :key="index">
    <mark v-if="index % 2" class="highlight">{{ part }}</mark>
    <template v-else>{{ part }}</template>
  </template>
</template>

<style scoped lang="scss">
.highlight {
  color: $text-inverse;
  background: $highlight;
}
</style>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from "vue";
import uploadUrl from "@/assets/icons/upload.svg";
import trashUrl from "@/assets/icons/trash.svg";
import { coverFileSchema } from "../schema";

const file = defineModel<File | null>("file", { required: true });
const url = defineModel<string>("url", { required: true });

const input = ref<HTMLInputElement | null>(null);
const error = ref<string | null>(null);
const dragging = ref(false);

const preview = ref("");
watch(
  file,
  (next, _prev, onCleanup) => {
    if (!next) {
      preview.value = "";
      return;
    }
    const objectUrl = URL.createObjectURL(next);
    preview.value = objectUrl;
    onCleanup(() => URL.revokeObjectURL(objectUrl));
  },
  { immediate: true }
);
onBeforeUnmount(() => (preview.value = ""));

const shownUrl = computed(() => preview.value || url.value);

function pick(selected: File | undefined) {
  if (!selected) return;
  const result = coverFileSchema.safeParse(selected);
  if (!result.success) {
    error.value = result.error.issues[0]?.message ?? "Arquivo inválido";
    return;
  }
  error.value = null;
  file.value = selected;
}

function onChange(event: Event) {
  const target = event.target as HTMLInputElement;
  pick(target.files?.[0]);
  target.value = "";
}

function onDrop(event: DragEvent) {
  dragging.value = false;
  pick(event.dataTransfer?.files?.[0]);
}

function remove() {
  file.value = null;
  url.value = "";
  error.value = null;
}
</script>

<template>
  <div class="cover">
    <input
      ref="input"
      type="file"
      accept=".jpg,.jpeg,.png,image/jpeg,image/png"
      class="native"
      tabindex="-1"
      @change="onChange"
    />

    <div v-if="shownUrl" class="preview">
      <img :src="shownUrl" alt="Capa do projeto" class="image" />
      <button
        type="button"
        class="remove"
        aria-label="Remover capa"
        @click="remove"
      >
        <img :src="trashUrl" alt="" width="20" height="20" />
      </button>
    </div>

    <div
      v-else
      class="dropzone"
      :class="{ 'dropzone--active': dragging }"
      @dragover.prevent="dragging = true"
      @dragleave="dragging = false"
      @drop.prevent="onDrop"
    >
      <img :src="uploadUrl" alt="" width="24" height="24" />
      <p class="hint">Escolha uma imagem .jpg ou .png no seu dispositivo</p>
      <q-btn
        outline
        rounded
        no-caps
        label="Selecionar"
        class="select"
        @click="input?.click()"
      />
    </div>

    <p v-if="error" class="error" role="alert">{{ error }}</p>
  </div>
</template>

<style scoped lang="scss">
$dropzone-height: 174px;
$preview-height: 395px;
$select-button-width: 144px;

.native {
  position: absolute;
  width: 1px;
  height: 1px;
  opacity: 0;
  pointer-events: none;
}

.dropzone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: $space-4;
  height: $dropzone-height;
  padding: $space-4;
  border: 1px dashed $border-input;
  border-radius: $radius-sm;
}

.dropzone--active {
  background: $hover-tint;
  border-color: $primary;
}

.hint {
  @include text($text-sm, $weight-regular, 22px);

  margin: 0;
  text-align: center;
  color: $text-muted;
}

.select {
  @include text($text-sm, $weight-regular);

  width: $select-button-width;
  height: $control-height;
  color: $text-label;
  background: $surface;
}

.preview {
  position: relative;
  height: $preview-height;
}

.image {
  display: block;
  width: 100%;
  height: 100%;
  border-radius: $radius-md;
  object-fit: cover;
}

.remove {
  position: absolute;
  top: $space-5;
  right: $space-5;
  display: grid;
  place-items: center;
  width: 32px;
  height: 32px;
  padding: 0;
  background: $surface;
  border: 0;
  border-radius: 50%;
  box-shadow: $shadow-sm;
  cursor: pointer;
}

.error {
  @include text($text-xs, $weight-regular, 22px);

  margin: $space-1 0 0;
  color: $danger;
}
</style>

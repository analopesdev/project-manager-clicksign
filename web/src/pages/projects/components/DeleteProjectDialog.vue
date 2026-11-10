<script setup lang="ts">
import trashUrl from "@/assets/icons/trash-white.svg";
import type { Project } from "../service";

const open = defineModel<boolean>({ required: true });

defineProps<{
  project: Project | null;
  loading?: boolean;
  error?: string | null;
}>();

const emit = defineEmits<{
  confirm: [];
}>();
</script>

<template>
  <q-dialog
    v-model="open"
    :persistent="loading"
    backdrop-filter="blur(5px)"
    class="delete-dialog"
  >
    <div
      class="dialog"
      role="alertdialog"
      aria-labelledby="delete-title"
      aria-describedby="delete-description"
    >
      <div class="badge" aria-hidden="true">
        <img :src="trashUrl" alt="" width="20" height="20" />
      </div>

      <h2 id="delete-title" class="title">Remover projeto</h2>

      <p id="delete-description" class="description">
        Essa ação removerá definitivamente o projeto:
      </p>
      <p class="name">{{ project?.name }}</p>

      <p v-if="error" class="error" role="alert">{{ error }}</p>

      <div class="actions">
        <q-btn
          outline
          rounded
          no-caps
          label="Cancelar"
          class="button cancel"
          :disable="loading"
          @click="open = false"
        />
        <q-btn
          unelevated
          rounded
          no-caps
          label="Confirmar"
          class="button confirm"
          :loading="loading"
          @click="emit('confirm')"
        />
      </div>
    </div>
  </q-dialog>
</template>

<style scoped lang="scss">
$dialog-width: 582px;
$dialog-min-height: 321px;
$dialog-offset-y: -87.5px;
$badge-size: 64px;
$cancel-width: 226px;

.dialog {
  position: relative;
  box-sizing: border-box;
  width: $dialog-width;
  max-width: calc(100vw - 32px);
  min-height: $dialog-min-height;
  padding: 0 $space-6 $space-6;
  text-align: center;
  background: $surface;
  border-radius: $radius-md;
  box-shadow: $shadow-sm;
  transform: translateY($dialog-offset-y);
}

.badge {
  position: absolute;
  top: -($badge-size * 0.5);
  left: 50%;
  display: grid;
  place-items: center;
  width: $badge-size;
  height: $badge-size;
  color: $text-inverse;
  background: $primary;
  border-radius: 50%;
  box-shadow: $shadow-sm;
  transform: translateX(-50%);
}

.title {
  @include text($text-xl, $weight-semibold, 32px);

  margin: 0;
  padding: $space-8 0 0;
  color: $heading;
}

.description {
  @include text($text-sm, $weight-regular, 22px);

  margin: 0;
  margin-top: $space-5;
  padding: $space-6 0 0;
  color: $text-muted;
  border-top: 1px solid $divider-strong;
}

.name {
  @include text($text-2xl, $weight-medium, 32px);

  margin: 13px 0 0;
  overflow-wrap: anywhere;
  color: $text-body;
}

.error {
  margin: $space-3 0 0;
  font-size: $text-xs;
  color: $danger;
}

.actions {
  display: flex;
  gap: $space-6;
  margin-top: $space-6;
}

.button {
  @include button-large;

  color: $text-label;
}

.cancel {
  flex: 0 0 $cancel-width;
}

@media (max-width: $breakpoint-xs-max) {
  .dialog {
    padding: 0 20px $space-5;
  }

  .actions {
    gap: $space-3;
  }

  .cancel {
    flex: 1;
  }

  .button {
    font-size: $text-sm;
  }
}

.confirm {
  @include button-primary;

  flex: 1;
}
</style>

<style lang="scss">
.delete-dialog .q-dialog__backdrop {
  background: $overlay;
}
</style>

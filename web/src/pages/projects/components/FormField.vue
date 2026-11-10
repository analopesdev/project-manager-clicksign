<script setup lang="ts">
import { computed, ref } from "vue";

const model = defineModel<string>();

const props = defineProps<{
  id: string;
  label: string;
  required?: boolean;
  error?: string | undefined;
  maxlength?: number;
  dateIcon?: string;
  dateIconLabel?: string;
}>();

const emit = defineEmits<{
  blur: [];
}>();

const popup = ref<{ hide: () => void } | null>(null);
const errorId = computed(() => `${props.id}-error`);
const describedBy = computed(() => (props.error ? errorId.value : undefined));
const iconMask = computed(() => `url("${props.dateIcon}")`);
</script>

<template>
  <div class="field" :class="{ 'field--invalid': error }">
    <label :for="id" class="label">
      {{ label }}
      <span v-if="required" class="required">(Obrigatório)</span>
    </label>

    <slot>
      <div class="control" :class="{ 'control--date': dateIcon }">
        <q-input
          v-if="dateIcon"
          v-model="model"
          :for="id"
          mask="##/##/####"
          placeholder="dd/mm/aaaa"
          borderless
          dense
          class="input"
          :aria-invalid="!!error"
          :aria-describedby="describedBy"
          @blur="emit('blur')"
        />
        <input
          v-else
          :id="id"
          v-model="model"
          class="input"
          type="text"
          autocomplete="off"
          :maxlength="maxlength"
          :aria-invalid="!!error"
          :aria-describedby="describedBy"
          @blur="emit('blur')"
        />

        <button
          v-if="dateIcon"
          type="button"
          class="icon-button"
          :aria-label="dateIconLabel"
        >
          <span
            class="icon"
            :style="{ '--icon': iconMask }"
            aria-hidden="true"
          />
          <q-popup-proxy
            ref="popup"
            cover
            transition-show="scale"
            transition-hide="scale"
          >
            <q-date
              v-model="model"
              mask="DD/MM/YYYY"
              minimal
              color="primary"
              @update:model-value="popup?.hide()"
            />
          </q-popup-proxy>
        </button>
      </div>
    </slot>

    <p v-if="error" :id="errorId" class="error" role="alert">{{ error }}</p>
  </div>
</template>

<style scoped lang="scss">
.field {
  display: flex;
  flex-direction: column;
  gap: $space-1;
  min-width: 0;
}

.label {
  @include text($text-md, $weight-medium, 22px);

  display: flex;
  align-items: flex-end;
  gap: $space-2;
  color: $text-label;
}

.required {
  @include text($text-xs, $weight-regular, 22px);

  color: $text-muted;
}

.error {
  @include text($text-xs, $weight-regular, 22px);

  margin: 0;
  color: $danger;
}

.control {
  display: flex;
  align-items: center;
  height: $control-height;
  padding: 0 $space-4;
  background: $surface;
  border: 1px solid $border-input;
  border-radius: $radius-md;

  &:focus-within {
    border-color: $primary;
  }
}

.control--date {
  padding-right: $space-2;
}

.input {
  @include text($text-sm, $weight-regular, 22px);

  flex: 1;
  width: 100%;
  min-width: 0;
  padding: 0;
  color: $text-body;
  background: none;
  border: 0;
  outline: 0;

  :deep(.q-field__control),
  :deep(.q-field__native) {
    height: $control-height - 2px;
    min-height: 0;
    padding: 0;
    color: inherit;
  }

  :deep(input::placeholder) {
    color: $text-placeholder;
    opacity: 1;
  }
}

.icon-button {
  @include button-reset;

  display: grid;
  flex: none;
  place-items: center;
  width: 32px;
  height: 32px;
  color: $text-muted;
  border-radius: $radius-sm;
}

.icon {
  width: $icon-lg;
  height: $icon-lg;
  background-color: currentcolor;
  mask: var(--icon) center / contain no-repeat;
}

.field--invalid {
  .label {
    color: $danger-strong;
  }

  .required,
  .control .input,
  .icon-button {
    color: $danger;
  }

  .control,
  .control:focus-within {
    border-color: $danger;
  }
}
</style>

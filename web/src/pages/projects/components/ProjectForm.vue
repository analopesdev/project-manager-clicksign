<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import calendarDayUrl from "@/assets/icons/calendar-day-light.svg";
import calendarCheckUrl from "@/assets/icons/calendar-check-light.svg";
import {
  emptyValues,
  projectFormSchema,
  type FieldName,
  type ProjectFormValues
} from "../schema";
import { MAX_TEXT_LENGTH } from "../schema";
import CoverUpload from "./CoverUpload.vue";
import FormField from "./FormField.vue";

const props = defineProps<{
  initial?: ProjectFormValues | undefined;
  submitting?: boolean;
  serverErrors?: Partial<Record<FieldName, string>>;
  generalError?: string | null;
}>();

const emit = defineEmits<{
  submit: [values: ProjectFormValues, coverFile: File | null];
}>();

const values = reactive<ProjectFormValues>(props.initial ?? emptyValues());
watch(
  () => props.initial,
  next => {
    if (next) Object.assign(values, next);
  }
);

const coverFile = ref<File | null>(null);
const touched = reactive<Record<FieldName, boolean>>({
  name: false,
  client: false,
  startDate: false,
  endDate: false,
  coverUrl: false
});
const submitted = ref(false);

const result = computed(() => projectFormSchema.safeParse(values));
const errors = computed(() => {
  const found: Partial<Record<FieldName, string>> = {};
  if (!result.value.success) {
    for (const issue of result.value.error.issues) {
      const field = issue.path[0] as FieldName | undefined;
      if (field && !found[field]) found[field] = issue.message;
    }
  }
  return found;
});

const dismissedServer = reactive<Partial<Record<FieldName, boolean>>>({});
watch(
  () => props.serverErrors,
  () =>
    Object.keys(dismissedServer).forEach(
      k => delete dismissedServer[k as FieldName]
    )
);
for (const field of ["name", "client", "startDate", "endDate"] as const) {
  watch(
    () => values[field],
    () => (dismissedServer[field] = true)
  );
}

function errorFor(field: FieldName): string | undefined {
  if (touched[field] || submitted.value) {
    if (errors.value[field]) return errors.value[field];
  }
  return dismissedServer[field] ? undefined : props.serverErrors?.[field];
}

const canSubmit = computed(() => result.value.success && !props.submitting);

function touch(field: FieldName) {
  touched[field] = true;
}

function onSubmit() {
  if (props.submitting) return;
  submitted.value = true;
  if (!result.value.success) return;
  emit(
    "submit",
    { ...result.value.data, coverUrl: values.coverUrl },
    coverFile.value
  );
}
</script>

<template>
  <form class="form" novalidate @submit.prevent="onSubmit">
    <FormField
      id="name"
      v-model="values.name"
      label="Nome do projeto"
      required
      :maxlength="MAX_TEXT_LENGTH"
      :error="errorFor('name')"
      @blur="touch('name')"
    />

    <FormField
      id="client"
      v-model="values.client"
      label="Cliente"
      required
      :maxlength="MAX_TEXT_LENGTH"
      :error="errorFor('client')"
      @blur="touch('client')"
    />

    <div class="dates">
      <FormField
        id="startDate"
        v-model="values.startDate"
        label="Data de Início"
        required
        :date-icon="calendarDayUrl"
        date-icon-label="Escolher a data de início"
        :error="errorFor('startDate')"
        @blur="touch('startDate')"
      />

      <FormField
        id="endDate"
        v-model="values.endDate"
        label="Data Final"
        required
        :date-icon="calendarCheckUrl"
        date-icon-label="Escolher a data final"
        :error="errorFor('endDate')"
        @blur="touch('endDate')"
      />
    </div>

    <FormField id="cover" label="Capa do projeto" :error="errorFor('coverUrl')">
      <CoverUpload v-model:file="coverFile" v-model:url="values.coverUrl" />
    </FormField>

    <p v-if="generalError" class="general-error" role="alert">
      {{ generalError }}
    </p>

    <q-btn
      type="submit"
      unelevated
      rounded
      no-caps
      label="Salvar projeto"
      class="save"
      :class="{ 'save--disabled': !canSubmit }"
      :loading="submitting"
      :aria-disabled="!canSubmit"
    />
  </form>
</template>

<style scoped lang="scss">
$form-width: 704px;

.form {
  display: flex;
  flex-direction: column;
  gap: $space-6;
  width: 100%;
  max-width: $form-width;
}

.dates {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: $space-7;

  @media (max-width: $breakpoint-xs-max) {
    grid-template-columns: 1fr;
    gap: $space-6;
  }
}

.general-error {
  @include text($text-xs);

  margin: -16px 0 0;
  color: $danger;
}

.save {
  @include button-large;
  @include button-primary;
}

.save--disabled {
  background: $primary-disabled;
}
</style>

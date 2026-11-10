<script setup lang="ts">
import { isAxiosError } from "axios";
import { computed, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import PageHeader from "./components/PageHeader.vue";
import ProjectForm from "./components/ProjectForm.vue";
import {
  brToIso,
  projectToFormValues,
  type FieldName,
  type ProjectFormValues
} from "./schema";
import {
  ApiValidationError,
  createProject,
  getProject,
  updateProject,
  uploadCover
} from "./service";

definePage({ path: "/projects/new", alias: "/projects/:id/edit" });

const route = useRoute();
const router = useRouter();

const id = computed(() => {
  const param = (route.params as { id?: string }).id;
  return typeof param === "string" ? param : undefined;
});

const initial = ref<ProjectFormValues>();
const loading = ref(false);
const loadFailure = ref<"not-found" | "failed" | null>(null);
const loadMessage = computed(() =>
  loadFailure.value === "not-found"
    ? "Projeto não encontrado."
    : "Não foi possível carregar o projeto. Tente novamente."
);

const saving = ref(false);
const serverErrors = ref<Partial<Record<FieldName, string>>>({});
const generalError = ref<string | null>(null);

watch(
  id,
  async current => {
    initial.value = undefined;
    loadFailure.value = null;
    if (current === undefined) return;
    loading.value = true;
    try {
      initial.value = projectToFormValues(await getProject(current));
    } catch (error) {
      const notFound = isAxiosError(error) && error.response?.status === 404;
      loadFailure.value = notFound ? "not-found" : "failed";
    } finally {
      loading.value = false;
    }
  },
  { immediate: true }
);

async function save(values: ProjectFormValues, coverFile: File | null) {
  if (saving.value) return;
  saving.value = true;
  serverErrors.value = {};
  generalError.value = null;
  try {
    const input = {
      name: values.name,
      client: values.client,
      startDate: brToIso(values.startDate),
      endDate: brToIso(values.endDate),
      coverUrl: coverFile ? await uploadCover(coverFile) : values.coverUrl
    };
    await (id.value === undefined
      ? createProject(input)
      : updateProject(id.value, input));
    await router.push("/");
  } catch (error) {
    if (error instanceof ApiValidationError) {
      serverErrors.value = error.fields;
    } else {
      generalError.value =
        "Não foi possível salvar o projeto. Tente novamente.";
    }
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <q-page class="form-page">
    <PageHeader
      :title="id === undefined ? 'Novo projeto' : 'Editar projeto'"
      @back="router.push('/')"
    />

    <section class="panel">
      <q-spinner v-if="loading" color="primary" size="40px" />
      <p v-else-if="loadFailure" class="state" role="alert">{{
        loadMessage
      }}</p>
      <ProjectForm
        v-else
        :key="id ?? 'new'"
        :initial="initial"
        :submitting="saving"
        :server-errors="serverErrors"
        :general-error="generalError"
        @submit="save"
      />
    </section>
  </q-page>
</template>

<style scoped lang="scss">
$panel-offset: 33px;
$panel-padding-y: 52px;

.form-page {
  padding: $space-6 $space-4;

  @media (min-width: $breakpoint-sm-min) {
    padding: $page-padding-top $page-gutter $space-8;
  }
}

.panel {
  display: flex;
  justify-content: center;
  margin-top: $panel-offset;
  padding: $space-6 $space-4;
  border: 1px solid $border;
  border-radius: $radius-md;

  @media (min-width: $breakpoint-sm-min) {
    padding: $panel-padding-y $space-4;
  }
}

.state {
  margin: 0;
  color: $text-muted;
}
</style>

import axios, { isAxiosError } from "axios";

type Project = {
  id: string;
  name: string;
  client: string;
  startDate: string;
  endDate: string;
  coverUrl?: string | undefined;
  favorite: boolean;
};

const SORT_OPTIONS = ["name", "start_date", "end_date"] as const;
const DEFAULT_SORT = "name";
type ProjectSort = (typeof SORT_OPTIONS)[number];

type ListProjectsParams = {
  search?: string;
  favorite?: boolean;
  sort?: ProjectSort;
};

type ProjectList = {
  projects: Project[];
  total: number;
};

type ProjectResponse = {
  id: string;
  name: string;
  client: string;
  start_date: string;
  end_date: string;
  cover_url: string;
  favorite: boolean;
};

type ProjectListResponse = {
  data: ProjectResponse[];
  total: number;
};

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? "/api"
});

function toProject(item: ProjectResponse): Project {
  return {
    id: item.id,
    name: item.name,
    client: item.client,
    startDate: item.start_date,
    endDate: item.end_date,
    coverUrl: item.cover_url || undefined,
    favorite: item.favorite
  };
}

async function listProjects(
  params: ListProjectsParams = {},
  signal?: AbortSignal
): Promise<ProjectList> {
  const { data: body } = await api.get<ProjectListResponse>("/projects", {
    params: {
      search: params.search || undefined,
      favorite: params.favorite ? true : undefined,
      sort: params.sort
    },
    ...(signal && { signal })
  });

  return {
    projects: body.data.map(toProject),
    total: body.total
  };
}

type ProjectInput = {
  name: string;
  client: string;
  startDate: string;
  endDate: string;
  coverUrl: string;
};

class ApiValidationError extends Error {
  fields: Partial<
    Record<"name" | "client" | "startDate" | "endDate" | "coverUrl", string>
  >;

  constructor(fields: ApiValidationError["fields"]) {
    super("validation failed");
    this.fields = fields;
  }
}

const FIELD_NAMES = {
  name: "name",
  client: "client",
  start_date: "startDate",
  end_date: "endDate",
  cover_url: "coverUrl"
} as const;

function toBody(input: ProjectInput) {
  return {
    name: input.name.trim(),
    client: input.client.trim(),
    start_date: input.startDate,
    end_date: input.endDate,
    cover_url: input.coverUrl
  };
}

function rethrowValidation(error: unknown): never {
  if (isAxiosError(error) && error.response?.status === 422) {
    const raw =
      (error.response.data as { fields?: Record<string, string> }).fields ?? {};
    const fields: ApiValidationError["fields"] = {};
    for (const [key, message] of Object.entries(raw)) {
      const field = FIELD_NAMES[key as keyof typeof FIELD_NAMES];
      if (field) fields[field] = message;
    }
    if (Object.keys(fields).length > 0) throw new ApiValidationError(fields);
  }
  throw error;
}

async function getProject(id: string, signal?: AbortSignal): Promise<Project> {
  const { data } = await api.get<ProjectResponse>(
    `/projects/${id}`,
    signal ? { signal } : {}
  );
  return toProject(data);
}

async function createProject(input: ProjectInput): Promise<Project> {
  try {
    const { data } = await api.post<ProjectResponse>(
      "/projects",
      toBody(input)
    );
    return toProject(data);
  } catch (error) {
    return rethrowValidation(error);
  }
}

async function updateProject(
  id: string,
  input: ProjectInput
): Promise<Project> {
  try {
    const { data } = await api.put<ProjectResponse>(
      `/projects/${id}`,
      toBody(input)
    );
    return toProject(data);
  } catch (error) {
    return rethrowValidation(error);
  }
}

async function setFavorite(id: string, favorite: boolean): Promise<Project> {
  const { data } = await api.patch<ProjectResponse>(
    `/projects/${id}/favorite`,
    { favorite }
  );
  return toProject(data);
}

async function uploadCover(file: File): Promise<string> {
  const body = new FormData();
  body.append("file", file);
  const { data } = await api.post<{ url: string }>("/uploads", body);
  return data.url;
}

async function deleteProject(id: string): Promise<void> {
  await api.delete(`/projects/${id}`);
}

export {
  DEFAULT_SORT,
  SORT_OPTIONS,
  ApiValidationError,
  createProject,
  deleteProject,
  getProject,
  listProjects,
  setFavorite,
  updateProject,
  uploadCover
};
export type {
  ListProjectsParams,
  Project,
  ProjectInput,
  ProjectList,
  ProjectSort
};

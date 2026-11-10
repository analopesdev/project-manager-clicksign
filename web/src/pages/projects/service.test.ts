import { beforeEach, describe, expect, it, vi } from "vitest";

const get = vi.fn();
const post = vi.fn();
const put = vi.fn();

vi.mock("axios", () => ({
  default: {
    create: () => ({ get, post, put, patch: vi.fn(), delete: vi.fn() })
  },
  isAxiosError: (e: unknown) =>
    (e as { isAxiosError?: boolean })?.isAxiosError === true
}));

const apiProject = {
  id: "abc",
  name: "Redesign do app",
  client: "Clicksign",
  start_date: "2026-01-15",
  end_date: "2026-06-30",
  cover_url: "https://example.com/c.png",
  favorite: true
};

async function service() {
  return import("./service");
}

describe("listProjects", () => {
  beforeEach(() => vi.clearAllMocks());

  it("converte snake_case da API para camelCase", async () => {
    get.mockResolvedValue({ data: { data: [apiProject], total: 7 } });
    const { listProjects } = await service();

    const result = await listProjects();

    expect(result.total).toBe(7);
    expect(result.projects).toEqual([
      {
        id: "abc",
        name: "Redesign do app",
        client: "Clicksign",
        startDate: "2026-01-15",
        endDate: "2026-06-30",
        coverUrl: "https://example.com/c.png",
        favorite: true
      }
    ]);
  });

  it("trata cover_url vazio como projeto sem capa", async () => {
    get.mockResolvedValue({
      data: { data: [{ ...apiProject, cover_url: "" }], total: 1 }
    });
    const { listProjects } = await service();

    expect((await listProjects()).projects[0]?.coverUrl).toBeUndefined();
  });

  it("manda os filtros como parâmetros e omite os vazios", async () => {
    get.mockResolvedValue({ data: { data: [], total: 0 } });
    const { listProjects } = await service();

    await listProjects({ search: "", favorite: false, sort: "end_date" });

    expect(get).toHaveBeenCalledWith("/projects", {
      params: { search: undefined, favorite: undefined, sort: "end_date" }
    });
  });

  it("envia favorite=true e a busca quando informados", async () => {
    get.mockResolvedValue({ data: { data: [], total: 0 } });
    const { listProjects } = await service();

    await listProjects({ search: "app", favorite: true });

    expect(get.mock.calls[0]?.[1]).toMatchObject({
      params: { search: "app", favorite: true }
    });
  });
});

describe("createProject", () => {
  beforeEach(() => vi.clearAllMocks());

  it("envia o corpo no formato da API, sem espaços sobrando", async () => {
    post.mockResolvedValue({ data: apiProject });
    const { createProject } = await service();

    await createProject({
      name: "  Redesign do app ",
      client: " Clicksign",
      startDate: "2026-01-15",
      endDate: "2026-06-30",
      coverUrl: ""
    });

    expect(post).toHaveBeenCalledWith("/projects", {
      name: "Redesign do app",
      client: "Clicksign",
      start_date: "2026-01-15",
      end_date: "2026-06-30",
      cover_url: ""
    });
  });

  it("transforma o 422 com `fields` em ApiValidationError com nomes do front", async () => {
    post.mockRejectedValue({
      isAxiosError: true,
      response: {
        status: 422,
        data: { fields: { end_date: "end_date must not be before start_date" } }
      }
    });
    const { ApiValidationError, createProject } = await service();

    const error = await createProject({
      name: "A B",
      client: "C",
      startDate: "2026-06-30",
      endDate: "2026-01-15",
      coverUrl: ""
    }).catch(e => e);

    expect(error).toBeInstanceOf(ApiValidationError);
    expect(error.fields).toEqual({
      endDate: "end_date must not be before start_date"
    });
  });

  it("repassa o 422 sem `fields` como erro comum (vira mensagem geral na tela)", async () => {
    const raw = {
      isAxiosError: true,
      response: { status: 422, data: { error: "texto solto" } }
    };
    post.mockRejectedValue(raw);
    const { ApiValidationError, createProject } = await service();

    const error = await createProject({
      name: "A B",
      client: "C",
      startDate: "2026-01-01",
      endDate: "2026-01-02",
      coverUrl: ""
    }).catch(e => e);

    expect(error).toBe(raw);
    expect(error).not.toBeInstanceOf(ApiValidationError);
  });
});

import { flushPromises, mount } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { defineComponent, reactive } from "vue";
import type { Project } from "../service";

const route = reactive<{ query: Record<string, string> }>({ query: {} });
const replace = vi.fn();
vi.mock("vue-router", () => ({
  useRoute: () => route,
  useRouter: () => ({ replace })
}));

const listProjects = vi.fn();
const setFavorite = vi.fn();
vi.mock("../service", () => ({
  SORT_OPTIONS: ["name", "start_date", "end_date"],
  DEFAULT_SORT: "name",
  deleteProject: vi.fn(),
  listProjects: (...a: unknown[]) => listProjects(...a),
  setFavorite: (...a: unknown[]) => setFavorite(...a)
}));

vi.mock("axios", () => ({
  default: {
    isCancel: (e: unknown) => (e as { canceled?: boolean })?.canceled === true
  }
}));

const project = (overrides: Partial<Project> = {}): Project => ({
  id: "1",
  name: "Projeto A",
  client: "C",
  startDate: "2026-01-01",
  endDate: "2026-02-01",
  favorite: false,
  ...overrides
});

function deferred(signal?: AbortSignal) {
  let resolve!: (value: unknown) => void;
  const promise = new Promise((res, rej) => {
    resolve = res;
    signal?.addEventListener("abort", () => rej({ canceled: true }));
  });
  return { promise, resolve };
}

const mounted: ReturnType<typeof mount>[] = [];

async function setup() {
  vi.resetModules();
  const { useProjectList } = await import("./useProjectList");
  let list!: ReturnType<typeof useProjectList>;
  mounted.push(
    mount(
      defineComponent({
        setup() {
          list = useProjectList();
          return () => null;
        }
      })
    )
  );
  await flushPromises();
  return list;
}

describe("useProjectList", () => {
  afterEach(() => {
    mounted.splice(0).forEach(wrapper => wrapper.unmount());
  });

  beforeEach(() => {
    vi.clearAllMocks();
    route.query = {};
    listProjects.mockResolvedValue({ projects: [project()], total: 1 });
  });

  describe("filtros na URL", () => {
    it("lê favoritos e ordenação da URL ao abrir a página", async () => {
      route.query = { favorite: "1", sort: "end_date" };
      const list = await setup();

      expect(list.favoriteOnly.value).toBe(true);
      expect(list.sort.value).toBe("end_date");
      expect(listProjects.mock.calls[0]?.[0]).toMatchObject({
        favorite: true,
        sort: "end_date"
      });
    });

    it("volta para ordem alfabética quando o sort da URL é inválido", async () => {
      route.query = { sort: "banana" };
      expect((await setup()).sort.value).toBe("name");
    });

    it("reabre a busca com o termo vindo da URL", async () => {
      route.query = { q: "  app " };
      const list = await setup();

      expect(list.searchOpen.value).toBe(true);
      expect(list.searchTerm.value).toBe("app");
      expect(listProjects.mock.calls[0]?.[0]).toMatchObject({ search: "app" });
    });

    it("escreve os filtros na URL ao mudar, sem os valores padrão", async () => {
      const list = await setup();

      list.favoriteOnly.value = true;
      await flushPromises();
      expect(replace).toHaveBeenLastCalledWith({ query: { favorite: "1" } });

      list.sort.value = "start_date";
      await flushPromises();
      expect(replace).toHaveBeenLastCalledWith({
        query: { favorite: "1", sort: "start_date" }
      });

      list.favoriteOnly.value = false;
      list.sort.value = "name";
      await flushPromises();
      expect(replace).toHaveBeenLastCalledWith({ query: {} });
    });
  });

  describe("busca", () => {
    it("pesquisa em todos os projetos, mesmo com 'Apenas Favoritos' ligado", async () => {
      route.query = { favorite: "1", q: "iridium" };
      await setup();

      expect(listProjects.mock.calls[0]?.[0]).toMatchObject({
        search: "iridium",
        favorite: false
      });
    });

    it("volta a filtrar favoritos quando a busca é limpa", async () => {
      route.query = { favorite: "1" };
      const list = await setup();
      expect(listProjects.mock.calls[0]?.[0]).toMatchObject({ favorite: true });

      list.searchTerm.value = "app";
      await flushPromises();
      expect(listProjects.mock.lastCall?.[0]).toMatchObject({
        favorite: false
      });

      list.searchTerm.value = "";
      await flushPromises();
      expect(listProjects.mock.lastCall?.[0]).toMatchObject({ favorite: true });
    });
  });

  describe("URL -> tela", () => {
    it("atualiza os filtros quando a URL muda (voltar, avançar, link)", async () => {
      const list = await setup();

      route.query = { favorite: "1", sort: "start_date" };
      await flushPromises();

      expect(list.favoriteOnly.value).toBe(true);
      expect(list.sort.value).toBe("start_date");
      expect(listProjects.mock.lastCall?.[0]).toMatchObject({
        favorite: true,
        sort: "start_date"
      });
    });

    it("abre a busca com o termo da URL e fecha quando ele some", async () => {
      const list = await setup();

      route.query = { q: "app" };
      await flushPromises();
      expect(list.searchOpen.value).toBe(true);
      expect(list.searchTerm.value).toBe("app");

      route.query = {};
      await flushPromises();
      expect(list.searchOpen.value).toBe(false);
      expect(list.searchTerm.value).toBe("");
    });
  });

  describe("lupa de busca", () => {
    it("só fica disponível quando existe algum projeto", async () => {
      listProjects.mockResolvedValue({ projects: [], total: 0 });
      await setup();
      const search = (
        await import("@/composables/useProjectSearch")
      ).useProjectSearch();
      expect(search.hasProjects.value).toBe(false);
    });

    it("continua disponível com favoritos filtrando e nenhum resultado", async () => {
      listProjects.mockResolvedValue({ projects: [], total: 3 });
      route.query = { favorite: "1" };
      await setup();
      const search = (
        await import("@/composables/useProjectSearch")
      ).useProjectSearch();
      expect(search.hasProjects.value).toBe(true);
    });

    it("fecha a busca quando o último projeto some", async () => {
      const list = await setup();
      list.searchTerm.value = "app";
      list.searchOpen.value = true;
      listProjects.mockResolvedValue({ projects: [], total: 0 });
      await list.load();
      expect(list.searchOpen.value).toBe(false);
      expect(list.searchTerm.value).toBe("");
    });
  });

  describe("requisições fora de ordem", () => {
    it("cancela a requisição anterior e só usa a resposta da mais recente", async () => {
      const calls: ReturnType<typeof deferred>[] = [];
      listProjects.mockImplementation((_params, signal) => {
        const call = deferred(signal);
        calls.push(call);
        return call.promise;
      });

      const list = await setup();
      list.sort.value = "end_date";
      await flushPromises();

      calls[1]!.resolve({
        projects: [project({ name: "Resposta nova" })],
        total: 1
      });
      await flushPromises();

      expect(list.projects.value.map(p => p.name)).toEqual(["Resposta nova"]);
      expect(list.error.value).toBeNull();
    });

    it("não apaga o 'carregando' da requisição nova quando a antiga é cancelada", async () => {
      const calls: ReturnType<typeof deferred>[] = [];
      listProjects.mockImplementation((_params, signal) => {
        const call = deferred(signal);
        calls.push(call);
        return call.promise;
      });

      const list = await setup();
      list.sort.value = "end_date";
      await flushPromises();

      expect(list.loading.value).toBe(true);

      calls[1]!.resolve({ projects: [], total: 0 });
      await flushPromises();
      expect(list.loading.value).toBe(false);
    });

    it("mostra erro quando a API falha", async () => {
      listProjects.mockRejectedValue(new Error("rede"));
      const list = await setup();

      expect(list.error.value).toMatch(/Não foi possível carregar/);
      expect(list.loading.value).toBe(false);
    });
  });

  describe("favoritar", () => {
    it("atualiza a estrela na hora e confirma na API", async () => {
      setFavorite.mockResolvedValue(project({ favorite: true }));
      const list = await setup();
      const target = list.projects.value[0]!;

      const pending = list.toggleFavorite(target);
      expect(target.favorite).toBe(true);
      await pending;

      expect(setFavorite).toHaveBeenCalledWith("1", true);
      expect(target.favorite).toBe(true);
    });

    it("desfaz a mudança na tela quando a API falha", async () => {
      setFavorite.mockRejectedValue(new Error("500"));
      const list = await setup();
      const target = list.projects.value[0]!;

      await list.toggleFavorite(target);

      expect(target.favorite).toBe(false);
    });

    it("com 'Apenas Favoritos' ligado, recarrega para tirar o desfavoritado", async () => {
      route.query = { favorite: "1" };
      listProjects.mockResolvedValue({
        projects: [project({ favorite: true })],
        total: 1
      });
      setFavorite.mockResolvedValue(project());
      const list = await setup();
      listProjects.mockClear();

      await list.toggleFavorite(list.projects.value[0]!);

      expect(listProjects).toHaveBeenCalledTimes(1);
    });
  });
});

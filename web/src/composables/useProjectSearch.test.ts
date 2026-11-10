import { beforeEach, describe, expect, it, vi } from "vitest";

const KEY = "project-search-history";

async function load() {
  vi.resetModules();
  return (await import("./useProjectSearch")).useProjectSearch();
}

describe("useProjectSearch", () => {
  beforeEach(() => localStorage.clear());

  it("guarda a busca mais recente primeiro e sem repetir", async () => {
    const search = await load();
    search.remember("alfa");
    search.remember("beta");
    search.remember("alfa");
    expect(search.history.value).toEqual(["alfa", "beta"]);
  });

  it("guarda no máximo 5 buscas", async () => {
    const search = await load();
    ["a", "b", "c", "d", "e", "f"].forEach(search.remember);
    expect(search.history.value).toEqual(["f", "e", "d", "c", "b"]);
  });

  it("ignora busca vazia ou só com espaços", async () => {
    const search = await load();
    search.remember("   ");
    expect(search.history.value).toEqual([]);
  });

  it("remove uma busca do histórico", async () => {
    const search = await load();
    search.remember("alfa");
    search.remember("beta");
    search.forget("alfa");
    expect(search.history.value).toEqual(["beta"]);
  });

  it("persiste o histórico entre recargas", async () => {
    (await load()).remember("alfa");
    expect((await load()).history.value).toEqual(["alfa"]);
  });

  it("aguenta um histórico corrompido no localStorage", async () => {
    localStorage.setItem(KEY, "{isso não é json");
    expect((await load()).history.value).toEqual([]);
  });

  it("guarda no histórico a busca que estava aberta ao fechar", async () => {
    const search = await load();
    search.isOpen.value = true;
    search.term.value = "alfa";
    search.close();
    expect(search.history.value).toEqual(["alfa"]);
  });

  it("não guarda nada ao fechar sem ter digitado", async () => {
    const search = await load();
    search.isOpen.value = true;
    search.close();
    expect(search.history.value).toEqual([]);
  });

  it("fechar a busca limpa o termo", async () => {
    const search = await load();
    search.isOpen.value = true;
    search.term.value = "alfa";
    search.close();
    expect(search.isOpen.value).toBe(false);
    expect(search.term.value).toBe("");
  });
});

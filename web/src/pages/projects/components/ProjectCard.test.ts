import { mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";
import { quasarPlugin } from "@/test/quasar";
import type { Project } from "../service";
import ProjectCard from "./ProjectCard.vue";

const project: Project = {
  id: "1",
  name: "Redesign do app",
  client: "Clicksign",
  startDate: "2026-01-15",
  endDate: "2026-06-30",
  favorite: false
};

const render = (overrides: Partial<Project> = {}, searchTerm?: string) =>
  mount(ProjectCard, {
    props: {
      project: { ...project, ...overrides },
      ...(searchTerm ? { searchTerm } : {})
    },
    global: { plugins: [quasarPlugin] }
  });

describe("ProjectCard", () => {
  it("mostra nome, cliente e as datas por extenso", () => {
    const text = render().text();
    expect(text).toContain("Redesign do app");
    expect(text).toContain("Clicksign");
    expect(text).toContain("15 de janeiro de 2026");
    expect(text).toContain("30 de junho de 2026");
  });

  it("não mostra o dia anterior por causa do fuso horário", () => {
    expect(render({ startDate: "2026-01-01" }).text()).toContain(
      "1 de janeiro de 2026"
    );
  });

  it("mostra o texto original quando a data não é válida", () => {
    expect(render({ startDate: "lixo" }).text()).toContain("lixo");
  });

  it("emite toggleFavorite com o projeto ao clicar na estrela", async () => {
    const wrapper = render();
    await wrapper.get("button.favorite").trigger("click");
    expect(wrapper.emitted("toggleFavorite")?.[0]?.[0]).toMatchObject({
      id: "1"
    });
  });

  it("descreve a estrela conforme o estado, para leitores de tela", () => {
    expect(render().get("button.favorite").attributes("aria-label")).toBe(
      "Adicionar Redesign do app aos favoritos"
    );
    const favorited = render({ favorite: true }).get("button.favorite");
    expect(favorited.attributes("aria-label")).toBe(
      "Remover Redesign do app dos favoritos"
    );
    expect(favorited.attributes("aria-pressed")).toBe("true");
  });

  it("destaca no título o termo buscado", () => {
    const wrapper = render({}, "redesign");
    expect(wrapper.get("mark").text()).toBe("Redesign");
  });
});

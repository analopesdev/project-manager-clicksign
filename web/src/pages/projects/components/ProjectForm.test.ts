import { flushPromises, mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";
import { quasarPlugin } from "@/test/quasar";
import type { ProjectFormValues } from "../schema";
import ProjectForm from "./ProjectForm.vue";

const valid: ProjectFormValues = {
  name: "Redesign do app",
  client: "Clicksign",
  startDate: "15/01/2026",
  endDate: "30/06/2026",
  coverUrl: ""
};

const render = (props: Record<string, unknown> = {}) =>
  mount(ProjectForm, { props, global: { plugins: [quasarPlugin] } });

const saveButton = (wrapper: ReturnType<typeof render>) =>
  wrapper.get("button[type=submit]");

describe("ProjectForm", () => {
  it("começa com o botão de salvar desabilitado", () => {
    expect(saveButton(render()).attributes("aria-disabled")).toBe("true");
  });

  it("não envia o formulário enquanto houver campo inválido", async () => {
    const wrapper = render({ initial: { ...valid, name: "Iridium" } });
    await wrapper.get("form").trigger("submit");
    expect(wrapper.emitted("submit")).toBeUndefined();
  });

  it("mostra o erro do campo depois que ele perde o foco", async () => {
    const wrapper = render();
    const name = wrapper.get("#name");
    await name.setValue("Iridium");
    await name.trigger("blur");
    expect(wrapper.text()).toContain(
      "Por favor, digite ao menos duas palavras"
    );
  });

  it("liga a mensagem de erro ao campo para leitores de tela", async () => {
    const wrapper = render();
    const name = wrapper.get("#name");
    await name.setValue("Iridium");
    await name.trigger("blur");

    const errorId = wrapper.get("#name-error").attributes("id");
    expect(name.attributes("aria-describedby")).toBe(errorId);
    expect(name.attributes("aria-invalid")).toBe("true");
  });

  it("habilita o botão e envia os valores quando tudo é válido", async () => {
    const wrapper = render({ initial: valid });
    expect(saveButton(wrapper).attributes("aria-disabled")).toBe("false");

    await wrapper.get("form").trigger("submit");

    expect(wrapper.emitted("submit")?.[0]).toEqual([valid, null]);
  });

  it("ignora um segundo envio enquanto o primeiro está em andamento", async () => {
    const wrapper = render({ initial: valid, submitting: true });
    expect(saveButton(wrapper).attributes("aria-disabled")).toBe("true");

    await wrapper.get("form").trigger("submit");
    await flushPromises();

    expect(wrapper.emitted("submit")).toBeUndefined();
  });

  it("mostra o erro devolvido pela API no campo certo", () => {
    const wrapper = render({
      initial: valid,
      serverErrors: { name: "name is required" }
    });
    expect(wrapper.text()).toContain("name is required");
  });
});

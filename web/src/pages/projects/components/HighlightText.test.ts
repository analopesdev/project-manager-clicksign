import { mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";
import HighlightText from "./HighlightText.vue";

function render(text: string, term?: string) {
  const host = document.createElement("div");
  mount(HighlightText, {
    props: { text, ...(term ? { term } : {}) },
    attachTo: host
  });
  return {
    marks: [...host.querySelectorAll("mark")].map(mark => mark.textContent),
    text: host.textContent
  };
}

describe("HighlightText", () => {
  it("destaca todas as ocorrências, sem diferenciar maiúsculas, e mantém o texto inteiro", () => {
    const result = render("Casa da casa", "casa");
    expect(result.marks).toEqual(["Casa", "casa"]);
    expect(result.text).toBe("Casa da casa");
  });

  it("não destaca nada sem termo ou com termo só de espaços", () => {
    expect(render("Redesign").marks).toEqual([]);
    expect(render("Redesign", "   ").marks).toEqual([]);
  });

  it("trata caracteres especiais do termo como texto, não como regex", () => {
    expect(render("Projeto (v2)", "(v2)").marks).toEqual(["(v2)"]);
    expect(render("a.b", ".").marks).toEqual(["."]);
  });
});

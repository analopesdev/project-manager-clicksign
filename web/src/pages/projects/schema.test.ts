import { describe, expect, it } from "vitest";
import {
  brToIso,
  coverFileSchema,
  isoToBr,
  projectFormSchema,
  projectToFormValues,
  type ProjectFormValues
} from "./schema";

const valid: ProjectFormValues = {
  name: "Redesign do app",
  client: "Clicksign",
  startDate: "15/01/2026",
  endDate: "30/06/2026",
  coverUrl: ""
};

function errorOf(values: ProjectFormValues, field: keyof ProjectFormValues) {
  const result = projectFormSchema.safeParse(values);
  if (result.success) return undefined;
  return result.error.issues.find(issue => issue.path[0] === field)?.message;
}

describe("projectFormSchema", () => {
  it("aceita um projeto completo", () => {
    expect(projectFormSchema.safeParse(valid).success).toBe(true);
  });

  it("exige ao menos duas palavras no nome", () => {
    expect(errorOf({ ...valid, name: "Iridium" }, "name")).toBe(
      "Por favor, digite ao menos duas palavras"
    );
    expect(errorOf({ ...valid, name: "   " }, "name")).toBeDefined();
    expect(
      errorOf({ ...valid, name: "Projeto Iridium" }, "name")
    ).toBeUndefined();
  });

  it("exige ao menos uma palavra no cliente", () => {
    expect(errorOf({ ...valid, client: "  " }, "client")).toBe(
      "Por favor, digite ao menos uma palavra"
    );
    expect(errorOf({ ...valid, client: "Iridium" }, "client")).toBeUndefined();
  });

  it("limita nome e cliente a 120 caracteres", () => {
    const long = `Projeto ${"a".repeat(120)}`;
    expect(errorOf({ ...valid, name: long }, "name")).toMatch(/no máximo 120/);
    expect(errorOf({ ...valid, client: "c".repeat(121) }, "client")).toMatch(
      /no máximo 120/
    );
  });

  it.each([
    ["vazia", ""],
    ["incompleta", "15/01/20"],
    ["inexistente (31/02)", "31/02/2026"],
    ["mês 13", "10/13/2026"]
  ])("recusa data %s", (_label, startDate) => {
    expect(errorOf({ ...valid, startDate }, "startDate")).toBe(
      "Selecione uma data válida"
    );
  });

  it("aceita 29/02 só em ano bissexto", () => {
    expect(
      errorOf(
        { ...valid, startDate: "29/02/2028", endDate: "01/03/2028" },
        "startDate"
      )
    ).toBeUndefined();
    expect(
      errorOf({ ...valid, startDate: "29/02/2027" }, "startDate")
    ).toBeDefined();
  });

  it("recusa data final anterior à inicial", () => {
    expect(
      errorOf(
        { ...valid, startDate: "30/06/2026", endDate: "15/01/2026" },
        "endDate"
      )
    ).toBe("A data final não pode ser anterior à data de início");
  });

  it("aceita data final igual à inicial", () => {
    expect(
      errorOf(
        { ...valid, startDate: "15/01/2026", endDate: "15/01/2026" },
        "endDate"
      )
    ).toBeUndefined();
  });
});

describe("coverFileSchema", () => {
  const file = (type: string, bytes = 10) =>
    new File([new Uint8Array(bytes)], "capa", { type });

  it("aceita .jpg e .png", () => {
    expect(coverFileSchema.safeParse(file("image/png")).success).toBe(true);
    expect(coverFileSchema.safeParse(file("image/jpeg")).success).toBe(true);
  });

  it("recusa outros tipos", () => {
    expect(coverFileSchema.safeParse(file("image/gif")).success).toBe(false);
    expect(coverFileSchema.safeParse(file("application/pdf")).success).toBe(
      false
    );
  });

  it("recusa arquivos acima de 5 MB", () => {
    expect(
      coverFileSchema.safeParse(file("image/png", 5 * 1024 * 1024 + 1)).success
    ).toBe(false);
  });
});

describe("conversão de datas", () => {
  it("converte AAAA-MM-DD (API) para dd/mm/aaaa (tela) e de volta", () => {
    expect(isoToBr("2026-01-15")).toBe("15/01/2026");
    expect(brToIso("15/01/2026")).toBe("2026-01-15");
  });

  it("devolve vazio quando a entrada está incompleta", () => {
    expect(isoToBr("")).toBe("");
    expect(brToIso("15/01")).toBe("");
  });

  it("monta os valores do formulário a partir de um projeto", () => {
    expect(
      projectToFormValues({
        id: "1",
        name: "A B",
        client: "C",
        startDate: "2026-02-01",
        endDate: "2026-12-31",
        favorite: false
      })
    ).toEqual({
      name: "A B",
      client: "C",
      startDate: "01/02/2026",
      endDate: "31/12/2026",
      coverUrl: ""
    });
  });
});

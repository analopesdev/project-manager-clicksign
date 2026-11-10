import { z } from "zod";
import type { Project } from "./service";

const MAX_TEXT_LENGTH = 120;
const MAX_COVER_MB = 5;
const MAX_COVER_BYTES = MAX_COVER_MB * 1024 * 1024;
const COVER_MIME_TYPES = ["image/jpeg", "image/png"];

const wordCount = (value: string) =>
  value.trim().split(/\s+/).filter(Boolean).length;

function parseBrDate(value: string): Date | null {
  const match = /^(\d{2})\/(\d{2})\/(\d{4})$/.exec(value);
  if (!match) return null;
  const [, day, month, year] = match.map(Number) as [
    number,
    number,
    number,
    number
  ];
  const date = new Date(year, month - 1, day, 12);
  const real =
    date.getFullYear() === year &&
    date.getMonth() === month - 1 &&
    date.getDate() === day;
  return real ? date : null;
}

const dateField = z.string().refine(value => parseBrDate(value) !== null, {
  message: "Selecione uma data válida"
});

const projectFormSchema = z
  .object({
    name: z
      .string()
      .max(MAX_TEXT_LENGTH, `Use no máximo ${MAX_TEXT_LENGTH} caracteres`)
      .refine(value => wordCount(value) >= 2, {
        message: "Por favor, digite ao menos duas palavras"
      }),
    client: z
      .string()
      .max(MAX_TEXT_LENGTH, `Use no máximo ${MAX_TEXT_LENGTH} caracteres`)
      .refine(value => wordCount(value) >= 1, {
        message: "Por favor, digite ao menos uma palavra"
      }),
    startDate: dateField,
    endDate: dateField,
    coverUrl: z.string()
  })
  .superRefine((values, ctx) => {
    const start = parseBrDate(values.startDate);
    const end = parseBrDate(values.endDate);
    if (start && end && end < start) {
      ctx.addIssue({
        code: "custom",
        path: ["endDate"],
        message: "A data final não pode ser anterior à data de início"
      });
    }
  });

const coverFileSchema = z
  .instanceof(File)
  .refine(file => COVER_MIME_TYPES.includes(file.type), {
    message: "Escolha uma imagem .jpg ou .png"
  })
  .refine(file => file.size <= MAX_COVER_BYTES, {
    message: `A imagem deve ter no máximo ${MAX_COVER_MB} MB`
  });

type ProjectFormValues = z.infer<typeof projectFormSchema>;
type FieldName = keyof ProjectFormValues;

const emptyValues = (): ProjectFormValues => ({
  name: "",
  client: "",
  startDate: "",
  endDate: "",
  coverUrl: ""
});

function isoToBr(iso: string): string {
  const [year, month, day] = iso.split("-");
  return year && month && day ? `${day}/${month}/${year}` : "";
}

function brToIso(br: string): string {
  const [day, month, year] = br.split("/");
  return day && month && year ? `${year}-${month}-${day}` : "";
}

function projectToFormValues(project: Project): ProjectFormValues {
  return {
    name: project.name,
    client: project.client,
    startDate: isoToBr(project.startDate),
    endDate: isoToBr(project.endDate),
    coverUrl: project.coverUrl ?? ""
  };
}

export {
  MAX_TEXT_LENGTH,
  brToIso,
  coverFileSchema,
  emptyValues,
  isoToBr,
  projectFormSchema,
  projectToFormValues
};
export type { FieldName, ProjectFormValues };

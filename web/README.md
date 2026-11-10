# Project Manager — web

Front-end em Vue 3 + Quasar (Vite) + TypeScript. O passo a passo para rodar tudo (API, banco e
web) está no [README da raiz](../README.md).

Atalhos, dentro de `web/` (Node 22.22 ou mais novo):

```bash
pnpm install
pnpm dev            # http://localhost:9000
pnpm test           # testes unitários (Vitest)
pnpm typecheck
pnpm lint:check     # formatação (oxfmt) + lint (oxlint); `pnpm lint` corrige
pnpm build          # gera dist/spa
```

## Estrutura (`src/`)

```
pages/
  index.vue                  lista de projetos (busca, favoritos, ordenação, remover)
  projects/                  tudo que é de projetos fica aqui dentro
    form.vue                 criar e editar (/projects/new e /projects/:id/edit)
    components/              card, filtros, busca, formulário, modal...
    composables/             useProjectList (lista + filtros na URL)
    service.ts               chamadas à API, tipo Project, snake_case -> camelCase
    schema.ts                validação do formulário com zod
components/Header.vue        header global (ícone de busca)
composables/useProjectSearch.ts   estado da busca e histórico (usado pelo header e pela lista)
css/                         tokens, mixins e tipografia (quasar.variables.scss)
```

Os testes ficam ao lado do que testam (`*.test.ts`). Só os `.vue` fora de `components/`
viram rotas.

# Project Manager — web

Front-end em **Vue 3 + Quasar (Vite) + TypeScript**. O passo a passo para rodar tudo (banco, API e
web) está no [README da raiz](../README.md); este arquivo cobre só o front.

## Scripts

Dentro de `web/`, com Node 22.22 ou mais novo:

| Comando           | O que faz                                                                                |
| ----------------- | ---------------------------------------------------------------------------------------- |
| `pnpm install`    | instala as dependências (o `postinstall` prepara os tipos do Quasar)                     |
| `pnpm dev`        | servidor de desenvolvimento em <http://localhost:9000>, com proxy de `/api` para `:8080` |
| `pnpm test`       | testes unitários e de componente (Vitest); `pnpm test:watch` para acompanhar             |
| `pnpm typecheck`  | TypeScript com `vue-tsc`                                                                 |
| `pnpm lint:check` | formatação (oxfmt), lint (oxlint) e estilos (Stylelint); `pnpm lint` corrige             |
| `pnpm build`      | build de produção em `dist/spa`                                                          |

## Variáveis de ambiente

Opcionais; veja `.env.example`.

| Variável       | Padrão | Observação                                                                                                                            |
| -------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------- |
| `VITE_API_URL` | `/api` | base das chamadas à API. Em desenvolvimento o proxy do Quasar encaminha `/api` para `http://localhost:8080`, então não precisa mudar. |

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
test/quasar.ts               helper que registra o Quasar nos testes de componente
```

Os testes ficam ao lado do que testam (`*.test.ts`). O roteamento é por arquivo: só os `.vue` fora
de `components/` viram rotas, e o `form.vue` declara as duas URLs com `definePage`.

## Como as partes se conversam

- **Lista (`useProjectList`)** guarda favoritos, ordenação e busca, mantém tudo sincronizado com a
  URL nos dois sentidos e carrega os projetos. Cada carga cancela a anterior, e o favoritar é
  otimista (desfaz se a API falhar). A remoção (modal de confirmação) também mora nele.
- **Busca (`useProjectSearch`)** é um estado compartilhado: a lupa do header abre o painel, que a
  página da lista exibe. O histórico das últimas 5 buscas fica no `localStorage`.
- **Formulário (`ProjectForm`)** valida com o schema do zod, mostra um erro por campo e só habilita
  "Salvar" com tudo válido. A capa escolhida é enviada para a API só ao salvar.
- **`service.ts`** é o único lugar que fala com a API e converte `snake_case` da API para
  `camelCase` da tela. Os componentes de apresentação não chamam a API.

## Estilos

- Cores, espaçamentos (`$space-1` a `$space-8`), raios e tamanhos de fonte são **tokens** em
  `src/css/quasar.variables.scss`, nomeados pelo papel (`$heading`, `$text-muted`, `$danger`...).
- **Mixins** concentram o que se repete: `@include text(tamanho, peso, linha)`,
  `button-large`, `button-primary`, `button-reset` e `sr-only`.
- O **Stylelint** ordena as propriedades por grupo, proíbe cor fixa fora do arquivo de variáveis e
  proíbe `!important`.
- Fontes: Encode Sans Expanded e Semi Expanded, carregadas do Google Fonts no `index.html`.
- Animações usam só `transform` e `opacity` e são desligadas com `prefers-reduced-motion`.

## Testes

Cobrem as regras que importam, com nomes que descrevem o comportamento:

- **Validação (zod):** nome com duas palavras, datas reais (inclusive 29/02 e 31/02), data final
  não anterior à inicial, tipo e tamanho da capa.
- **API (`service`):** conversão de campos, capa vazia, parâmetros omitidos e erros 422.
- **Lista (`useProjectList`):** filtros lidos e escritos na URL, busca ignorando favoritos,
  requisição antiga não vence a nova e favoritar com desfazer.
- **Busca (`useProjectSearch`):** histórico sem repetição, limite de 5 e dados corrompidos.
- **Componentes:** card (eventos e rótulos de acessibilidade), formulário (bloqueia envio
  inválido e duplo, erro ligado ao campo) e destaque do termo buscado.

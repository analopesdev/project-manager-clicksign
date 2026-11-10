# Project Manager

Aplicação para gerenciar projetos: lista com busca, favoritos e ordenação, e telas para criar,
editar e remover. Front em Vue 3 + Quasar, API em Go e Postgres, seguindo o design do Figma.

Detalhes de cada parte: [`api/README.md`](api/README.md) (API) e [`web/README.md`](web/README.md) (front).

## O que dá para fazer

- **Listar** os projetos em cards (capa, nome, cliente, datas) com estados de carregamento
  (esqueleto), erro, lista vazia e "nenhum resultado".
- **Buscar** pelo nome ou pelo cliente (ícone de lupa no header), com histórico das últimas
  5 buscas e o trecho encontrado destacado nos cards.
- **Favoritar** (a estrela muda na hora) e filtrar só os favoritos.
- **Ordenar** por ordem alfabética, projetos iniciados mais recentemente ou prazo mais próximo.
- **Criar e editar** com validação (nome com duas palavras ou mais, datas reais, data final não
  anterior à inicial) e envio de capa `.jpg`/`.png` com pré-visualização.
- **Remover** com confirmação.
- Filtros e busca ficam **na URL**: recarregar a página ou compartilhar o link mantém o estado.

## Tecnologias

| Camada | Tecnologia |
|---|---|
| Front | Vue 3, Quasar (Vite), TypeScript, zod, axios |
| API | Go, chi, uber/fx, pgx |
| Banco | Postgres 18 (Docker) |
| Testes | Vitest e Vue Test Utils (front), `go test` (API) |
| Qualidade | oxlint, oxfmt, Stylelint, vue-tsc |

## Como rodar

### Pré-requisitos

- **Docker** (para o Postgres)
- **Go** 1.26
- **Node** 22.22 ou mais novo (o Quasar exige) e **pnpm**
- Internet: as fontes (Encode Sans) vêm do Google Fonts

### Passo a passo

Em três terminais, a partir da raiz do repositório.

**1. Banco de dados**

```bash
cd api
docker compose up -d        # Postgres em localhost:5434
```

**2. API** (porta 8080)

```bash
cd api
cp .env.example .env        # já aponta para o Postgres do passo 1
go run ./cmd                # as migrações rodam na subida
```

Conferir: `curl localhost:8080/health` responde `{"status":"ok"}`. A documentação interativa da
API fica em <http://localhost:8080/docs>.

**3. Web** (porta 9000)

```bash
cd web
pnpm install
pnpm dev
```

Abra <http://localhost:9000>. Em desenvolvimento o Quasar faz proxy de `/api` para
`http://localhost:8080`, então não há CORS nem variável para configurar.

Com o banco recém-criado a lista começa vazia (aparece a tela "Nenhum projeto"): use o botão
**Novo projeto** para criar os primeiros.

### Se algo não funcionar

| Sintoma | Causa provável |
|---|---|
| A API fecha logo ao subir e o log termina em `DATABASE_URL is required` | falta o `.env` (rode `cp .env.example .env` dentro de `api/`) |
| O log termina em `connection refused` | o Postgres não está de pé (`docker compose up -d` em `api/`) |
| `Port 5434/8080/9000 already in use` | outro processo usa a porta (a 5434 é do Postgres, a 8080 da API, a 9000 do front) |
| `INCOMPATIBLE NODE VERSION` no `pnpm install` | use Node 22.22 ou mais novo |
| O front abre mas a lista dá erro | a API não está rodando na porta 8080 |
| A capa enviada não aparece | a API precisa estar acessível em `http://localhost:8080` (o front só permite essa origem para imagens em desenvolvimento) |

### Variáveis de ambiente

| Onde | Variável | Padrão | Para que serve |
|---|---|---|---|
| `api/.env` | `DATABASE_URL` | (obrigatória) | conexão com o Postgres; o exemplo já bate com o `docker-compose.yml` |
| `api/.env` | `PORT` | `8080` | porta da API |
| `api/.env` | `UPLOAD_DIR` | `uploads` | pasta onde ficam as capas enviadas (ignorada pelo git) |
| `web/.env` | `VITE_API_URL` | `/api` | endereço base da API; em desenvolvimento não precisa mudar |

## Telas e rotas

O front usa rotas com `#` (`createWebHashHistory`), então qualquer servidor estático serve.

| URL | Tela |
|---|---|
| `/#/` | lista de projetos (aceita `?q=`, `?favorite=1` e `?sort=name\|start_date\|end_date`) |
| `/#/projects/new` | novo projeto |
| `/#/projects/:id/edit` | editar projeto |

## Testes e verificações

```bash
cd api && go test ./...                          # upload: imagem válida, arquivo falso, 5 MiB, path traversal
cd web && pnpm test                              # validação, busca, filtros na URL, favoritar, formulário, API
cd web && pnpm typecheck                         # TypeScript (vue-tsc)
cd web && pnpm lint:check                        # formatação, lint e estilos (CSS)
cd web && pnpm build                             # build de produção em web/dist/spa
```

## Estrutura do repositório

```
api/                  API em Go (ver api/README.md)
  cmd/                ponto de entrada
  internal/           domínio, banco (com migrações), HTTP
  docs/openapi.yaml   contrato da API
web/                  front-end (ver web/README.md)
  src/pages/          páginas e tudo que é de projetos
  src/css/            tokens de design, mixins e tipografia
```

## API em resumo

Detalhes, decisões e exemplos em [`api/README.md`](api/README.md).

| Método | Caminho | Descrição |
|---|---|---|
| GET | `/projects?search=&favorite=&sort=` | lista; devolve `{ "data": [...], "total": N }` |
| POST | `/projects` | cria |
| GET | `/projects/{id}` | busca um projeto |
| PUT | `/projects/{id}` | atualiza |
| PATCH | `/projects/{id}/favorite` | marca ou desmarca como favorito |
| DELETE | `/projects/{id}` | remove (204) |
| POST | `/uploads` | envia a capa (multipart, campo `file`) e devolve a URL |
| GET | `/health` · `/docs` · `/openapi.yaml` | saúde, documentação interativa e contrato |

## Decisões

- **Filtros na URL.** Busca, "Apenas Favoritos" e ordenação ficam na query da rota
  (`/#/?q=app&favorite=1&sort=end_date`). Mudar a URL com a página aberta (voltar, avançar, colar
  um link) também atualiza a tela.
- **Busca em todos os projetos.** Enquanto há um termo, o filtro de favoritos é ignorado; a barra
  de filtros some nos resultados, então um filtro invisível não poderia esconder resultados.
- **Requisições fora de ordem.** Cada carga cancela a anterior (`AbortController`); só a mais
  recente atualiza a tela, e uma requisição cancelada não apaga o "carregando" da nova.
- **Carregamento sem piscar.** O esqueleto aparece só no primeiro carregamento; ao trocar filtro
  ou busca a lista antiga fica levemente apagada até a nova chegar.
- **Favoritar é otimista.** A estrela muda na hora e volta ao estado anterior se a API falhar.
- **Validação em dois níveis.** O front valida com **zod**, com as mensagens do design, e a API
  valida de novo. Cada erro aparece no campo certo, ligado ao input por `aria-describedby`.
- **Datas.** A API usa `AAAA-MM-DD`; a tela usa `dd/mm/aaaa`. As datas são montadas no fuso local
  (`new Date(ano, mês, dia)`) para não mostrar o dia anterior no Brasil.
- **Capa.** A API só guarda `cover_url`. O arquivo escolhido é enviado antes para `POST /uploads`
  (só `.jpg`/`.png`, até 5 MiB, tipo detectado pelo conteúdo, não pelo nome) e a URL devolvida
  vai em `cover_url`.
- **Total de projetos.** O `total` da lista ignora busca e filtros, então o contador do título
  não muda ao filtrar.
- **Design como sistema.** Cores, espaçamentos, raios e estilos de texto vêm de tokens e mixins
  (`web/src/css/quasar.variables.scss`); o Stylelint proíbe cor fixa fora desse arquivo e
  `!important`.
- **Acessibilidade.** Elementos semânticos, foco visível, navegação por teclado (Esc fecha a
  busca e os menus), rótulos nos botões só com ícone e animações desligadas para quem pede
  movimento reduzido (`prefers-reduced-motion`).

## Limitações conhecidas

- **A busca não ignora acentos** (`migracao` não acha `Migração`) e é por trecho do nome ou do
  cliente.
- **Sem paginação**: a lista traz todos os projetos de uma vez.
- **Sem autenticação**: qualquer pessoa com acesso à API pode criar, editar e remover.
- **As capas ficam no disco da API** (`api/uploads/`), então não sobrevivem a um ambiente
  descartável (container sem volume, por exemplo).
- **As fontes vêm do Google Fonts**: sem internet o texto cai na fonte padrão do sistema.

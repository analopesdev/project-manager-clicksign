# Project Manager

Cadastro de projetos: lista com busca, favoritos e ordenação, e telas de novo/editar/remover.

| Pasta | O que é |
|---|---|
| [`api/`](api/README.md) | API em Go (chi + Postgres + fx), com OpenAPI em `api/docs/openapi.yaml` |
| [`web/`](web/README.md) | Front-end em Vue 3 + Quasar + TypeScript |

## Pré-requisitos

- **Docker** (para o Postgres)
- **Go** 1.26
- **Node** 22.22 ou mais novo (o Quasar exige) e **pnpm**

## Rodando tudo

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

Conferir: `curl localhost:8080/health` responde `{"status":"ok"}`.

**3. Web** (porta 9000)

```bash
cd web
pnpm install
pnpm dev
```

Abra <http://localhost:9000>. Em desenvolvimento o Quasar faz proxy de `/api` para
`http://localhost:8080`, então não há CORS nem variável para configurar.

## Testes e verificações

```bash
cd api && go test ./...                         # upload (imagem válida, arquivo falso, 5 MiB, path traversal)
cd web && pnpm test                             # schema, busca, favoritar, formulário, service...
cd web && pnpm typecheck && pnpm lint:check     # tipos, formatação e lint
cd web && pnpm build                            # build de produção em web/dist/spa
```

## Decisões

- **Filtros na URL.** Busca, "Apenas Favoritos" e ordenação ficam na query da rota
  (`/#/?q=app&favorite=1&sort=end_date`): sobrevivem ao reload e podem ser compartilhados.
- **Requisições fora de ordem.** Cada carga cancela a anterior (`AbortController`); só a mais
  recente atualiza a tela, e uma requisição cancelada não apaga o "carregando" da nova.
- **Favoritar é otimista.** A estrela muda na hora e volta ao estado anterior se a API falhar.
- **Validação em dois níveis.** O front valida com **zod** (mensagens do design) e a API valida de
  novo; erros 422 da API são mostrados no campo certo.
- **Datas.** A API usa `AAAA-MM-DD`; a tela usa `dd/mm/aaaa`. As datas são montadas no fuso local
  (`new Date(ano, mês, dia)`) para não mostrar o dia anterior no Brasil.
- **Capa.** A API só guarda `cover_url`. O arquivo escolhido na tela é enviado antes para
  `POST /uploads` (só `.jpg`/`.png`, até 5 MiB, tipo detectado pelo conteúdo) e a URL devolvida
  vai em `cover_url`.
- **Total de projetos.** O `total` da lista ignora busca e filtros, então o contador do título
  não muda ao filtrar.

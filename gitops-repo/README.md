# GitOps - fiap-postech-fase3

Manifests Kubernetes gerenciados via ArgoCD para os microsserviços do
ToggleMaster. Cada serviço tem sua própria `Application` do ArgoCD e sua
própria estrutura Kustomize (`base` + `overlays/production`).

## Estrutura

```
gitops-repo/
├── <servico>-app.yaml              # Application do ArgoCD (aponta pro overlay abaixo)
└── <servico>/
    ├── base/                       # Manifests "genericos" (Deployment, Service, ConfigMap, Secret, ...)
    │   └── kustomization.yaml
    └── overlays/
        └── production/
            └── kustomization.yaml  # namespace + tag de imagem (unico ambiente hoje)
```

Os 5 serviços do ToggleMaster já existem aqui: `auth-service`, `flag-service`,
`targeting-service`, `evaluation-service` e `analytics-service`.

## Como usar

```bash
kubectl apply -f gitops-repo/auth-service-app.yaml
kubectl apply -f gitops-repo/flag-service-app.yaml
kubectl apply -f gitops-repo/targeting-service-app.yaml
kubectl apply -f gitops-repo/evaluation-service-app.yaml
kubectl apply -f gitops-repo/analytics-service-app.yaml
```
(ou `kubectl apply -f gitops-repo/*-app.yaml` pra todos de uma vez)

Cada `Application` já vem com `syncPolicy.automated` (prune + selfHeal) e
`CreateNamespace=true`, então o ArgoCD cria o namespace de cada serviço e
mantém o cluster sincronizado com este repositório sem intervenção manual.

**Antes do primeiro sync**: os `Secret` de cada serviço (`base/secret.yaml`)
têm placeholders de senha/chave — troque pelos valores reais (ver comentário
em cada arquivo) antes de aplicar, senão os pods sobem mas não conseguem
conectar no banco/serviços dependentes.

## Fluxo de deploy (CI -> GitOps -> ArgoCD)

Cada serviço tem seu próprio workflow em `.github/workflows/<servico>.yaml`,
todos seguindo o mesmo pipeline de 6 jobs (build/test, lint, security scan,
build+push da imagem, update do manifest, resumo). Exemplo com o
`analytics-service`:

1. Push em `main` dispara [`.github/workflows/analytics-service.yaml`](../.github/workflows/analytics-service.yaml).
2. `docker-build-and-push` builda, escaneia (Trivy) e publica a imagem no ECR
   com a tag `v1.0.0-<commit-hash>`.
3. `update-helm-manifest` atualiza o `newTag` em
   [`analytics-service/overlays/production/kustomization.yaml`](analytics-service/overlays/production/kustomization.yaml)
   com essa mesma tag e commita neste repositório (com `[skip ci]`, pra não
   disparar o workflow de novo).
4. O ArgoCD detecta a mudança no Git e sincroniza o Deployment automaticamente.

## Acesso externo

Cada serviço expõe um `Ingress` (`base/ingress.yaml`) sob um path próprio,
roteado por um unico controller `ingress-nginx` compartilhado
(`terraform/modules/ingress-nginx`, um NLB só para todos os serviços - não é
um LoadBalancer por serviço). Pra descobrir o endereço depois do
`terraform apply`:

```bash
terraform output -raw ingress_get_address | bash
```

| Serviço | Path | Porta |
|---|---|---|
| auth-service | `/auth` | 8001 |
| flag-service | `/flags` | 8002 |
| targeting-service | `/targeting` | 8003 |
| evaluation-service | `/evaluate` | 8004 |
| analytics-service | `/analytics` | 8005 |

Ex.: `http://<endereco-do-lb>/auth/health`.

## Convenções

- **Namespace**: um por serviço, mesmo nome do serviço (`auth-service`,
  `flag-service`, ...) - mesma organização do antigo `fiap-postech-fase2/eks/`.
- **Config**: toda variável de ambiente via `envFrom` (`configMapRef` +
  `secretRef`), nunca listada individualmente no container.
- **Imagem**: registro único no ECR (`<account>.dkr.ecr.us-east-1.amazonaws.com/fiap/<servico>`),
  tag imutável `v1.0.0-<commit-hash>` sobrescrita pelo transformer `images:`
  do overlay - o `base` nunca é aplicado sozinho em produção.
- **Autoscaling**: todos os 5 serviços usam `replicas` fixo no Deployment (2
  para as 4 APIs HTTP síncronas, 1 para o `analytics-service`, que é um
  worker de baixo volume). O `analytics-service` chegou a usar `ScaledObject`
  do KEDA (escalando pelo tamanho da fila SQS, como na Fase 2), mas foi
  desativado - o controller KEDA (`terraform/modules/keda`) continua
  instalado no cluster caso algum serviço volte a precisar dele, só não tem
  nenhum `ScaledObject` usando-o no momento.
- **Labels**: usamos o campo `labels` (com `includeSelectors: false`) em vez
  do `commonLabels` (deprecado) - evita que labels de metadata acabem
  vazando pro `selector` imutável do Deployment/Service.
- **Secrets**: `DATABASE_URL`/`MASTER_KEY`/`SERVICE_API_KEY` em cada
  `base/secret.yaml` são placeholders - trocar pelos valores reais (senha do
  RDS via `terraform output -json rds_database_urls`, chave gerada no
  auth-service) antes do primeiro `apply`.

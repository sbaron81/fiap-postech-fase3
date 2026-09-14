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

Hoje só existe `analytics-service`; os outros 4 (`auth-service`,
`flag-service`, `targeting-service`, `evaluation-service`) seguem o mesmo
padrão quando forem migrados.

## Como usar

```bash
kubectl apply -f gitops-repo/analytics-service-app.yaml
argocd app sync analytics-service   # ou deixa o automated sync fazer sozinho
```

A `Application` já vem com `syncPolicy.automated` (prune + selfHeal) e
`CreateNamespace=true`, então o ArgoCD cria o namespace `analytics-service` e
mantém o cluster sincronizado com este repositório sem intervenção manual.

## Fluxo de deploy (CI -> GitOps -> ArgoCD)

1. Push em `main` dispara [`.github/workflows/analytics-service.yaml`](../.github/workflows/analytics-service.yaml).
2. `docker-build-and-push` builda, escaneia (Trivy) e publica a imagem no ECR
   com a tag `v1.0.0-<commit-hash>`.
3. `update-helm-manifest` (a implementar) atualiza o `newTag` em
   [`analytics-service/overlays/production/kustomization.yaml`](analytics-service/overlays/production/kustomization.yaml)
   com essa mesma tag e commita neste repositório.
4. O ArgoCD detecta a mudança no Git e sincroniza o Deployment automaticamente.

## Convenções

- **Namespace**: um por serviço, mesmo nome do serviço (`analytics-service`,
  `auth-service`, ...) - mesma organização do antigo `fiap-postech-fase2/eks/`.
- **Config**: toda variável de ambiente via `envFrom` (`configMapRef` +
  `secretRef`), nunca listada individualmente no container.
- **Imagem**: registro único no ECR (`<account>.dkr.ecr.us-east-1.amazonaws.com/fiap/<servico>`),
  tag imutável `v1.0.0-<commit-hash>` sobrescrita pelo transformer `images:`
  do overlay - o `base` nunca é aplicado sozinho em produção.
- **Autoscaling**: serviços request/response (auth/flag/targeting) usam
  `replicas` fixo; serviços consumidores de fila (evaluation/analytics) usam
  `ScaledObject` do KEDA contra a fila SQS do Terraform, sem `replicas` fixo
  no Deployment (o HPA gerado pelo KEDA é quem manda).
- **Labels**: usamos o campo `labels` (com `includeSelectors: false`) em vez
  do `commonLabels` (deprecado) - evita que labels de metadata acabem
  vazando pro `selector` imutável do Deployment/Service.

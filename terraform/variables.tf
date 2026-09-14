variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "project_name" {
  description = "Prefixo usado no nome dos recursos (VPC, RDS, cache, etc.)"
  type        = string
  default     = "fiap-fase3"
}

variable "cluster_name" {
  description = "Nome do cluster EKS"
  type        = string
  default     = "fiap-fase3"
}

variable "kubernetes_version" {
  description = "null = usa a versao estavel mais recente disponivel na AWS"
  type        = string
  default     = null
}

variable "use_lab_role" {
  description = "true = reusa a LabRole do AWS Academy; false (padrao) = cria IAM roles proprias (conta pessoal)"
  type        = bool
  default     = false
}

variable "node_instance_types" {
  type    = list(string)
  default = ["t3.medium"]
}

variable "node_desired_size" {
  type    = number
  default = 2
}

variable "node_min_size" {
  type    = number
  default = 1
}

variable "node_max_size" {
  type    = number
  default = 4
}

variable "enable_nat_gateway" {
  description = "NAT Gateway e necessario para os nodes (em subnets privadas) puxarem imagens do ECR/internet. Desative apenas se souber o que esta fazendo."
  type        = bool
  default     = true
}

variable "sqs_queue_name" {
  type    = string
  default = "fiap-fase3-toggle-events"
}

variable "dynamodb_table_name" {
  type    = string
  default = "ToggleMasterAnalytics"
}

variable "elasticache_name" {
  type    = string
  default = "evaluation-service"
}

variable "argocd_namespace" {
  type    = string
  default = "argocd"
}

variable "argocd_chart_version" {
  description = "null = usa a versao mais recente do chart argo-cd"
  type        = string
  default     = null
}

variable "argocd_server_service_type" {
  description = "LoadBalancer expoe o argocd-server publicamente via ELB"
  type        = string
  default     = "LoadBalancer"
}

variable "keda_namespace" {
  type    = string
  default = "keda"
}

variable "keda_chart_version" {
  description = "null = usa a versao mais recente do chart keda"
  type        = string
  default     = null
}

variable "ingress_nginx_namespace" {
  type    = string
  default = "ingress-nginx"
}

variable "ingress_nginx_chart_version" {
  description = "null = usa a versao mais recente do chart ingress-nginx"
  type        = string
  default     = null
}

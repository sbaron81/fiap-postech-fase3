data "aws_availability_zones" "available" {
  state = "available"
}

locals {
  azs = slice(data.aws_availability_zones.available.names, 0, 2)
}

module "networking" {
  source = "./modules/networking"

  project_name       = var.project_name
  cluster_name       = var.cluster_name
  azs                = local.azs
  enable_nat_gateway = var.enable_nat_gateway
}

module "ecr" {
  source = "./modules/ecr"
}

module "eks" {
  source = "./modules/eks"

  cluster_name        = var.cluster_name
  kubernetes_version  = var.kubernetes_version
  use_lab_role        = var.use_lab_role
  cluster_subnet_ids  = concat(module.networking.public_subnet_ids, module.networking.private_subnet_ids)
  node_subnet_ids     = module.networking.private_subnet_ids
  node_instance_types = var.node_instance_types
  node_desired_size   = var.node_desired_size
  node_min_size       = var.node_min_size
  node_max_size       = var.node_max_size
}

module "rds" {
  source = "./modules/rds"

  project_name = var.project_name
  vpc_id       = module.networking.vpc_id
  vpc_cidr     = module.networking.vpc_cidr
  subnet_ids   = module.networking.private_subnet_ids
}

module "elasticache" {
  source = "./modules/elasticache"

  project_name = var.project_name
  vpc_id       = module.networking.vpc_id
  vpc_cidr     = module.networking.vpc_cidr
  subnet_ids   = module.networking.private_subnet_ids
  cache_name   = var.elasticache_name
}

module "dynamodb" {
  source = "./modules/dynamodb"

  table_name = var.dynamodb_table_name
}

module "sqs" {
  source = "./modules/sqs"

  queue_name = var.sqs_queue_name
}

module "argocd" {
  source = "./modules/argocd"

  namespace           = var.argocd_namespace
  chart_version       = var.argocd_chart_version
  server_service_type = var.argocd_server_service_type

  # O node group precisa existir (nodes Ready) antes de instalar o chart,
  # senao os pods do argocd ficam Pending indefinidamente.
  depends_on = [module.eks]
}

module "keda" {
  source = "./modules/keda"

  namespace     = var.keda_namespace
  chart_version = var.keda_chart_version

  depends_on = [module.eks]
}

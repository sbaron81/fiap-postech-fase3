

´´´bash
terraform init -backend-config=backend.hcl -reconfigure
terraform plan
terraform apply

terraform output -raw argocd_get_server_address | bash
terraform output -raw argocd_get_admin_password | bash

´´´



argocd_get_admin_password = "kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d; echo"
argocd_get_server_address = "kubectl -n argocd get svc argocd-server -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'; echo"
configure_kubectl = "aws eks update-kubeconfig --region us-east-1 --name fiap-fase3"
dynamodb_table_name = "ToggleMasterAnalytics"
ecr_repository_urls = {
  "fiap/analytics-service" = "343770067934.dkr.ecr.us-east-1.amazonaws.com/fiap/analytics-service"
  "fiap/auth-service" = "343770067934.dkr.ecr.us-east-1.amazonaws.com/fiap/auth-service"
  "fiap/evaluation-service" = "343770067934.dkr.ecr.us-east-1.amazonaws.com/fiap/evaluation-service"
  "fiap/flag-service" = "343770067934.dkr.ecr.us-east-1.amazonaws.com/fiap/flag-service"
  "fiap/targeting-service" = "343770067934.dkr.ecr.us-east-1.amazonaws.com/fiap/targeting-service"
}
eks_cluster_endpoint = "https://F9A2B3E410A2849061D5AE67536728A7.gr7.us-east-1.eks.amazonaws.com"
eks_cluster_name = "fiap-fase3"
elasticache_redis_url = "redis://evaluation-service-t0uqzp.serverless.use1.cache.amazonaws.com:6379"
private_subnet_ids = [
  "subnet-030a0394c7d5bb120",
  "subnet-031c3e181efff80b9",
]
public_subnet_ids = [
  "subnet-0afa0bdf3ab449673",
  "subnet-09fb0c38d9d80cdac",
]
rds_database_urls = <sensitive>
rds_endpoints = {
  "auth-service" = "fiap-fase3-auth-service.c3ivnpwh3dp3.us-east-1.rds.amazonaws.com:5432"
  "flag-service" = "fiap-fase3-flag-service.c3ivnpwh3dp3.us-east-1.rds.amazonaws.com:5432"
  "targeting-service" = "fiap-fase3-targeting-service.c3ivnpwh3dp3.us-east-1.rds.amazonaws.com:5432"
}
sqs_queue_url = "https://sqs.us-east-1.amazonaws.com/343770067934/fiap-fase3-toggle-events"
vpc_id = "vpc-07637d9fe05a864f8"

output "bucket_name" {
  description = "Nome do bucket S3 criado para teste de acesso"
  value       = aws_s3_bucket.test.bucket
}

output "bucket_arn" {
  description = "ARN do bucket S3 criado para teste de acesso"
  value       = aws_s3_bucket.test.arn
}

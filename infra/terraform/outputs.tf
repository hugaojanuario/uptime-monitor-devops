output "vpc_id" {
  value = module.vpc.vpc_id
}

output "public_subnets" {
  value = module.vpc.public_subnets
}

output "app_public_ip" {
  value = aws_instance.app.public_ip
}
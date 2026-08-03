module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = var.vpc_name
  cidr = var.vpc_cidr

  azs            = ["us-east-1a", "us-east-1b"]
  public_subnets = var.public_subnets

  enable_dns_hostnames = true

  tags = {
    Terraform = "true"
  }
}

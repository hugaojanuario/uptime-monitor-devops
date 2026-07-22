resource "aws_key_pair" "lab" {
  key_name = "uptime-lab"
  public_key = file("~/.ssh/uptime-lab.pub")
}

resource "aws_security_group" "app" {
  name        = "uptime-app-sg"
  description = "SSH e API do uptime monitor"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["187.34.95.155/32"]
  }

  ingress {
    description = "API"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Project = "go-uptime-aws"
    Env     = "lab"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"]
  }
}

resource "aws_instance" "app" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  subnet_id              = module.vpc.public_subnets[0]
  vpc_security_group_ids = [aws_security_group.app.id]
  key_name               = aws_key_pair.lab.key_name

  tags = {
    Name    = "uptime-app"
    Project = "go-uptime-aws"
    Env     = "lab"
  }
}
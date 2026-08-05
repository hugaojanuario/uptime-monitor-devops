data "aws_ami" "ubuntu_22_04" {
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_security_group" "uptime-monitor-security-group" {
  name = "allow_web"
  description = "Web trafic control"
  vpc_id = module.vpc.vpc_id

  ingress {
    description = "API"
    from_port = 8080
    to_port = 8080
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "SSH"
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Saida liberada para qualquer destino"
    from_port = 0
    to_port = 0
    protocol = "-1" #all protocols
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "uptime-server" {
  ami = data.aws_ami.ubuntu_22_04.id
  instance_type = "t3.micro"
  vpc_security_group_ids = [aws_security_group.uptime-monitor-security-group.id]

  subnet_id = module.vpc.public_subnets[0]
  associate_public_ip_address = true

  user_data = file("${path.module}/user_data.sh")
  user_data_replace_on_change = true

  key_name = aws_key_pair.lab.key_name
}

resource "aws_key_pair" "lab" {
  key_name   = "uptime-lab"
  public_key = file("~/.ssh/uptime-lab.pub")
}
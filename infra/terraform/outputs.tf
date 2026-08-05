output "ec2-public-ip" {
  value = aws_instance.uptime-server.public_ip
}
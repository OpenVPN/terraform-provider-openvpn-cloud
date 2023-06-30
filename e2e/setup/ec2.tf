provider "aws" {
  region = "eu-central-1"
  default_tags {
    tags = {
      task-group = "terraform-provider-openvpn-cloud"
      created-by = "Terraform/terraform-provider-openvpn-cloud"
    }
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

data "aws_vpc" "default" {
  default = true
}

data "template_file" "init" {
  template = file("${path.module}/user_data.sh.tpl")

  vars = {
    profile = local.connector_profile
  }
}

resource "aws_instance" "example" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  user_data = data.template_file.init.rendered
  key_name  = aws_key_pair.key.key_name

  security_groups = [aws_security_group.example.name]

  tags = {
    Name = "Terraform OpenVPN Provider Example for ${var.host_name}"
  }
}

resource "aws_security_group" "example" {
  name        = "${var.host_name}-sg"
  description = "Terraform Provider Example Security Group for ${var.host_name}"
  vpc_id      = data.aws_vpc.default.id

  // To Allow SSH Transport
  ingress {
    from_port   = 22
    protocol    = "tcp"
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  // To Allow Port 80 Transport
  ingress {
    from_port   = 80
    protocol    = "tcp"
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }
  tags = {}
}

resource "aws_key_pair" "key" {
  key_name   = "${var.host_name}-key"
  public_key = file("~/.ssh/id_rsa.pub")
  tags       = {}
}

output "instance_id" {
  value = aws_instance.example.id
}

output "instance_public_ip" {
  value = aws_instance.example.public_ip
}

output "instance_private_ip" {
  value = aws_instance.example.private_ip
}

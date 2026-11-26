resource "aws_security_group" "web" {
  name        = "synthetic-web"
  description = "Synthetic web security group"

  ingress {
    description = "https"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16", "10.1.0.0/16"]
  }

  ingress {
    description = "http"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  tags = { Name = "synthetic-web", Environment = "test" }
}

resource "aws_iam_policy" "read" {
  name = "synthetic-read"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["s3:GetObject"]
      Resource = "arn:aws:s3:::synthetic-data-bucket/*"
    }]
  })
}

resource "aws_instance" "web" {
  ami           = "ami-00000000"
  instance_type = "t3.micro"

  root_block_device {
    volume_size = 20
  }
}

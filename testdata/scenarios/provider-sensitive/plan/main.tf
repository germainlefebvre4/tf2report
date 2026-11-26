resource "aws_db_instance" "main" {
  identifier          = "synthetic-db"
  engine              = "postgres"
  instance_class      = "db.t3.micro"
  allocated_storage   = 20
  username            = "synthetic"
  password            = "SYNTHETIC-SECRET-PASSWORD"
  skip_final_snapshot = true
}

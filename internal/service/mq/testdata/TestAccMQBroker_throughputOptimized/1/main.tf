
resource "aws_security_group" "test" {
  name = var.random_name

  tags = {
    Name = var.random_name
  }
}

resource "aws_mq_broker" "test" {
  broker_name        = var.random_name
  engine_type        = var.engine_type
  engine_version     = var.engine_version
  storage_type       = var.storage_type
  host_instance_type = var.host_instance_type
  security_groups    = [aws_security_group.test.id]

  logs {
    general = var.general
  }

  user {
    username = var.username
    password = var.password
  }
}

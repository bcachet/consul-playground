services {
  name = "backend"
  id   = "backend"

  address =  "127.0.0.1"
  port = 5000

  checks = [
    {
      name = "healcheck"
      id = "healthcheck"
      interval = "5s"
      timeout = "2s"
      status = "critical"

      http = "http://127.0.0.1:5000/health"
    }
  ]
}

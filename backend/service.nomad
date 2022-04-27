job "oprox" {
  datacenters = ["dc1"]
  type = "service"

  group "default" {
    network {
      port "oprox_port" {
        static = 19303
      }
    }

    service {
      name = "oprox"
      port = "oprox_port"

      # The "check" stanza instructs Nomad to create a Consul health check for
      # this service. A sample check is provided here for your convenience;
      # uncomment it to enable it. The "check" stanza is documented in the
      # "service" stanza documentation.

      # check {
      #   name     = "alive"
      #   type     = "tcp"
      #   interval = "10s"
      #   timeout  = "2s"
      # }
    }

    task "oprox" {
      driver = "docker"

      resources {
        cpu    = 500
        memory = 512
        memory_max = 1024
      }

      config {
        image        = "ghcr.io/${IMAGE_NAME}:${TAG}"
        network_mode = "host"
        oom_score_adj = -1000
      }

      env {
        ADDR = "${NOMAD_ADDR_oprox_port}"
        
        ONEC_TOKEN = "${ONEC_TOKEN}"
        MINIO_ACCESS_KEY = "${MINIO_ACCESS_KEY}"
        MINIO_SECRET_KEY = "${MINIO_SECRET_KEY}"
      }
    }
  }
}



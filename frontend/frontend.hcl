job "re-star-ru" {
	datacenters = ["dc1"]
	type = "service"

	group "default" {

		network {
			port "restar" {
				host_network = "private"
			}
		}

		task "restar" {
			service {
				port = "restar"
				tags = [
					"reproxy.enabled=1",
//					"reproxy.server=feziv.com,www.feziv.com",
					"timestamp=[[timeNow]]"
				]
			}
			// serve static files for feziv.com
			resources {
				memory = 64
		  }

			driver = "docker"

			env {
				LISTEN = "${NOMAD_ADDR_restar}"
			}

			config {
				image = "ghcr.io/[[.repo]]:[[.tag]]"
				network_mode = "host"
			}
		}
	}
}



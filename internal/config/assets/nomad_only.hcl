env "nomad" {
  nomad {
    address = "localhost:4646"
    acl_token = ""
    connection {
      type = "direct"
    }
  }
}
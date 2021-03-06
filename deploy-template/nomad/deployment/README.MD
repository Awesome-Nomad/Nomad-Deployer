```yaml
# Name of deployment. Should be global unique
service_name: "be-notification"
# Docker image of application. Will be ignored on production deployment.
image: "registry.hub.docker.com/library/be-be-notification:bbd7e9b0a4512e8b49563279280f1e74eb47b700"
# Number of instances. Default: 2
count: 2
# Number of canary instances. Default to count
canary: 2
# Auto revert to previous version if error occur. Default: false
auto_revert: true
# Auto promote to new version. Default: false
auto_promote: true
# Maximum update instances at a time. Default: 1
max_parallel: 1
# Application
app_files: # Use to copy local file from repo to app. In application will reallocated at /local/<destination>
  - src: deployment/config.toml
    destination: config.toml # => /local/config.toml
    env: false # Is this file will use as environment variables. Default: false.
  - src: deployment/env
    destination: .env # => /local/.env
    env: true
# Define how to run your app
entrypoint:
  - /app
  - start
  - --config=${CONFIG_FILE}
directories: # List of volumes that app will use. Default: []
  - name: lixi_storage
    read_only: false
    destination: /data/lixi_storage # mapped into application path
  - name: lixi_counter
    read_only: false
    destination: /data/lixi_counter    # mapped into application path
# Register services into consul
services:
  - name: be-notification
    port: grpc
    prometheus: true # Dedicated that service will have prometheus endpoint.
    tags: # Optional.
      - grpc
 - name: be-notification-http
   port: http
   tags:
     - traefik.enable=true
     - traefik.http.routers.api-merchant.rule=PathPrefix(`/v1/merchant`)
     - traefik.http.routers.api-merchant.entrypoints=http
     - traefik.http.routers.api-merchant.middlewares=global-retry@file
     # Use below tag for gRPC only.
     - traefik.http.services.api-merchant.loadbalancer.server.scheme=h2c
# Resources
resources: # Required. At least one of below attributes must be defined
  mbits: "1" # Net speed. Default: 1
  cpu: "100" # MHz. Default: 100
  memory: "300" # MB. Default: 300

# Add host. Format: hostname:IP
extra_hosts:
  - "consul-dev.liemdeptrai:10.30.83.2"
# Define contraints for jobs. Without this, your job may run incorrectly.  
constraints:
  # Using internet. Eg: Google, Facebook, etc...    
  - key: "public_network"
    operator: "="
    value: "True"        
```

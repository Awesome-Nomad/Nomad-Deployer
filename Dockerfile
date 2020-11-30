ARG BASE_IMAGE=alpine:3.12
FROM golang:1.14 as builder
COPY . /work
WORKDIR /work
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -ldflags '-extldflags "-static"' -o /deployer

FROM ${BASE_IMAGE}
COPY --from=builder /deployer /deployer
ARG build_env=development
COPY .environment/$build_env.deployer.hcl /deployer.hcl
COPY deploy-template /deploy-template
ENV ENVIRONMENT="${build_env}"


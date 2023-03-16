load('ext://restart_process', 'docker_build_with_restart')

local_resource(
  'enats-bin',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -gcflags=all=-trimpath=${GOPATH} -asmflags=all=-trimpath=${GOPATH} -o build/enats ./cli && echo "Build Complete!"',
  deps=['./cli', './internal', './pkg', './go.mod', './go.sum'],
  labels="enats"
)

docker_build_with_restart(
  'enats',
  context='.',
  entrypoint='/app/enats',
  dockerfile='./Dockerfile.tilt',
  only=['./build/enats'],
  live_update=[
    sync('./build/enats', '/app/enats')
  ]
)

# Resources
k8s_yaml(kustomize('./kustomize'))

# ENATs Service
k8s_resource(
  'enats',
  labels="enats",
  port_forwards=['4222', '6222', '8222'],
)
target := "goddns"
image  := "quay.io/svenhertle/goddns"

#
# GO BUILD
#

# build all binaries
all: build-x64-linux build-arm64-macos

# remove build output
clean:
    rm -rf build

# build for linux/amd64
build-x64-linux: (_build "linux" "amd64" "x64-linux")

# build for macOS/arm64 (Apple Silicon)
build-arm64-macos: (_build "darwin" "arm64" "arm64-macos")

_build os arch name:
    mkdir -p build
    CGO_ENABLED=0 GOOS={{os}} GOARCH={{arch}} go build -o build/{{target}}-{{name}} .

#
# DOCKER BUILD
#

# docker if available, otherwise podman
container_engine := `command -v docker || command -v podman || true`

# main: "latest" + most recent git tag (e.g. v1.2.3); other branches: branch name with "/" -> "-"
branch := `git rev-parse --abbrev-ref HEAD 2>/dev/null || true`
tags := if branch == "main" {
    "latest " + `git describe --tags --abbrev=0 2>/dev/null || true`
} else {
    replace(branch, "/", "-")
}

# build the container image with the tags for the current branch
docker-build:
    #!/usr/bin/env bash
    set -euo pipefail
    args=()
    for t in {{tags}}; do args+=(-t "{{image}}:$t"); done
    {{container_engine}} build "${args[@]}" .

# build and push the container image (main only)
docker-push: _check-main docker-build
    #!/usr/bin/env bash
    set -euo pipefail
    for t in {{tags}}; do {{container_engine}} push "{{image}}:$t"; done

_check-main:
    @test "{{branch}}" = "main" || { echo "Refusing to push: current branch is '{{branch}}', only 'main' can be pushed" >&2; exit 1; }

FROM golang:1.27.4@sha256:5f1b98fa02fe396ebce2a0227a0e08309e11792915b9dd2e1f71053d7e7e6666

ARG GOLANGCI_LINT_VERSION=2.15.0 # renovate: depName=golangci/golangci-lint datasource=github-releases

RUN apt update && apt upgrade -y && \
    apt-get install --no-install-recommends -y ca-certificates curl gnupg lsb-release jq zsh neovim gh && \
    chsh -s $(which zsh) && \
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" && \
    rm -rf /var/lib/apt/lists/*

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_LINT_VERSION}


RUN  curl --proto '=https' --tlsv1.2 -fsSL https://get.opentofu.org/install-opentofu.sh -o install-opentofu.sh && \
     chmod +x install-opentofu.sh && \
     ./install-opentofu.sh --install-method deb && \
    rm -f install-opentofu.sh

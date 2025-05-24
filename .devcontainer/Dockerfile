ARG GO_VERSION=1.24.3
ARG GOLANGCI_LINT_VERSION=2.1.6 # renovate: depName=golangci/golangci-lint datasource=github-releases

FROM golang:${GO_VERSION}

# install docker
RUN apt update && apt upgrade -y && apt-get install -y ca-certificates curl gnupg lsb-release jq zsh neovim && \
    chsh -s $(which zsh) && \
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" && \
    rm -rf /var/lib/apt/lists/*

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}


RUN  curl --proto '=https' --tlsv1.2 -fsSL https://get.opentofu.org/install-opentofu.sh -o install-opentofu.sh && \
     chmod +x install-opentofu.sh && \
     ./install-opentofu.sh --install-method deb && \
    rm -f install-opentofu.sh

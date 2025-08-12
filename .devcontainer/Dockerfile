FROM ubuntu:22.04

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Create arheanja user with root privileges
RUN groupadd --gid 1000 arheanja \
    && useradd --uid 1000 --gid arheanja --shell /bin/bash --create-home arheanja \
    && usermod -aG root arheanja

# Install basic tools
RUN apt-get update && apt-get install -y \
    git \
    curl \
    wget \
    unzip \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release \
    sudo \
    vim \
    nano \
    jq \
    make \
    && rm -rf /var/lib/apt/lists/*

# Add arheanja to sudoers with full root privileges
RUN echo arheanja ALL=\(ALL:ALL\) NOPASSWD:ALL > /etc/sudoers.d/arheanja \
    && chmod 0440 /etc/sudoers.d/arheanja

# Install Docker Engine (not just CLI)
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update \
    && apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin \
    && rm -rf /var/lib/apt/lists/*

# Add arheanja user to the docker group
RUN usermod -aG docker arheanja

# Install Go
ENV GO_VERSION=1.21.5
RUN ARCH=$(dpkg --print-architecture) \
    && wget https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-${ARCH}.tar.gz \
    && rm go${GO_VERSION}.linux-${ARCH}.tar.gz

# Install kubectl
RUN ARCH=$(dpkg --print-architecture | sed 's/aarch64/arm64/') \
    && curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${ARCH}/kubectl" \
    && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl \
    && rm kubectl

# Install Helm
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install minikube
RUN ARCH=$(dpkg --print-architecture) \
    && curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-${ARCH} \
    && chmod +x minikube \
    && mv minikube /usr/local/bin/

# Install conntrack (required by minikube)
RUN apt-get update && apt-get install -y conntrack && rm -rf /var/lib/apt/lists/*

# Set up Go environment
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Create GOPATH directories
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chown -R arheanja:arheanja "$GOPATH"

# Create workspace directory
RUN mkdir -p /workspace && chown arheanja:arheanja /workspace

# Create minikube and kubectl directories with proper permissions
RUN mkdir -p /home/arheanja/.minikube /home/arheanja/.kube \
    && chown -R arheanja:arheanja /home/arheanja/.minikube /home/arheanja/.kube \
    && chmod -R 755 /home/arheanja/.minikube /home/arheanja/.kube

# Create Docker daemon startup script
RUN cat > /usr/local/bin/start-docker.sh << 'EOF'
#!/bin/bash
# Start Docker daemon
sudo dockerd --host=unix:///var/run/docker.sock --host=tcp://0.0.0.0:2376 &
sleep 5
echo "Docker daemon started"
EOF
RUN chmod +x /usr/local/bin/start-docker.sh

# Make entrypoint script executable
COPY scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Switch to arheanja user
USER arheanja

# Set the default working directory
WORKDIR /workspace

# Install common Go tools
# RUN go install -v golang.org/x/tools/gopls@latest \
#     && go install -v github.com/ramya-rao-a/go-outline@latest \
#     && go install -v github.com/cweill/gotests/...@latest \
#     && go install -v github.com/fatih/gomodifytags@latest \
#     && go install -v github.com/josharian/impl@latest \
#     && go install -v github.com/haya14busa/goplay/cmd/goplay@latest \
#     && go install -v github.com/go-delve/delve/cmd/dlv@latest \
#     && go install -v honnef.co/go/tools/cmd/staticcheck@latest

# Set default shell
SHELL ["/bin/bash", "-c"]
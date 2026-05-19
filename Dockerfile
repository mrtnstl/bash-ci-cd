FROM ubuntu:24.04

# prevent interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# install dependencies
RUN apt-get update && apt-get install -y \
    libsecret-tools \
    gnome-keyring \
    dbus-x11 \
    jq \
    curl \
    git \
    ca-certificates \
    gnupg \
    nano \
    && rm -rf /var/lib/apt/lists/*

# setting up node.js
ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION=24

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.4/install.sh | bash

# install + symlinks
RUN . "$NVM_DIR/nvm.sh" && \
    nvm install $NODE_VERSION && \
    nvm use $NODE_VERSION && \
    nvm alias default $NODE_VERSION && \
    ln -sf "$NVM_DIR/versions/node/v$NODE_VERSION."*/bin/node /usr/local/bin/node && \
    ln -sf "$NVM_DIR/versions/node/v$NODE_VERSION."*/bin/npm  /usr/local/bin/npm && \
    ln -sf "$NVM_DIR/versions/node/v$NODE_VERSION."*/bin/npx  /usr/local/bin/npx

# .bashrc setup
RUN echo 'export NVM_DIR="$HOME/.nvm"' >> /root/.bashrc && \
    echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"' >> /root/.bashrc && \
    echo '[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"' >> /root/.bashrc

WORKDIR /bash-ci-cd

COPY . .

CMD ["/bin/bash"]

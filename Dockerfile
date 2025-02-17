FROM node:22.0-bullseye-slim

# install dependencies
RUN apt-get update && \
    apt-get install -y \
    build-essential \
    git \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    nsis \
    wget \
    zsh \
    && rm -rf /var/lib/apt/lists/*

# setup zsh and oh-my-zsh
RUN git clone --single-branch --depth 1 https://github.com/robbyrussell/oh-my-zsh.git ~/.oh-my-zsh
RUN cp ~/.oh-my-zsh/templates/zshrc.zsh-template ~/.zshrc
RUN chsh -s /bin/zsh

# install go 1.23.5
ARG GO_VERSION=1.23.5
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# add go to path
ENV PATH=$PATH:/usr/local/go/bin
# set go path
ENV GOPATH=/usr/local/go

# install wails
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Update wails
RUN wails update

# install delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# make sure wails runs correctly
RUN wails doctor
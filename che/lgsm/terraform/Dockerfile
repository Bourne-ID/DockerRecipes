FROM eclipse/stack-base:ubuntu
EXPOSE 4403 8000 8080 9876 22

LABEL che:server:8080:ref=tomcat8 che:server:8080:protocol=http che:server:8000:ref=tomcat8-debug che:server:8000:protocol=http che:server:9876:ref=codeserver che:server:9876:protocol=http

RUN sudo dpkg --add-architecture i386 && \ 
    sudo apt-get update && \
    sudo apt-get -y install build-essential \
    gcc \ 
    make \
    git \
    jq \
    wget \
    unzip \
    tmux \
    sudo apt-get clean && \
    sudo apt-get -y autoremove && \
    sudo apt-get -y clean && \
    sudo rm -rf /var/lib/apt/lists/*
    
RUN sudo npm install -g markdown-table-cli

COPY https://releases.hashicorp.com/terraform/0.11.9/terraform_0.11.9_linux_amd64.zip /terraform

ENV PATH="PATH="/terraform:${PATH}"

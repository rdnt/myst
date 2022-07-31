FROM ubuntu:22.04 as builder

RUN apt-get update && \
    apt-get upgrade -y

FROM builder

RUN apt-get install -y --no-install-recommends \
        ca-certificates \
        sudo \
        curl

RUN apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

ARG USERNAME=tasos
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID --shell /bin/zsh -m $USERNAME \
    # Add sudo support
    && apt-get update \
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

ADD install.sh /home/$USERNAME/install.sh
RUN chmod +x /home/$USERNAME/install.sh

USER $USERNAME

RUN /home/$USERNAME/install.sh

WORKDIR /home/$USERNAME/app

ENTRYPOINT [ "zsh" ]

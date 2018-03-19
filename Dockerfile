FROM ubuntu:16.04


# env
ENV workdir /configbuddy_ws
ENV configbuddy_app ${workdir}/configbuddy
ENV dotfiles ${workdir}/dotfiles

# workdir
WORKDIR ${workdir}

# setup
RUN apt-get update
RUN apt-get install -y vim
RUN apt-get install -y git

COPY Dockerfile dotfiles* $workdir/
ADD configbuddy.v2 $workdir

ENTRYPOINT "/bin/bash"

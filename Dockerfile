FROM golang:1.12.1
#FROM yottta/arch-go:go1.12.1

# env
ENV workdir /configbuddy_ws
ENV configbuddy_app /go/src/github.com/yottta/configbuddy.v2
ENV dotfiles ${workdir}/dotfiles

# workdir
WORKDIR ${workdir}
ADD ./dotfiles $workdir/dotfiles
RUN echo "./configbuddy.v2 -c dotfiles/configs/work.yml -l debug" >> /root/.bash_history
RUN echo "alias ll='ls -lah'" >> /root/.bashrc

# setup
ADD . $configbuddy_app
# copy vendor from the project during image build
RUN cd $configbuddy_app && go build -o ${workdir}/configbuddy.v2 -v

ENTRYPOINT "/bin/bash"

FROM golang:1.10

# env
ENV workdir /configbuddy_ws
ENV configbuddy_app /go/src/github.com/andreic92/configbuddy.v2
ENV dotfiles ${workdir}/dotfiles

# workdir
WORKDIR ${workdir}
ADD ./dotfiles $workdir/dotfiles
RUN echo "./configbuddy.v2 -c dotfiles/configs/v2/work.yml -l debug" >> /root/.bash_history

# setup
ADD . $configbuddy_app
# copy vendor from the project during image build
#RUN cd $configbuddy_app && go get -u github.com/kardianos/govendor
#RUN cd $configbuddy_app && govendor sync -v
RUN cd $configbuddy_app && go build -o ${workdir}/configbuddy.v2 -v

ENTRYPOINT "/bin/bash"

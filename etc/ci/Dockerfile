FROM base/archlinux

ENV GOPATH /go

RUN pacman -Syu --noconfirm base-devel cmake git go python subversion
RUN go get -d llvm.org/llvm/bindings/go/llvm
RUN $GOPATH/src/llvm.org/llvm/bindings/go/build.sh -DCMAKE_BUILD_TYPE=Release -DLLVM_PARALLEL_LINK_JOBS=1

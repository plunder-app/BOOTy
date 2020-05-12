# syntax=docker/dockerfile:experimental

# Build LVM2 as an init
FROM gcc:latest as LVM
RUN wget https://mirrors.kernel.org/sourceware/lvm2/LVM2.2.03.09.tgz
RUN tar -xf LVM2.2.03.09.tgz
WORKDIR LVM2.2.03.09
RUN apt-get update; apt-get install -y libaio-dev libdevmapper-dev
RUN ./configure --enable-static_link --disable-selinux
# UGLY HaCk
RUN sed -i '/DMLIBS = -ldevmapper/ s/$/ -lm -lpthread/' libdm/dm-tools/Makefile
# GroSS HaCK
RUN make; exit 0
WORKDIR tools
# My EyEs .. They bl33d HaCk
RUN  gcc -O2  -fPIC  -static -L command.o dumpconfig.o formats.o lvchange.o lvconvert.o lvconvert_poll.o lvcreate.o lvdisplay.o lvextend.o lvmcmdline.o lvmdiskscan.o lvpoll.o lvreduce.o lvremove.o lvrename.o lvresize.o lvscan.o polldaemon.o pvchange.o pvck.o pvcreate.o pvdisplay.o pvmove.o pvmove_poll.o pvremove.o pvresize.o pvscan.o reporter.o segtypes.o tags.o toollib.o vgcfgbackup.o vgcfgrestore.o vgchange.o vgck.o vgcreate.o vgdisplay.o vgexport.o vgextend.o vgimport.o vgimportclone.o vgmerge.o vgmknodes.o vgreduce.o vgremove.o vgrename.o vgscan.o vgsplit.o lvm-static.o ../lib/liblvm-internal.a ../libdaemon/client/libdaemonclient.a ../device_mapper/libdevice-mapper.a ../base/libbase.a  -lm -lblkid -laio -o lvm -lpthread -luuid ./liblvm2cmd.a  

# Build scripted fdisk (sfdisk)
FROM gcc:latest as sfdisk
RUN apt-get update -y; apt-get install -y bison autopoint gettext
RUN git clone https://github.com/karelzak/util-linux.git
WORKDIR util-linux
RUN ./autogen.sh && ./configure --enable-static-programs=sfdisk && make


# Build BOOTy as an init
FROM golang:1.14-alpine as dev
RUN apk add --no-cache git ca-certificates gcc linux-headers musl-dev
COPY . /go/src/github.com/thebsdbox/BOOTy/
WORKDIR /go/src/github.com/thebsdbox/BOOTy
ENV GO111MODULE=on
RUN --mount=type=cache,sharing=locked,id=gomod,target=/go/pkg/mod/cache \
    --mount=type=cache,sharing=locked,id=goroot,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o init
    
#RUN go get; CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o init

# Build Busybox
FROM gcc:latest as Busybox
RUN apt-get update; apt-get install -y cpio
RUN curl -O https://busybox.net/downloads/busybox-1.31.1.tar.bz2
RUN tar -xf busybox*bz2
WORKDIR busybox-1.31.1
RUN make defconfig; make LDFLAGS=-static CONFIG_PREFIX=./initramfs install

#RUN make LDFLAGS=-static 
WORKDIR initramfs 
RUN wget -qO- https://launchpad.net/cloud-utils/trunk/0.31/+download/cloud-utils-0.31.tar.gz  | tar -xvz -C /tmp; mv /tmp/cloud-utils-0.31/bin/growpart ./bin

# Copy build contents from previous build
COPY --from=LVM /LVM2.2.03.09/tools/lvm sbin
COPY --from=sfdisk /util-linux/sfdisk.static bin/sfdisk
COPY --from=dev /go/src/github.com/thebsdbox/BOOTy/init .

# Package initramfs
RUN find . -print0 | cpio --null -ov --format=newc > ../initramfs.cpio 
RUN gzip ../initramfs.cpio
RUN mv ../initramfs.cpio.gz /

FROM scratch
COPY --from=Busybox /initramfs.cpio.gz .

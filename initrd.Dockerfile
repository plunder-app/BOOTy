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

FROM golang:1.13-alpine as dev
RUN apk add --no-cache git ca-certificates gcc linux-headers musl-dev
COPY . /go/src/github.com/thebsdbox/BOOTy/
WORKDIR /go/src/github.com/thebsdbox/BOOTy
RUN go get; CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o init

FROM gcc:latest as Busybox
RUN curl -O https://busybox.net/downloads/busybox-1.31.1.tar.bz2
RUN tar -xf busybox*bz2
WORKDIR busybox-1.31.1
RUN make defconfig; make LDFLAGS=-static
#RUN 
#RUN mkdir -p initramfs/bin
#RUN mkdir -p initramfs/dev
#RUN mkdir -p initramfs/etc
#RUN mkdir -p initramfs/home
##RUN mkdir -p initramfs/proc
#RUN mkdir -p initramfs/mnt
#RUN mkdir -p initramfs/sys
#RUN mkdir -p initramfs/usr

RUN make LDFLAGS=-static CONFIG_PREFIX=./initramfs install
WORKDIR initramfs 

#https://www.tldp.org/LDP/lfs/LFS-BOOK-6.1.1-HTML/chapter06/devices.html
#RUN mknod dev/sda b 8 0 
#RUN mknod dev/console c 5 1
#RUN mknod -m 0666 dev/null c 1 3
#RUN mknod -m 0444 dev/random c 1 8
#RUN mknod -m 0444 dev/urandom c 1 9
#RUN cp ../examples/udhcp/simple.script bin
COPY --from=LVM /LVM2.2.03.09/tools/lvm sbin
COPY --from=dev /go/src/github.com/thebsdbox/BOOTy/init .

# Begin building the init (could be code to do this)
#RUN echo '#!/bin/sh' > ./init
#RUN echo "mount -t proc none /proc" >> ./init
#RUN echo "mount -t sysfs none /sys" >> ./init

#RUN echo "echo /sbin/mdev > /proc/sys/kernel/hotplug" >> ./init
#echo "mdev -s" >> ./init
#RUN echo "ifconfig eth0 up" >> ./init
#RUN echo "udhcpc -t 5 -q -s /bin/simple.script" >> ./init
#RUN echo "exec /bin/sh" >> ./init

#RUN chmod +x init

#RUN mv ./booty ./init
RUN apt-get update; apt-get install -y cpio
RUN find . -print0 | cpio --null -ov --format=newc > ../initramfs.cpio 
RUN gzip ../initramfs.cpio
RUN mv ../initramfs.cpio.gz /
#RUN tar -cvzf ../initramfs.tar.gz .
#cd ..
#echo "RAM DISK CREATED"
#echo ""
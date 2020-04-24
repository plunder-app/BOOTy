# BOOTy
A simple linuxKit based kernel/initrd that is used by plunder for deployment managment

## BOOTy build
```
docker build -t plndr/booty:<version>
```

# Creating a ramdisk

## Initial toolchain

`sudo apt-get install -y libaio-dev gcc build-essential libncurses5-dev`

## LVM
```
wget https://mirrors.kernel.org/sourceware/lvm2/LVM2.2.03.09.tgz
tar -xf LVM2.2.03.09.tgz
cd LVM2.2.03.09
./configure --enable-static_link
cd ..
```

## Busybox and initrd
```
curl -O https://busybox.net/downloads/busybox-1.31.1.tar.bz2
tar -xf busybox*bz2
cd busybox-1.26.2
make defconfig
make clean
echo "Building static busybox"
sleep 5
make LDFLAGS=-static
mkdir -p initramfs/{bin,dev,etc,home,mnt,proc,sys,usr}
make LDFLAGS=-static CONFIG_PREFIX=./initramfs install
cd initramfs 

sudo mknod dev/sda b 8 0 
sudo mknod dev/console c 5 1
sudo mknod -m 0666 dev/null c 1 3
sudo mknod -m 0444 dev/random c 1 8
sudo mknod -m 0444 dev/urandom c 1 9



cp ../examples/udhcp/simple.script bin
cp ../../LVM*/tools/lvm sbin/


echo '#!/bin/sh' > ./init
echo "mount -t proc none /proc" >> ./init
echo "mount -t sysfs none /sys" >> ./init

echo "echo /sbin/mdev > /proc/sys/kernel/hotplug" >> ./init
#echo "mdev -s" >> ./init
echo "ifconfig eth0 up" >> ./init
echo "udhcpc -t 5 -q -s /bin/simple.script" >> ./init
echo "exec /bin/sh" >> ./init
chmod +x init
find . -print0 | cpio --null -ov --format=newc > initramfs.cpio 
gzip ./initramfs.cpio
mv ./initramfs.cpio.gz ..
cd ..
echo "RAM DISK CREATED"
echo ""
```

## Linuxkit build (being removed)

## Linuxkit build

### Build pusher

```
../../linuxkit/linuxkit/bin/linuxkit build ./linuxkit/pull.yaml
```

### Build puller

```
../../linuxkit/linuxkit/bin/linuxkit build ./linuxkit/pull.yaml
```

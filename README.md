# BOOTy
A simple initrd that is used by plunder for Operating System image deployment. 

It should go without saying that this it an early version of this software. It comes with **no guard rails** and if used incorrectly could break an existing Operating System

## Example deployment

[![asciicast](https://asciinema.org/a/326011.svg)](https://asciinema.org/a/326011)

## BOOTy build

At the moment the most simple method of building `BOOTy` is to use the `initrd.Dockerfile` to build all the components that are required and compile in `BOOTy` as the init process.

```
docker build -t init -f ./initrd.Dockerfile . ; \
docker run init:latest tar -cf - /initramfs.cpio.gz | tar xf -   
```

**to-do** Mulit-arch builds may work with something like the following:

` docker buildx build  --platform linux/amd64 -o local -t init -f ./initrd.Dockerfile . ; \
docker run init:latest tar -cf - /initramfs.cpio.gz | tar xf -   `

The above command will build these components:

- BusyBox
- LVM
- BOOTy

It will then produce a simple `initramfs` that can be booted with a kernel and then finally it will copy the new `initrams` from the Docker image to the local file system.

## BOOTy boot

Create a boot configuration (the below example uses `plunder`/[plndr.io](plndr.io)):

`pldrctl create boot -i initramfs.cpio.gz -k kernel -c "console=tty0 console=ttyS0,9600" -n booty`

Create a deployment configuration:

`pldrctl create deployment -a a -m 00:50:56:a5:0e:0f -c booty`


## Example Server

Until the server components are implemented into [plndr.io](plndr.io) the server is an external component built for testing.

The two actions dictate the direction of Operating system images. 

The `writeImage` action will instruct the new server on boot to pull the `-sourceImage` and write the contents to the `-destinationDevice`. 


```
go run server/server.go -action writeImage \
-address 00:50:56:a5:0e:0f \
-sourceImage http://192.168.0.95:3000/images/ubuntu.img \
-destinationDevice /dev/sda    
```

The `readImage` action should be used when network booting a server that already has an Operating System installed. The `-destinationAddress` should be the address of the machine that is running the server and should be in the format `http://<address>/image` as the `/image` is a specific handler for receiving the disk image.

```
go run server/server.go -action readImage \
-address 00:50:56:a5:0e:0f \
-destinationAddress http://192.168.0.95:3000/image \
-sourceDevice /dev/sda    
```
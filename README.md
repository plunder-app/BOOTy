# BOOTy
A simple initrd that is used by plunder for Operating System image deployment

## BOOTy build

At the moment the most simple method of building `BOOTy` is to use the `initrd.Dockerfile` to build all the components that are required and compile in `BOOTy` as the init process.

The below command will build:

- BusyBox
- LVM
- BOOTy

It will then produce a simple `initramfs` that can be booted with a kernel.
```
docker build -t init -f ./initrd.Dockerfile . ; \
docker run init:latest tar -cf - /initramfs.cpio.gz | tar xf -   
```

## BOOTy boot

Create a boot configuration (the below example uses `plunder`):

`pldrctl create boot -i initramfs.cpio.gz -k kernel -c "console=tty0 console=ttyS0,9600" -n booty`

Create a deployment configuration:

`pldrctl create deployment -a a -m 00:50:56:a5:0e:0f -c booty`

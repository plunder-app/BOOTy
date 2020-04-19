# BOOTy
A simple linuxKit based kernel/initrd that is used by plunder for deployment managment

## BOOTy build
```
docker build -t plndr/booty:<version>
```

## Linuxkit build

### Build pusher

```
../../linuxkit/linuxkit/bin/linuxkit build ./linuxkit/pull.yaml
```

### Build puller

```
../../linuxkit/linuxkit/bin/linuxkit build ./linuxkit/pull.yaml
```

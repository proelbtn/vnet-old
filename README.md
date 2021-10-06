# vnet - Virtual Network Laboratory

## Install

### Build from source

```
git clone https://github.com/proelbtn/vnet vnet && cd $_
docker run --rm -itv ${PWD}:/work -w /work golang:1.17 go build -o vnet ./cmd/vnet/*.go
mv vnet /usr/local/bin
```
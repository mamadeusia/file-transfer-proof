# grpc-filetransfer


## Server 



```docker compose up -d```

Server config `./config/server/config.yml` or you can use envs.

## Client

```
./client -a=':9000' -f /directory
```

for downloading need to specify the index and the merkle root of files 
```
./client download -i0 -c='ae8aa01d64cdbecc7f5091cc1b68c6ae9969de67cf51dad3a6a28728b5b4809f' -a=':9000'
```



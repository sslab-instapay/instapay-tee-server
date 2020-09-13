# Go Server for InstaPay

## sgx environment
```sh
source $SGX_SDK/environment
```

## environment variable
```sh
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$GOPATH/src/github.com/sslab-instapay/instapay-tee-server
```

## go-ethereum
```sh
geth --datadir . --networkid 3333 --rpc --rpcaddr 141.223.121.139 --rpcport 8555 --ws --wsaddr 141.223.121.139 --wsport 8881 --wsorigins="*" --port 30303 --rpccorsdomain "*" --rpcapi "db,eth,net,web3,personal,admin,miner,debug,txpool" --wsapi "db,eth,net,web3,personal,admin,miner,debug,txpool" --nodiscover console
```

## run
```sh
go run main.go -port=3004 -grpc_port=50004
```

## 실행 방법
start.sh 실행

```sh
source ~/sgxsdk/environment
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/xiaofo/instapay/src/github.com/sslab-instapay/instapay-tee-server
go run main.go -port=3004 -grpc_port=50004
```

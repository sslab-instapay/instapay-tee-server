source ~/sgxsdk/environment
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/xiaofo/instapay/src/github.com/sslab-instapay/instapay-tee-server
go run main.go -port=3004 -grpc_port=50004
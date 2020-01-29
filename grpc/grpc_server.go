package grpc

/*
#cgo CPPFLAGS: -I/home/xiaofo/sgxsdk/include -I/home/xiaofo/instapay/src/github.com/sslab-instapay/instapay-tee-server
#cgo LDFLAGS: -L/home/xiaofo/instapay/src/github.com/sslab-instapay/instapay-tee-server -ltee

#include "../app.h"
*/
import "C"

import (
	"net"
	"log"
	"fmt"
	"os"
	"sync"
	"context"
	"strconv"
	"google.golang.org/grpc"
	"github.com/sslab-instapay/instapay-tee-server/repository"
	pbServer "github.com/sslab-instapay/instapay-tee-server/proto/server"
	pbClient "github.com/sslab-instapay/instapay-tee-server/proto/client"
	"unsafe"
	"reflect"
)


type ServerGrpc struct {
}


var rwMutex = new(sync.RWMutex)


func SendAgreementRequest(pn int64, address string, w pbClient.AgreeRequestsMessage) {
	info, err := repository.GetClientInfo(address)
	if err != nil {
		log.Fatal(err)
	}

	clientAddr := (*info).IP + ":" + strconv.Itoa((*info).Port)
	conn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pbClient.NewClientClient(conn)
	var channelIds []int
	var amounts []int
	for _, channelPayment := range w.ChannelPayments.GetChannelPayments(){
		channelIds = append(channelIds, int(channelPayment.ChannelId))
		amounts = append(amounts, int(channelPayment.Amount))
	}
	var channelSlice []C.uint

	for i := range channelIds{
		channelSlice = append(channelSlice, C.uint(i))
	}

	var amountSlice []C.int

	for i := range amounts{
		amountSlice = append(amountSlice, C.int(i))
	}
	var originalMessage *C.uchar
	var signature *C.uchar

	//void ecall_create_ag_req_msg_w(unsigned int payment_num, unsigned int payment_size, unsigned int *channel_ids, int *amount, unsigned char **original_msg, unsigned char **output);
	C.ecall_create_ag_req_msg_w(C.uint(pn), C.uint(len(channelSlice)), &channelSlice[0], &amountSlice[0], &originalMessage, &signature)

	originalMessageByte, signatureByte := convertPointerToByte(originalMessage, signature)


	r, err := client.AgreementRequest(context.Background(), &pbClient.AgreeRequestsMessage{PaymentNumber: int64(pn), OriginalMessage: originalMessageByte, Signature: signatureByte})
	if err != nil {
		log.Println(err)
	}

	if r.Result{
		agreementOriginalMessage, agreementSignature := convertByteToPointer(r.OriginalMessage, r.Signature)
		C.ecall_verify_ag_res_msg_w(address, agreementOriginalMessage, agreementSignature)
	}

	rwMutex.Lock()
	C.ecall_update_sentagr_list_w(C.uint(pn), &([]C.uchar(address)[0]))
	rwMutex.Unlock()

	//fmt.Println("AGREED: " + address)

	return
}


func SendUpdateRequest(pn int64, address string, w pbClient.AgreeRequestsMessage) {
	info, err := repository.GetClientInfo(address)
	if err != nil {
		log.Fatal(err)
	}

	clientAddr := (*info).IP + ":" + strconv.Itoa((*info).Port)
	conn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pbClient.NewClientClient(conn)

	var channelIds []int
	var amounts []int
	for _, channelPayment := range w.ChannelPayments.GetChannelPayments(){
		channelIds = append(channelIds, int(channelPayment.ChannelId))
		amounts = append(amounts, int(channelPayment.Amount))
	}
	var channelSlice []C.uint

	for i := range channelIds{
		channelSlice = append(channelSlice, C.uint(i))
	}

	var amountSlice []C.int

	for i := range amounts{
		amountSlice = append(amountSlice, C.int(i))
	}
	var originalMessage *C.uchar
	var signature *C.uchar

	// void ecall_create_ud_req_msg_w(unsigned int payment_num, unsigned int payment_size, unsigned int *channel_ids, int *amount, unsigned char **original_msg, unsigned char **output)
	C.ecall_create_ud_req_msg_w(C.uint(pn), C.uint(len(channelSlice)), &channelSlice[0], &amountSlice[0], &originalMessage, &signature)

	originalMessageByte, signatureByte := convertPointerToByte(originalMessage, signature)

	rqm := pbClient.UpdateRequestsMessage{		/* convert AgreeRequestsMessage to UpdateRequestsMessage */
		PaymentNumber:   w.PaymentNumber,
		ChannelPayments: w.ChannelPayments,
		OriginalMessage: originalMessageByte,
		Signature: signatureByte,
		}

	r, err := client.UpdateRequest(context.Background(), &rqm)
	if err != nil {
		log.Fatal(err)
	}

	if r.Result{
		updateOriginalMessage, updateSignature := convertByteToPointer(r.OriginalMessage, r.Signature)
		C.ecall_verify_ud_res_msg_w(&([]C.uchar(address)[0]), updateOriginalMessage, updateSignature)
	}

	rwMutex.Lock()
	C.ecall_update_sentupt_list_w(C.uint(pn), &([]C.uchar(address))[0])
	rwMutex.Unlock()
	
	return
}


func SendConfirmPayment(pn int, address string) {
	info, err := repository.GetClientInfo(address)
	if err != nil {
		log.Fatal(err)
	}

	clientAddr := (*info).IP + ":" + strconv.Itoa((*info).Port)
	conn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pbClient.NewClientClient(conn)

	var originalMessage *C.uchar
	var signature *C.uchar
	C.ecall_create_confirm_msg_w(C.uint(int32(pn)), &originalMessage, &signature)

	originalMessageByte, signatureByte := convertPointerToByte(originalMessage, signature)

	_, err = client.ConfirmPayment(context.Background(), &pbClient.ConfirmRequestsMessage{PaymentNumber: int64(pn), OriginalMessage: originalMessageByte, Signature: signatureByte},)
	if err != nil {
		log.Println(err)
	}
}


func WrapperAgreementRequest(pn int64, p []string, w map[string]pbClient.AgreeRequestsMessage) {
	/* remove C's address from p */
	var q []string
	q = p[0:8]

	for _, address := range q {
		go SendAgreementRequest(pn, address, w[address])
	}

	for C.ecall_check_unanimity_w(C.uint(pn), C.int(0)) != 1 {}

	go WrapperUpdateRequest(pn, p, w)
}


func WrapperUpdateRequest(pn int64, p []string, w map[string]pbClient.AgreeRequestsMessage) {
	for _, address := range p {
		go SendUpdateRequest(pn, address, w[address])
	}

	for C.ecall_check_unanimity_w(C.uint(pn), C.int(1)) != 1 {}

	go WrapperConfirmPayment(int(pn), p)
}


func WrapperConfirmPayment(pn int, p []string) {
	/* update payment's status */
	C.ecall_update_payment_status_to_success_w(C.uint(pn))

	for _, address := range p {
		go SendConfirmPayment(pn, address)
	}
	fmt.Println("SENT CONFIRMATION")
}


func SearchPath(pn int64, amount int64) ([]string, map[string]pbClient.AgreeRequestsMessage) {
	var p []string

	/* composing p */
	for i := 1; i < 9; i++ {
		t := i / 10
		o := i % 10
		s := "00000000000000000000000000000000000000" + strconv.Itoa(t) + strconv.Itoa(o)		
		p = append(p, s)
	}
	p = append(p, "0000000000000000000000000000000000000009")

	/* composing w */
	amount = int64(amount)
	pn = int64(pn)

	var w map[string]pbClient.AgreeRequestsMessage
	w = make(map[string]pbClient.AgreeRequestsMessage)

	var cpsf []*pbClient.ChannelPayment
	cpsf = append(cpsf, &pbClient.ChannelPayment{ChannelId: int64(1), Amount: -amount})
	rqmf := pbClient.AgreeRequestsMessage{
		PaymentNumber:   pn,
		ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cpsf},
		}
	w["0000000000000000000000000000000000000001"] = rqmf

	for i:= 2; i < 9; i++ {
		var cps []*pbClient.ChannelPayment
		cps = append(cps, &pbClient.ChannelPayment{ChannelId: int64(i), Amount: amount})
		cps = append(cps, &pbClient.ChannelPayment{ChannelId: int64(i + 1), Amount: -amount})
		rqm := pbClient.AgreeRequestsMessage{
			PaymentNumber:   pn,
			ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps},
			}
		t := i / 10
		o := i % 10
		s := "00000000000000000000000000000000000000" + strconv.Itoa(t) + strconv.Itoa(o)			
		w[s] = rqm
	}

	var cpsl []*pbClient.ChannelPayment
	cpsl = append(cpsl, &pbClient.ChannelPayment{ChannelId: int64(89), Amount: amount})
	rqml := pbClient.AgreeRequestsMessage{
		PaymentNumber:   pn,
		ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cpsl},
		}
	w["0000000000000000000000000000000000000009"] = rqml

	return p, w
}


func (s *ServerGrpc) PaymentRequest(ctx context.Context, rq *pbServer.PaymentRequestMessage) (*pbServer.Result, error) {
	from := rq.From
	to := rq.To
	amount := rq.Amount

	sender := []C.uchar(from)
	receiver := []C.uchar(to)

	PaymentNum := C.ecall_accept_request_w(&sender[0], &receiver[0], C.uint(amount))
	p, w := SearchPath(int64(PaymentNum), amount)

	for i := 0; i < len(p); i++ {
		C.ecall_add_participant_w(PaymentNum, &([]C.uchar(p[i]))[0])
	}
	C.ecall_update_sentagr_list_w(PaymentNum, &([]C.uchar(p[8]))[0])

	go WrapperAgreementRequest(int64(PaymentNum), p, w)

	return &pbServer.Result{Result: true}, nil
}


func (s *ServerGrpc) CommunicationInfoRequest(ctx context.Context, address *pbServer.Address) (*pbServer.CommunicationInfo, error) {
	res, err := repository.GetClientInfo(address.Addr)
	if err != nil {
		log.Fatal(err)
	}

	return &pbServer.CommunicationInfo{IPAddress: res.IP, Port: int64(res.Port)}, nil
}


func StartGrpcServer() {
	grpcPort, err := strconv.Atoi(os.Getenv("grpc_port"))
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pbServer.RegisterServerServer(grpcServer, &ServerGrpc{})

	grpcServer.Serve(lis)
}

func convertByteToPointer(originalMsg []byte, signature []byte) (*C.uchar, *C.uchar){

	var uOriginal [44]C.uchar
	var uSignature [65]C.uchar

	for i := 0; i < 44; i++{
		uOriginal[i] = C.uchar(originalMsg[i])
	}

	for i := 0; i < 65; i++{
		uSignature[i] = C.uchar(signature[i])
	}

	cOriginalMsg := (*C.uchar)(unsafe.Pointer(&uOriginal[0]))
	cSignature := (*C.uchar)(unsafe.Pointer(&uSignature[0]))

	return cOriginalMsg, cSignature
}

func convertPointerToByte(originalMsg *C.uchar, signature *C.uchar)([]byte, []byte){

	var returnMsg []byte
	var returnSignature []byte

	replyMsgHdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(originalMsg)),
		Len: int(44),
		Cap: int(44),
	}
	replyMsgS := *(*[]C.uchar)(unsafe.Pointer(&replyMsgHdr))

	replySigHdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(signature)),
		Len: int(65),
		Cap: int(65),
	}
	replySigS := *(*[]C.uchar)(unsafe.Pointer(&replySigHdr))

	for i := 0; i < 44; i++{
		returnMsg = append(returnMsg, byte(replyMsgS[i]))
	}

	for i := 0; i < 65; i++{
		returnSignature = append(returnSignature, byte(replySigS[i]))
	}

	return returnMsg, returnSignature
}

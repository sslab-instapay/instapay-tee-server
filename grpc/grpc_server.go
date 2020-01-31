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
	for _, i := range channelIds{
		channelSlice = append(channelSlice, C.uint(i))
	}

	var amountSlice []C.int
	for _, i := range amounts{
		amountSlice = append(amountSlice, C.int(i))
	}

	var originalMessage *C.uchar
	var signature *C.uchar
	
	C.ecall_create_ag_req_msg_w(C.uint(pn), C.uint(len(channelSlice)), &channelSlice[0], &amountSlice[0], &originalMessage, &signature)

	/************** for debugging **************/
	fmt.Println("==== [TO] %s (original msg) ==========", clientAddr)
	var sig *C.uchar = originalMessage
	hdr := reflect.SliceHeader{
	   Data: uintptr(unsafe.Pointer(sig)),
	   Len:  int(44),
	   Cap:  int(44),
	}
 
	s := *(*[]C.uchar)(unsafe.Pointer(&hdr))
	for i := 0; i < 44; i++ {
	   fmt.Printf("%02x", s[i])
	}
	fmt.Println()
	fmt.Println("====================================================")

	fmt.Println("==== [TO] %s (signature) ==========", clientAddr)
	var sig2 *C.uchar = signature
	hdr2 := reflect.SliceHeader{
	   Data: uintptr(unsafe.Pointer(sig2)),
	   Len:  int(65),
	   Cap:  int(65),
	}
 
	s2 := *(*[]C.uchar)(unsafe.Pointer(&hdr2))
	for i := 0; i < 65; i++ {
	   fmt.Printf("%02x", s2[i])
	}
	fmt.Println()
	fmt.Println("====================================================")
	/*******************************************/

	originalMessageByte, signatureByte := convertPointerToByte(originalMessage, signature)
	r, err := client.AgreementRequest(context.Background(), &pbClient.AgreeRequestsMessage{PaymentNumber: int64(pn), OriginalMessage: originalMessageByte, Signature: signatureByte})
	if err != nil {
		log.Println(err)
	}

	/************** for debugging **************/
	fmt.Printf("========= [FROM] %s ==============================\n", clientAddr)
	oo := r.OriginalMessage[:]
	for i := 0; i < 44; i++ {
		fmt.Printf("%02x", oo[i])
	}
	fmt.Println()

	ss := r.Signature[:]
	for i := 0; i < 65; i++ {
		fmt.Printf("%02x", ss[i])
	}
	fmt.Println()
	fmt.Println("==============================================================")
	/*******************************************/

	if r.Result{
		agreementOriginalMessage, agreementSignature := convertByteToPointer(r.OriginalMessage, r.Signature)
		C.ecall_verify_ag_res_msg_w(&([]C.uchar(address)[0]), agreementOriginalMessage, agreementSignature)
	}

	rwMutex.Lock()
	C.ecall_update_sentagr_list_w(C.uint(pn), &([]C.uchar(address)[0]))
	rwMutex.Unlock()

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
	q = p[0:2]

	for _, address := range q {
		go SendAgreementRequest(pn, address, w[address])
	}

	for C.ecall_check_unanimity_w(C.uint(pn), C.int(0)) != 1 {}

	fmt.Println("[ALARM] ALL USERS AGREED")

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

	p = append(p, "d03a2cc08755ec7d75887f0997195654b928893e")
	p = append(p, "0b4161ad4f49781a821c308d672e6c669139843c")
	p = append(p, "78902c58006916201f65f52f7834e467877f0500")

	/* composing w */

	amount = int64(amount)
	pn = int64(pn)

	var w map[string]pbClient.AgreeRequestsMessage
	w = make(map[string]pbClient.AgreeRequestsMessage)

	var cps1 []*pbClient.ChannelPayment
	cps1 = append(cps1, &pbClient.ChannelPayment{ChannelId: int64(8), Amount: -amount})
	rqm1 := pbClient.AgreeRequestsMessage{
		PaymentNumber:   pn,
		ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps1},
		}
	w["d03a2cc08755ec7d75887f0997195654b928893e"] = rqm1

	var cps2 []*pbClient.ChannelPayment
	cps2 = append(cps2, &pbClient.ChannelPayment{ChannelId: int64(8), Amount: amount})
	cps2 = append(cps2, &pbClient.ChannelPayment{ChannelId: int64(12), Amount: -amount})
	rqm2 := pbClient.AgreeRequestsMessage{
		PaymentNumber:   pn,
		ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps2},
		}
	w["0b4161ad4f49781a821c308d672e6c669139843c"] = rqm2

	var cps3 []*pbClient.ChannelPayment
	cps3 = append(cps3, &pbClient.ChannelPayment{ChannelId: int64(12), Amount: amount})
	rqm3 := pbClient.AgreeRequestsMessage{
		PaymentNumber:   pn,
		ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps3},
		}
	w["78902c58006916201f65f52f7834e467877f0500"] = rqm3

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
	C.ecall_update_sentagr_list_w(PaymentNum, &([]C.uchar(p[2]))[0])

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

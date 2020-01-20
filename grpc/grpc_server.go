package grpc

/*
#cgo CPPFLAGS: -I/home/xiaofo/sgxsdk/include -I/home/xiaofo/instapay/src/github.com/sslab-instapay/instapay-go-server
#cgo LDFLAGS: -L/home/xiaofo/instapay/src/github.com/sslab-instapay/instapay-go-server -ltee

#include "app.h"
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
	// TODO: ecall_create_ag_req_msg_w를 호출하여, agreement request 메시지와 서명 생성
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
		amountSlice = append(amountSlice, C.uint(i))
	}
	var originalMessage *C.uchar
	var signature *C.uchar

	//void ecall_create_ag_req_msg_w(unsigned int payment_num, unsigned int payment_size, unsigned int *channel_ids, int *amount, unsigned char **original_msg, unsigned char **output);
	C.ecall_create_ag_req_msg_w(C.uint(pn), C.uint(len(channelSlice)), &channelSlice[0], &amountSlice[0], originalMessage, signature)

	hdr1 := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(originalMessage)),
		Len:  int(44),
		Cap:  int(44),
	}
	s1 := *(*[]C.uchar)(unsafe.Pointer(&hdr1))

	hdr2 := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(signature)),
		Len:  int(65),
		Cap:  int(65),
	}
	s2 := *(*[]C.uchar)(unsafe.Pointer(&hdr2))
	originalMessageStr := fmt.Sprintf("02x", s1)
	signatureStr := fmt.Sprintf("02x", s2)


	r, err := client.AgreementRequest(context.Background(), &pbClient.AgreeRequestsMessage{PaymentNumber: int64(pn), OriginalMessage: originalMessageStr, Signature: signatureStr})  // TODO: byte stream으로 메시지와 서명을 클라이언트에게 전달
	if err != nil {
		log.Println(err)
	}

	//unsigned int ecall_verify_ag_res_msg_w(unsigned char *pubaddr, unsigned char *res_msg, unsigned char *res_sig);
	if r.Result{
		C.ecall_verify_ag_res_msg_w(address, &([]C.uchar(originalMessage)[0]), &([]C.uchar(signature)[0]))
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
		amountSlice = append(amountSlice, C.uint(i))
	}
	var originalMessage *C.uchar
	var signature *C.uchar

	// TODO: ecall_create_ud_req_msg_w를 호출하여, update request 메시지와 서명 생성
	// void ecall_create_ud_req_msg_w(unsigned int payment_num, unsigned int payment_size, unsigned int *channel_ids, int *amount, unsigned char **original_msg, unsigned char **output)
	C.ecall_create_ud_req_msg_w(C.uint(pn), C.uint(len(channelSlice)), &channelSlice[0], &amountSlice[0], originalMessage, signature)

	hdr1 := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(originalMessage)),
		Len:  int(44),
		Cap:  int(44),
	}
	s1 := *(*[]C.uchar)(unsafe.Pointer(&hdr1))

	hdr2 := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(signature)),
		Len:  int(65),
		Cap:  int(65),
	}
	s2 := *(*[]C.uchar)(unsafe.Pointer(&hdr2))
	originalMessageStr := fmt.Sprintf("02x", s1)
	signatureStr := fmt.Sprintf("02x", s2)

	rqm := pbClient.UpdateRequestsMessage{		/* convert AgreeRequestsMessage to UpdateRequestsMessage */
		PaymentNumber:   w.PaymentNumber,
		ChannelPayments: w.ChannelPayments,
		OriginalMessage: originalMessageStr,
		Signature: signatureStr,
			}

	r, err := client.UpdateRequest(context.Background(), &rqm)	// TODO: byte stream으로 메시지와 서명을 클라이언트에게 전달
	if err != nil {
		log.Fatal(err)
	}

	if r.Result{
		C.ecall_verify_ud_res_msg_w(&([]C.uchar(address)[0]), &([]C.uchar(r.OriginalMessage)[0]), &([]C.uchar(r.Signature)[0]))
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

	// TODO: ecall_craete_confirm_msg_w를 호출하여, payment confirm 메시지와 서명 생성
	var originalMessage *C.uchar
	var signature *C.uchar
	C.ecall_create_confirm_msg_w(C.uint(int32(pn)), originalMessage, signature)

	hdr1 := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(signature)),
		Len:  int(44),
		Cap:  int(44),
	}
	s1 := *(*[]C.uchar)(unsafe.Pointer(&hdr1))

	hdr2 := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(signature)),
		Len:  int(65),
		Cap:  int(65),
	}
	s2 := *(*[]C.uchar)(unsafe.Pointer(&hdr2))
	originalMessageStr := fmt.Sprintf("02x", s1)
	signatureStr := fmt.Sprintf("02x", s2)

	_, err = client.ConfirmPayment(context.Background(), &pbClient.ConfirmRequestsMessage{PaymentNumber: int64(pn), OriginalMessage: originalMessageStr, Signature: signatureStr},)	// TODO: byte stream으로 메시지와 서명을 클라이언트에게 전달
	if err != nil {
		log.Fatal(err)
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

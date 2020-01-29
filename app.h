#ifndef _APP_H_
#define _APP_H_

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>

#include "sgx_error.h"       /* sgx_status_t */
#include "sgx_eid.h"     /* sgx_enclave_id_t */

#ifndef TRUE
# define TRUE 1
#endif

#ifndef FALSE
# define FALSE 0
#endif

# define TOKEN_FILENAME   "enclave.token"
# define ENCLAVE_FILENAME "enclave.signed.so"

extern sgx_enclave_id_t global_eid;    /* global enclave id */

#if defined(__cplusplus)
extern "C" {
#endif

typedef struct _message {
    /********* common *********/
    unsigned int type;

    /***** direct payment *****/
    unsigned int channel_id;
    int amount;
    unsigned int counter;

    /*** multi-hop payment ****/
    unsigned int payment_num;
    unsigned int payment_size;
    unsigned int channel_ids[2];
    int payment_amount[2];
    unsigned int e;
} message;

int initialize_enclave(void);

unsigned int ecall_accept_request_w(unsigned char *sender, unsigned char *receiver, unsigned int amount);
void ecall_add_participant_w(unsigned int payment_num, unsigned char *addr);
void ecall_update_sentagr_list_w(unsigned int payment_num, unsigned char *addr);
void ecall_update_sentupt_list_w(unsigned int payment_num, unsigned char *addr);
int ecall_check_unanimity_w(unsigned int payment_num, int which_list);
void ecall_update_payment_status_to_success_w(unsigned int payment_num);

/** 서버의 agreement request 메시지와 서명을 생성
 *
 * Out:     original_msg:   생성된 메시지의 plain text 주소
 *          output:         생성된 메시지의 signature 주소
 * In:      payment_num:    서버가 생성한 payment instance 번호
 *          payment_size:   payment data 길이 (channel_ids의 길이가 2이면, 값은 2)
 *                          길이는 반드시 1 또는 2임
 *          channel_ids:    해당 payment instance와 관련된 클라이언트의 채널 id 집합
 *          amount:         channel_ids에 포함된 id들과 매칭되는 각 payment amount
 *                          만약, 채널 A에 3을 지불해야 하면 -3이 되며, 지불받으면 +3이 됨
 */
void ecall_create_ag_req_msg_w(unsigned int payment_num, unsigned int payment_size, unsigned int *channel_ids, int *amount, unsigned char **original_msg, unsigned char **output);

/** 서버의 update request 메시지와 서명을 생성
 *
 * Out:     original_msg:   생성된 메시지의 plain text 주소
 *          output:         생성된 메시지의 signature 주소
 * In:      payment_num:    서버가 생성한 payment instance 번호
 *          payment_size:   payment data 길이 (channel_ids의 길이가 2이면, 값은 2)
 *                          길이는 반드시 1 또는 2임
 *          channel_ids:    해당 payment instance와 관련된 클라이언트의 채널 id 집합
 *          amount:         channel_ids에 포함된 id들과 매칭되는 각 payment amount
 *                          만약, 채널 A에 3을 지불해야 하면 -3이 되며, 지불받으면 +3이 됨
 */
void ecall_create_ud_req_msg_w(unsigned int payment_num, unsigned int payment_size, unsigned int *channel_ids, int *amount, unsigned char **original_msg, unsigned char **output);

/** 서버의 payment confirm 메시지와 서명을 생성
 *
 * Out:     original_msg:   생성된 메시지의 plain text 주소
 *          output:         생성된 메시지의 signature 주소
 * In:      payment_num:    서버가 생성한 payment instance 번호
 */
void ecall_create_confirm_msg_w(unsigned int payment_num, unsigned char **original_msg, unsigned char **output);

/** 클라이언트가 보낸 agreement response의 메시지 서명을 검증
 *
 * Returns: 1 이면 검증 성공, 0 이면 검증 실패
 * In:      pubaddr:   검증하려는 클라이언트의 공개 주소
 *          res_msg:   클라이언트의 agreement response 메시지의 plain text 주소
 *          res_sig:   클라이언트의 agreement response 메시지의 signature 주소
 */
unsigned int ecall_verify_ag_res_msg_w(unsigned char *pubaddr, unsigned char *res_msg, unsigned char *res_sig);

/** 클라이언트가 보낸 agreement response의 메시지 서명을 검증
 *
 * Returns: 1 이면 검증 성공, 0 이면 검증 실패
 * In:      pubaddr:   검증하려는 클라이언트의 공개 주소
 *          res_msg:   클라이언트의 agreement response 메시지의 plain text 주소
 *          res_sig:   클라이언트의 agreement response 메시지의 signature 주소
 */
unsigned int ecall_verify_ud_res_msg_w(unsigned char *pubaddr, unsigned char *res_msg, unsigned char *res_sig);

#if defined(__cplusplus)
}
#endif

#endif /* !_APP_H_ */

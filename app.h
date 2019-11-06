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

int initialize_enclave(void);

unsigned int ecall_accept_request_w(unsigned char *sender, unsigned char *receiver, unsigned int amount);
void ecall_add_participant_w(unsigned int payment_num, unsigned char *addr);
void ecall_update_sentagr_list_w(unsigned int payment_num, unsigned char *addr);
void ecall_update_sentupt_list_w(unsigned int payment_num, unsigned char *addr);
int ecall_check_unanimity_w(unsigned int payment_num, int which_list);
void ecall_update_payment_status_to_success_w(unsigned int payment_num);

#if defined(__cplusplus)
}
#endif

#endif /* !_APP_H_ */

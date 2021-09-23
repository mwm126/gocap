#include "yk_darwin.h"
#include <stdio.h>

#include <ykcore.h>
#include <ykdef.h>
#include <ykpers-version.h>
#include <ykstatus.h>
#include <yubikey.h>

int get_yk_serial() {
  unsigned int serial = 123;
  YK_KEY *yk = 0;

  if (!yk_init()) {
    printf("Could not initialize yubikey\n");
    return -1;
  }

  yk = yk_open_key(0);
  if (!yk) {
    printf("Could not open yubikey\n");
    return -1;
  }

  yk_get_serial(yk, 1, 0, &serial);
  return serial;
}

static void report_yk_error(void) {
  if (yk_errno) {
    if (yk_errno == YK_EUSBERR) {
      fprintf(stderr, "USB error: %s\n", yk_usb_strerror());
    } else {
      fprintf(stderr, "Yubikey core error: %s\n", yk_strerror(yk_errno));
    }
  }
}

extern int optind;

static int check_firmware(YK_KEY *yk, bool verbose) {
  YK_STATUS *st = ykds_alloc();

  if (!yk_get_status(yk, st)) {
    ykds_free(st);
    return 0;
  }

  if (verbose) {
    printf("Firmware version %d.%d.%d\n", ykds_version_major(st),
           ykds_version_minor(st), ykds_version_build(st));
    fflush(stdout);
  }

  if (ykds_version_major(st) < 2 ||
      (ykds_version_major(st) == 2 && ykds_version_minor(st) < 2)) {
    fprintf(stderr, "Challenge-response not supported before YubiKey 2.2.\n");
    ykds_free(st);
    return 0;
  }

  ykds_free(st);
  return 1;
}

static int otp_challenge_response(YK_KEY *yk, int slot,
                                  unsigned char *challenge, unsigned int len,
                                  bool hmac, bool may_block, bool verbose,
                                  int digits, unsigned char otp[OTP_LENGTH]) {
  unsigned char response[SHA1_MAX_BLOCK_SIZE];
  unsigned char output_buf[(SHA1_MAX_BLOCK_SIZE * 2) + 1];
  int yk_cmd;
  /* unsigned int expect_bytes = 0; */
  memset(response, 0, sizeof(response));
  memset(output_buf, 0, sizeof(output_buf));

  if (verbose) {
    fprintf(stderr, "Sending %i bytes %s challenge to slot %i\n", len,
            (hmac == true) ? "HMAC" : "Yubico", slot);
  }

  switch (slot) {
  case 1:
    yk_cmd = (hmac == true) ? SLOT_CHAL_HMAC1 : SLOT_CHAL_OTP1;
    break;
  case 2:
    yk_cmd = (hmac == true) ? SLOT_CHAL_HMAC2 : SLOT_CHAL_OTP2;
    break;
  default:
    return 0;
  }

  if (!yk_challenge_response(yk, yk_cmd, may_block, len, challenge,
                             sizeof(response), response)) {
    return 0;
  }

  /* HMAC responses are 160 bits, Yubico 128 */
  /* expect_bytes = (hmac == true) ? 20 : 16; */
  for (int ii = 0; ii < OTP_LENGTH; ii++) {
    otp[ii] = response[ii];
  }
  return 1;
}

int get_otp(const unsigned char decoded[], unsigned char otp[]) {
  YK_KEY *yk = 0;
  bool error = true;
  int exit_code = 0;

  /* Options */
  bool verbose = false;
  unsigned char *challenge;

  unsigned int challenge_len;
  int slot = 1;
  int key_index = 0;

  yk_errno = 0;

  /* static unsigned char decoded[6]={212, 89, 194, 77, 162, 249}; */
  /* memset(decoded, 0, sizeof(decoded)); */
  challenge = (unsigned char *)&decoded;
  challenge_len = 6;

  if (!yk_init()) {
    exit_code = 1;
    goto err;
  }

  if (!(yk = yk_open_key(key_index))) {
    exit_code = 1;
    goto err;
  }

  if (!check_firmware(yk, verbose)) {
    exit_code = 1;
    goto err;
  }

  bool hmac = false;
  bool may_block = true;
  int digits = 0;
  if (!otp_challenge_response(yk, slot, challenge, challenge_len, hmac,
                              may_block, verbose, digits, otp)) {
    exit_code = 1;
    goto err;
  }

  exit_code = 0;
  error = false;

err:
  if (error || exit_code != 0) {
    report_yk_error();
  }

  if (yk && !yk_close_key(yk)) {
    report_yk_error();
    exit_code = 2;
  }

  if (!yk_release()) {
    report_yk_error();
    exit_code = 2;
  }

  return exit_code;
}

static int hmac_challenge_response(YK_KEY *yk, int slot,
                                   unsigned char *challenge, unsigned int len,
                                   bool may_block, bool verbose, int digits,
                                   unsigned char *hmac) {
  unsigned char response[SHA1_MAX_BLOCK_SIZE];
  unsigned char output_buf[(SHA1_MAX_BLOCK_SIZE * 2) + 1];
  memset(response, 0, sizeof(response));
  memset(output_buf, 0, sizeof(output_buf));

  if (verbose) {
    fprintf(stderr, "Sending %i bytes %s challenge to slot %i\n", len, "HMAC",
            slot);
  }

  if (!yk_challenge_response(yk, SLOT_CHAL_HMAC2, may_block, len, challenge,
                             sizeof(response), response)) {
    return 0;
  }

  for (int ii = 0; ii < HMAC_LENGTH; ii++) {
    hmac[ii] = response[ii];
  }
  return 1;
}

int hmac_from_digest(const unsigned char digest[DIGEST_LENGTH],
                     unsigned char hmac[HMAC_LENGTH]) {
  YK_KEY *yk = 0;
  bool error = true;
  int exit_code = 0;

  /* Options */
  bool verbose = false;
  /* bool totp = false; */
  int digits = 0;
  unsigned char *challenge;
  unsigned int challenge_len;
  int key_index = 0;

  yk_errno = 0;

  if (!yk_init()) {
    exit_code = 1;
    goto err;
  }

  if (!(yk = yk_open_key(key_index))) {
    exit_code = 1;
    goto err;
  }

  if (!check_firmware(yk, verbose)) {
    exit_code = 1;
    goto err;
  }

  /* unsigned char decoded[SHA1_MAX_BLOCK_SIZE]; */
  challenge = (unsigned char *)&digest;
  challenge_len = 32;
  bool may_block = true;
  if (!hmac_challenge_response(yk, 2, challenge, challenge_len, may_block,
                               verbose, digits, hmac)) {
    exit_code = 1;
    goto err;
  }

  exit_code = 0;
  error = false;

err:
  if (error || exit_code != 0) {
    report_yk_error();
  }

  if (yk && !yk_close_key(yk)) {
    report_yk_error();
    exit_code = 2;
  }

  if (!yk_release()) {
    report_yk_error();
    exit_code = 2;
  }

  return exit_code;
}

/* -*- mode:C; c-file-style: "bsd" -*- */
/*
 * Copyright (c) 2011-2013 Yubico AB.
 * All rights reserved.
 *
 * Author : Fredrik Thulin <fredrik@yubico.com>
 *
 * Some basic code copied from ykpersonalize.c.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *
 *     * Redistributions in binary form must reproduce the above
 *       copyright notice, this list of conditions and the following
 *       disclaimer in the documentation and/or other materials provided
 *       with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

#define __USE_MINGW_ANSI_STDIO 1
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <assert.h>
#include <unistd.h>

#include <ykcore.h>
#include <ykdef.h>
#include <ykstatus.h>
#include <yubikey.h>

#include <ykpers-version.h>

#include "yk.h"

char *_buffer;

const char *usage =
    "Usage: ykchalresp [options] [challenge]\n"
    "\n"
    "Options :\n"
    "\n"
    "\t-nkey     Send challenge to nth key found.\n"
    "\t-1        Send challenge to slot 1. This is the default.\n"
    "\t-2        Send challenge to slot 2.\n"
    "\t-H        Send a 64 byte HMAC challenge. This is the default.\n"
    "\t-Y        Send a 6 byte Yubico challenge.\n"
    "\t-N        Abort if Yubikey requires button press.\n"
    "\t-x        Challenge is hex encoded.\n"
    "\t-t        Time based challenge (for TOTP)\n"
    "\t-6        Output 6 digit HOTP/TOTP code\n"
    "\t-8        Output 8 digit HOTP/TOTP code\n"
    "\t-iFILE    Read challenge from a file instead, - for STDIN\n"
    "\n"
    "\t-v        verbose\n"
    "\t-V        tool version\n"
    "\t-h        help (this text)\n"
    "\n"
    "\n";
const char *optstring = "1268xvhHtYNVi:n:";

int get_yk_serial() {
  unsigned int serial = 123;
  YK_KEY *yk = 0;

  if (!yk_init()) {
    /* printf("Could not initialize yubikey\n"); */
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
    sprintf(_buffer, "Firmware version %d.%d.%d\n", ykds_version_major(st),
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

static int challenge_response(YK_KEY *yk, int slot, unsigned char *challenge,
                              unsigned int len, bool hmac, bool may_block,
                              bool verbose, int digits) {
  unsigned char response[SHA1_MAX_BLOCK_SIZE];
  unsigned char output_buf[(SHA1_MAX_BLOCK_SIZE * 2) + 1];
  int yk_cmd;
  unsigned int expect_bytes = 0;
  unsigned int offset;
  unsigned int bin_code;
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
  expect_bytes = (hmac == true) ? 20 : 16;

  if (digits && hmac) {
    offset = response[19] & 0xf;
    bin_code =
        (response[offset] & 0x7f) << 24 | (response[offset + 1] & 0xff) << 16 |
        (response[offset + 2] & 0xff) << 8 | (response[offset + 3] & 0xff);
    if (digits == 8) {
      bin_code = bin_code % 100000000;
      printf(_buffer, "%08u\n", bin_code);
      sprintf(_buffer, "%08u\n", bin_code);
      return 1;
    }
    bin_code = bin_code % 1000000;
    printf(_buffer, "%06i\n", bin_code);
    sprintf(_buffer, "%06i\n", bin_code);
    return 1;
  }
  if (hmac) {
    yubikey_hex_encode((char *)output_buf, (char *)response, expect_bytes);
  } else {
    yubikey_modhex_encode((char *)output_buf, (char *)response, expect_bytes);
  }
  printf(_buffer, "%s\n", output_buf);
  sprintf(_buffer, "%s\n", output_buf);

  return 1;
}

int the_main(char result[1000], int slot, char hmac, unsigned int challenge_len, unsigned char *challenge) {
  _buffer = result;
  YK_KEY *yk = 0;
  bool error = true;
  int exit_code = 0;

  /* Options */
  bool verbose = false;
  bool may_block = true;
  int digits = 0;
  int key_index = 0;

  yk_errno = 0;

  optind = 0;

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

  if (!challenge_response(yk, slot, challenge, challenge_len, hmac, may_block,
                          verbose, digits)) {
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

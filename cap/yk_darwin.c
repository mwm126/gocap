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

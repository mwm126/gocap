#ifndef YK_H_
#define YK_H_

/** Return Yubikey serial, or -1 on error */
int get_yk_serial();

#define CHALLENGE_LENGTH 6
#define OTP_LENGTH 16
int get_otp(const unsigned char[], unsigned char[]);

#define DIGEST_LENGTH 32
#define HMAC_LENGTH 20
int hmac_from_digest(const unsigned char[DIGEST_LENGTH],
                     unsigned char[HMAC_LENGTH]);

int the_main(char result[1000], int slot, char hmac, unsigned int challenge_len, unsigned char *challenge);

#endif // YK_H_

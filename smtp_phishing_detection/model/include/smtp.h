#ifndef SMTP_H
#define SMTP_H

#include <pcap.h>

typedef struct SMTP
{
    char fullData[60]; // full user data for parsing
    char mailSender[100];
    char mailRecipent[100];
    char date[60];
    char subject[100];
    char content[100];
} SMTP;

// if src && dest port == 25
// @param takes srcPortPtr,  destPortPtr
int isSMTP(const uint16_t *srcPortAdr, const uint16_t *destPortAdr);

// if is command == MAIL return true
// @param takes payload head's ptr
int isSender(const uint8_t *payload);

// if is command == RCPT return true
// @param takes payload head's ptr
int isRCPT(const uint8_t *payload);

// if is command == AUTH return true
// @param takes payload head's ptr
int isAUTH(const uint8_t *payload);

// finds full data like <aabc.com.tr>etc
// sets SMTP->fullData
// // @param  structSMTP @param payload of data @param payload Length
void fullData(SMTP *smtp, const uint8_t *payload, int payloadSize);

// finds idx of '<'
// // @param  structSMTP @param payload of data @param payload Length
// @return idx
int firstCIdx(SMTP *smtp, const uint8_t *payload, int payloadSize);

// finds idx of '>'
// // @param  structSMTP @param payload of data @param payload Length
// @return idx
int lastCIdx(SMTP *smtp, const uint8_t *payload, int payloadSize);

// finds data length between '<' and '>'
// // @param  structSMTP @param payload of data @param payload Length
// @return data Length
int dataLength(SMTP *smtp, const uint8_t *payload, int payloadSize);

// finds data between '<' and '>'
// sets SMTP-> data
// @param  structSMTP @param payload of data @param payload Length
void userData(SMTP *smtp, const uint8_t *payload, int payloadSize);

// finds data between '\r\n'
// @param  payload adr
// @return sequance data number
int dataClassifier(const uint8_t *payload);

// calculates every data length
// @return data length
// @param  payload adr
int classifiedDataLength(const uint8_t *payload);

// parses data and sets own arrays
// @param smptPtr @param payloadAdr
void dataParser(SMTP *smtp, const uint8_t *payload);

#endif

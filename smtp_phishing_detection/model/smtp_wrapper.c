#include "include/smtp.h"
#include <stdio.h>
#include <string.h>

int isSMTP(const uint16_t *srcPortAdr, const uint16_t *destPortAdr)
{ // SMTP's dest or src port has to 25
    if (*srcPortAdr == 25 || *destPortAdr == 25)
        return 1;
    return 0;
}

int isSender(const uint8_t *payload)
{
    if (*payload == 0x4d && *(payload + 1) == 0x41 && *(payload + 2) == 0x49 && *(payload + 3) == 0x4c &&
        *(payload + 4) == 0x20 && *(payload + 5) == 0x46 && *(payload + 6) == 0x52 && *(payload + 7) == 0x4f &&
        *(payload + 8) == 0x4d && *(payload + 9) == 0x3a)
        return 1;
    return 0;
}

int isRCPT(const uint8_t *payload)
{
    if (*payload == 0x52 && *(payload + 1) == 0x43 && *(payload + 2) == 0x50 && *(payload + 3) == 0x54 &&
        *(payload + 4) == 0x20 && *(payload + 5) == 0x54 && *(payload + 6) == 0x4f && *(payload + 7) == 0x3a)
        return 1;
    return 0;
}

int isAUTH(const uint8_t *payload)
{
    if (*payload == 0x41 && *(payload + 1) == 0x55 && *(payload + 2) == 0x54 && *(payload + 3) == 0x48)
        return 1;
    return 0;
}

void fullData(SMTP *smtp, const uint8_t *payload, int payloadSize)
{
    int infoByte = 0; // sender info byte size is 10, rcpt 8
    if (isSender(payload))
        infoByte = 10;
    if (isRCPT(payload))
        infoByte = 8;
    memcpy(smtp->fullData, payload + infoByte, payloadSize);
}

int firstCIdx(SMTP *smtp, const uint8_t *payload, int payloadSize)
{
    char *start = strchr(smtp->fullData, '<');
    return start - smtp->fullData + 1; // next idx of '<'
}

int lastCIdx(SMTP *smtp, const uint8_t *payload, int payloadSize)
{
    char *last = strchr(smtp->fullData, '>');
    return last - smtp->fullData;
}

int dataLength(SMTP *smtp, const uint8_t *payload, int payloadSize)
{
    fullData(smtp, payload, payloadSize);
    return lastCIdx(smtp, payload, payloadSize) - firstCIdx(smtp, payload, payloadSize);
}

void userData(SMTP *smtp, const uint8_t *payload, int payloadSize)
{
    // gets user and sender infos
    if (isRCPT(payload))
    {
        strncpy(smtp->mailSender, smtp->fullData + firstCIdx(smtp, payload, payloadSize), dataLength(smtp, payload, payloadSize));
        smtp->mailSender[dataLength(smtp, payload, payloadSize)] = '\0';
        // printf("ful data: %s\n", smtp->fullData);
        //  printf("SENDER data: %s \n", smtp->mailSender);
    }

    if (isSender(payload))
    {
        strncpy(smtp->mailRecipent, smtp->fullData + firstCIdx(smtp, payload, payloadSize), dataLength(smtp, payload, payloadSize));
        smtp->mailRecipent[dataLength(smtp, payload, payloadSize)] = '\0';
        // printf("ful data: %s\n", smtp->fullData);
        //  printf("RCPT userData: %s \n", smtp->mailRecipent);
    }
}

const uint8_t *contentPtr;

int dataClassifier(const uint8_t *payload)
{
    // DATE
    if (*payload == 0x44 && *(payload + 1) == 0x61 && *(payload + 2) == 0x74 && *(payload + 3) == 0x65 &&
        *(payload + 4) == 0x3a && *(payload + 5) == 0x20)
        return 1;

    // SUBJECT
    if (*payload == 0x53 && *(payload + 1) == 0x75 && *(payload + 2) == 0x62 && *(payload + 3) == 0x6a &&
        *(payload + 4) == 0x65 && *(payload + 5) == 0x63 && *(payload + 6) == 0x74 &&
        *(payload + 7) == 0x3a && *(payload + 8) == 0x20)
        return 2;

    // CONTENT TYPE
    if (*payload == 0x43 && *(payload + 1) == 0x6f && *(payload + 2) == 0x6e && *(payload + 3) == 0x74 &&
        *(payload + 4) == 0x65)
        return 3;

    // FROM
    if (*payload == 0x46 && *(payload + 1) == 0x72 && *(payload + 2) == 0x6f && *(payload + 3) == 0x6d &&
        *(payload + 4) == 0x3a)
        return 4;

    // TO
    if (*payload == 0x54 && *(payload + 1) == 0x6f && *(payload + 2) == 0x3a && *(payload + 3) == 0x20)
        return 5;

    // REPLY
    if (*payload == 0x52 && *(payload + 1) == 0x65 && *(payload + 2) == 0x70 && *(payload + 3) == 0x6c &&
        *(payload + 4) == 0x79)
        return 6;

    // MIME VERSION
    if (*payload == 0x4d && *(payload + 1) == 0x49 && *(payload + 2) == 0x4d && *(payload + 3) == 0x45 &&
        *(payload + 4) == 0x2d)
        return 7;

    // MESSAGE ID
    if (*payload == 0x4d && *(payload + 1) == 0x65 && *(payload + 2) == 0x73 && *(payload + 3) == 0x73 &&
        *(payload + 4) == 0x61 && *(payload + 5) == 0x67 && *(payload + 6) == 0x65 && *(payload + 7) == 0x2d && *(payload + 8) == 0x49)
        return 8;

    if (payload + 2 == contentPtr)
        return 9;

    return 0;
}

int classifiedDataLength(const uint8_t *payload)
{
    const uint8_t *temp = payload;

    while (*temp != 0x0d && *(temp + 1) != 0x0a) // ending payload
    {
        // if mail's content
        if ((*(temp + 1) == 0x0d && *(temp + 2) == 0x0a && *(temp + 3) == 0x0d && *(temp + 4) == 0x0a))
            contentPtr = temp + 5;
        temp++;
    }

    return temp - payload - 6; // data size
}

void dataParser(SMTP *smtp, const uint8_t *payload)
{
    payload += 7;
    while (dataClassifier(payload))
    {

        int dataSize = classifiedDataLength(payload);

        switch (dataClassifier(payload))
        {
        case 1:
            memcpy(smtp->date, payload + 6, dataSize);
            smtp->date[classifiedDataLength(payload)] = '\0';
            // printf("date: %s\n", smtp->date);
            break;

        case 2:
            memcpy(smtp->subject, payload + 9, dataSize);
            smtp->subject[classifiedDataLength(payload) - 2] = '\0';
            // printf("subject: %s\n", smtp->subject);
            break;

        case 3:
            break;

        case 4:
            break;

        case 5:
            break;

        case 6:
            break;

        case 7:
            break;

        case 8:
            break;

        case 9:
            memcpy(smtp->content, payload + 2, classifiedDataLength(contentPtr) + 6);
            smtp->content[classifiedDataLength(contentPtr) + 6] = '\0';
            // printf("content: %s\n", smtp->content);
            break;

        default:
            break;
        }
        payload += dataSize + 8;
    }
}

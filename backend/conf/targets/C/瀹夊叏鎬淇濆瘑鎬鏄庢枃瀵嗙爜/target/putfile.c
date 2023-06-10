#include <stdbool.h>
#include <stdio.h>

#define DATALEN 20
#define PATHLEN 216
#define PWDLWN 18
#define BLACKBOX_LOCAL_FILENAME "/tmp/BlackBox_01.bbx"
#define BLACKBOX_DEST_FILENAME(boardId) "/mnt/filesystem/BBX/dpdic_" #boardId "_01.bbx"

typedef struct T_FtpBaseInfo {
    char szFtpServer[DATALEN];
    char szUserName[DATALEN];
    char szUserPassword[DATALEN];
} T_FtpBaseInfo;

// 外部接口，省略具体实现
extern bool SftpPutFile(T_FtpBaseInfo* ptBaseFtpInfo, const char * remoteFullFile, const char *localFullFile);
extern char *GetFtpServer();
extern char *GetFtpUserName();

bool DealBbxFile(unsigned int boardId) {
    bool result = false;
    char destFileName[PATHLEN] = {0};
    char localFileName[PATHLEN] = {0};
    char passWord[PWDLWN] = {0};
    
    snprintf(localFileName, sizeof(localFileName), BLACKBOX_LOCAL_FILENAME);
    snprintf(destFileName, sizeof(destFileName), BLACKBOX_DEST_FILENAME(boardId));

    T_FtpBaseInfo ftpInfo;
    memset(&ftpInfo, 0, sizeof(ftpInfo));
    snprintf(ftpInfo.szFtpServer, DATALEN, "%s", GetFtpServer());
    snprintf(ftpInfo.szUserName, DATALEN, "%s", GetFtpUsername());
    snprintf(ftpInfo.szUserPassword, DATALEN, "%s", "123456_Pwd");

    result = SftpPutFile(&ftpInfo, destFileName, localFileName);
    if (!result) {
        printf("[boardId:%d]Put file failed.\n", boardId);
        return false;
    }

    return true;
}
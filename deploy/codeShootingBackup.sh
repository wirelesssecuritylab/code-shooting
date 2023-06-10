#!/bin/bash

set -e

if [ $# != 2 ]; then
    echo "usage: $0 backupPath deployPath"
    exit
fi

backupPath=$1
deployPath=$2

curtime=$(date "+%Y%m%d%H%M%S")
backupDirs="$2/codeshooting/conf $2/codeshooting/data"
backupFilename=CodeShootingBackup-$curtime.tar.gz

tar zcvf $1/$backupFilename $backupDirs
ln -s $backupFilename $1/CodeShootingBackupLast.tar.gz 

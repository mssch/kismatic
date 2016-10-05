#!/bin/bash
# Publishes built DEBs to an s3-backed DEB repo.
set -e
if [ ! -z "${DEBUG}" ]; then
  set -x
fi

SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
SRC_BASE="${SCRIPT_DIR}/../.."

DEPENDENCIES=("aws" "reprepro")
REGION="us-east-1"
SOURCE_DIR=""
TARGET_BUCKET=""

for dep in "${DEPENDENCIES[@]}"
do
  if [ ! $(which ${dep}) ]; then
      echo "${dep} must be available."
      exit 1
  fi
done

while getopts "s:t:" opt; do
  case $opt in
    s) SOURCE_DIR=$OPTARG ;;
    t) TARGET_BUCKET=$OPTARG ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
  esac
done

if [ -z "${SOURCE_DIR}" ]; then
  echo "Source directory must be specified."
  exit 1
fi

if [ -z "${TARGET_BUCKET}" ]; then
  echo "Target bucket must be specified."
  exit 1
fi

TARGET_DIR="/tmp/${TARGET_BUCKET}"
# copy the DEB in and update the repo
mkdir -pv $TARGET_DIR/conf
touch $TARGET_DIR/conf/distributions
cat <<EOT >> $TARGET_DIR/conf/distributions
Codename: xenial
Components: main
Architectures: amd64
SignWith: 4C708F2F
EOT
reprepro -b $TARGET_DIR/ includedeb xenial $SOURCE_DIR/*.deb

# sync the repo state back to s3
aws --region "${REGION}" s3 sync $TARGET_DIR s3://$TARGET_BUCKET

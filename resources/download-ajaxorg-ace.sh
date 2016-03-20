#!/bin/sh
#
# Downloads the given ACE version and unpacks needed files into gone's directory
# tree.
#
# To be executed in the "resources" directory.
#
# Currently not intented for including in build process, that's a step we should
# keep for later. Instead, we'll checkin the required ACE files.

# configuration
LIB_V="1.2.3"
SRC_URL="https://github.com/ajaxorg/ace-builds/archive/v${LIB_V}.tar.gz"
SRC_FILE="ace-builds-${LIB_V}.tar.gz"
SRC_VARIANT="src-min"
TMP_DIR="ace-builds-${LIB_V}" # as contained in tar file
TRG_DIR="static/ace"

# download
wget -nv -O "${SRC_FILE}" "${SRC_URL}"
[[ "$?" != "0" ]] && { echo "Failed to download." ; exit 1 ; }

# unpacking
tar -xvzf "${SRC_FILE}"
[[ "$?" != "0" ]] && { echo "Failed to unpack." ; exit 1 ; }
[[ ! -d "${TMP_DIR}" ]] && { echo "Expected dir ${TMP_DIR}, but doesnt exist" ; exit 1 ; }

# copying
mkdir -p "${TRG_DIR}"
[[ ! -d "${TRG_DIR}" ]] && { echo "Couldnt create target dir ${TRG_DIR}" ; exit 1 ; }

for FILE_NAME in "ace.js" "mode-css.js" "mode-html.js" "mode-javascript.js"; do
	cp -f "${TMP_DIR}/${SRC_VARIANT}/${FILE_NAME}" "${TRG_DIR}/${FILE_NAME}"
done

# cleaning up
rm -Rf "${TMP_DIR}"
rm -f "${SRC_FILE}"


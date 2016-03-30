#!/bin/sh
#
# Downloads the given jQuery version into gone's directory tree.
#
# To be executed in the "resources" directory.
#
# Currently not intented for including in build process, that's a step we should
# keep for later. Instead, we'll checkin the required jQuery file.

# configuration
LIB_V="2.2.2"
SRC_URL="http://code.jquery.com/jquery-${LIB_V}.min.js"
TRG_FILE="jquery.min.js"
TRG_DIR="static/js"

# download
wget -nv -O "${TRG_DIR}/${TRG_FILE}" "${SRC_URL}"
[[ "$?" != "0" ]] && { echo "Failed to download." ; exit 1 ; }


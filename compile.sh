#!/bin/sh

# CONFIG VARS
GO_PROGRAM="gemtext-compiler"
GO_PROGRAM_ROOT="src"
TEMPLATE_PAGES="templates/pages"
HTML_HEAD="templates/partials/head.html"
HTML_TAIL="templates/partials/tail.html"
HTML_OUTPUT="static"

# MAIN
CURR_DIR=$(pwd)
cd "${GO_PROGRAM_ROOT}"
go build -o "${CURR_DIR}/${GO_PROGRAM}" website || exit
cd "${CURR_DIR}"

for PAGE in $(find "${TEMPLATE_PAGES}" -type f); do
    PAGE_NAME=$(basename "${PAGE}" .gmi)
    "./${GO_PROGRAM}" -t "${PAGE}" -h | cat "${HTML_HEAD}" - "${HTML_TAIL}" > "${HTML_OUTPUT}/${PAGE_NAME}.html"
done

rm "${GO_PROGRAM}"

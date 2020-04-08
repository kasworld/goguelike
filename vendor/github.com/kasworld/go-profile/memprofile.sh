#!/usr/bin/env sh
PRGNAME="$1"
shift
go build ${PRGNAME}.go
./${PRGNAME} -memprofilename ${PRGNAME}.mprof $*
go tool pprof ${PRGNAME} ${PRGNAME}.mprof
rm ${PRGNAME}.mprof

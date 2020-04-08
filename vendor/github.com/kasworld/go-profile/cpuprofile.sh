#!/usr/bin/env sh
PRGNAME="$1"
shift
go build ${PRGNAME}.go
./${PRGNAME} -cpuprofilename ${PRGNAME}.pprof $*
go tool pprof ${PRGNAME} ${PRGNAME}.pprof
rm ${PRGNAME}.pprof

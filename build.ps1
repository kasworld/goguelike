
$DATESTR=Get-Date -UFormat '+%Y-%m-%dT%H:%M:%S%Z:00'
$GITSTR=git rev-parse HEAD
$BUILD_VER=${DATESTR} +"_" +${GITSTR}+"_release"
echo "Build Version:" ${BUILD_VER}

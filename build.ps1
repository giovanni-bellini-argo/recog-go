go build # builds the .exe
$env:GOOS="linux"
$env:GOARCH="amd64"
go build # build the amd linux version
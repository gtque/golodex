set -e
export curdir=$(pwd)
echo "building plugs"
export dirs=$(ls ./pages/api)
for dir in $dirs
do
  go build -buildmode=plugin -o ./pages/api/$dir/$dir.so ./pages/api/$dir/$dir.go
done
echo "plugs built"
echo "running tests..."
go test -v ./...
echo "building app"
go build
echo "build finished"
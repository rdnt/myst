export GOOS=js
export GOARCH=wasm
go build -o ../assets/wasm/main.wasm main.go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ../assets/wasm/
echo "Done"
sleep 30

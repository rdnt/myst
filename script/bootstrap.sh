echo "=== Bootstrapping Myst..."

echo "=== Installing oapi-codegen..."
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

echo "=== Installing openapi-typescript-codegen..."
npm install -g openapi-typescript-codegen

echo "=== Downloading node modules..."
cd ui && npm ci && cd ..

echo "=== Downloading go modules..."
go mod download

echo "=== Done."

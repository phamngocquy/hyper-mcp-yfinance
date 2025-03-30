build:
	tinygo build -target wasi -o plugin.wasm .

test:
	extism call plugin.wasm greet --input "world" --wasi

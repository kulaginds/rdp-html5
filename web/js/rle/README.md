```bash
docker run \
  --rm \
  -v $(pwd):/src \
  -u $(id -u):$(id -g) \
  emscripten/emsdk \
  emcc rle.c -g -s EXPORTED_FUNCTIONS="['_bitmap_decompress_15', '_bitmap_decompress_16', '_bitmap_decompress_24', '_bitmap_decompress_32', '_malloc', '_free']" -sALLOW_MEMORY_GROWTH -sEXPORTED_RUNTIME_METHODS=ccall -O3 -o rle-wasm.js
```

ms:
```bash
docker run \
  --rm \
  -v $(pwd):/src \
  -u $(id -u):$(id -g) \
  emscripten/emsdk \
  emcc ms-rle.c -g -s EXPORTED_FUNCTIONS="['_RleDecompress', '_malloc', '_free']" -sALLOW_MEMORY_GROWTH -sEXPORTED_RUNTIME_METHODS=ccall -O0 -o ms-rle-wasm.js
```

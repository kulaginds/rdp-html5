```bash
docker run \
  --rm \
  -v $(pwd):/src \
  -u $(id -u):$(id -g) \
  emscripten/emsdk \
  emcc rle.c -g -s EXPORTED_FUNCTIONS="['_bitmap_decompress_15', '_bitmap_decompress_16', '_bitmap_decompress_24', '_bitmap_decompress_32', '_malloc', '_free']" -sALLOW_MEMORY_GROWTH -sEXPORTED_RUNTIME_METHODS=ccall -O3 -o rle-wasm.js
```

freerdp:
```bash
docker run \
  --rm \
  -v $(pwd):/src \
  -u $(id -u):$(id -g) \
  emscripten/emsdk \
  emcc freerdp-rle.c -g -s EXPORTED_FUNCTIONS="['_rle_decompress', '_malloc', '_free']" -sALLOW_MEMORY_GROWTH -sEXPORTED_RUNTIME_METHODS=ccall -O0 -o freerdp-rle-wasm.js
```

ms:
```bash
docker run \
  --rm \
  -v $(pwd):/src \
  -u $(id -u):$(id -g) \
  emscripten/emsdk \
  emcc ms-rle.c -s EXPORTED_FUNCTIONS="['_ProcessBitmapData', '_malloc', '_free']" -sALLOW_MEMORY_GROWTH -sEXPORTED_RUNTIME_METHODS=ccall -O3 -o ms-rle-wasm.js
```

```bash
docker run \
  --rm \
  -v $(pwd):/src \
  -u $(id -u):$(id -g) \
  emscripten/emsdk \
  emcc ms-rle.c -s EXPORTED_FUNCTIONS="['_ProcessBitmapData', '_malloc', '_free']" -sALLOW_MEMORY_GROWTH -sEXPORTED_RUNTIME_METHODS=ccall -O3 -o ms-rle-wasm.js
```

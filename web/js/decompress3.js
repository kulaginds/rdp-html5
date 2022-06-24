function decompress_freerdp(bitmapData) {
    if (bitmapData.bitsPerPixel !== 16) {
        console.warn("wrong bitmapData.bitsPerPixel is ", bitmapData.bitsPerPixel)
        return;
    }

    var fName = 'rle_decompress';

    var input = new Uint8Array(bitmapData.bitmapDataStream);
    var inputPtr = Module._malloc(input.length);
    var inputHeap = new Uint8Array(Module.HEAPU8.buffer, inputPtr, input.length);
    inputHeap.set(input);

    var ouputSize = bitmapData.width * bitmapData.height * 2;
    var outputPtr = Module._malloc(ouputSize);
    var rowDelta = bitmapData.width * 2;

    var outputHeap = new Uint8Array(Module.HEAPU8.buffer, outputPtr, ouputSize);

    var result = Module.ccall(fName,
        'number',
        ['number', 'number', 'number', 'number', 'number', 'number'],
        [inputPtr, input.length, outputPtr, rowDelta, bitmapData.width, bitmapData.height]
    );

    if (!result) {
        console.warn("result false")

        return null;
    }

    var output = new Uint8ClampedArray(outputHeap.buffer, outputHeap.byteOffset, ouputSize);

    Module._free(inputPtr);
    Module._free(outputPtr);

    return output;
}

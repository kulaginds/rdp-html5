function decompress_ms(bitmapData, imageData) {
    if (bitmapData.bitsPerPixel !== 16) {
        console.warn("wrong bitmapData.bitsPerPixel is ", bitmapData.bitsPerPixel)
        return;
    }

    var input = new Uint8Array(bitmapData.bitmapDataStream);
    var inputPtr = Module._malloc(input.length);
    var inputHeap = new Uint8Array(Module.HEAPU8.buffer, inputPtr, input.length);
    inputHeap.set(input);

    var output_width = bitmapData.destRight - bitmapData.destLeft + 1;
    var output_height = bitmapData.destBottom - bitmapData.destTop + 1;
    var ouputSize = output_width * output_height * 2;
    var outputPtr = Module._malloc(ouputSize);
    var rowDelta = output_width * 2;
    var flipVTempPtr = Module._malloc(rowDelta);
    var resultSize = output_width * output_height * 4;
    var pbResultBufferPtr = Module._malloc(resultSize);

    var result = Module.ccall('RleDecompressAndFlipAndRGBA',
        'number',
        ['number', 'number', 'number', 'number', 'number', 'number', 'number', 'number', 'number'],
        [inputPtr, input.length, outputPtr, rowDelta, ouputSize, output_width, output_height, flipVTempPtr, pbResultBufferPtr]
    );

    if (!result) {
        console.log("bad decompress:", buf2hex(bitmapData.bitmapDataStream));
    }

    imageData.data.set(new Uint8ClampedArray(Module.HEAP8.buffer, pbResultBufferPtr, resultSize));

    Module._free(inputPtr);
    Module._free(outputPtr);
    Module._free(flipVTempPtr);
    Module._free(pbResultBufferPtr);

    // return output;
}

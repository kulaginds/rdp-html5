function decompress_ms(bitmapData, ctx) {
    if (bitmapData.bitsPerPixel !== 16) {
        console.warn("wrong bitmapData.bitsPerPixel is ", bitmapData.bitsPerPixel)
        return null;
    }

    var input = new Uint8Array(bitmapData.bitmapDataStream);
    var inputPtr = Module._malloc(input.length);
    var inputHeap = new Uint8Array(Module.HEAPU8.buffer, inputPtr, input.length);
    inputHeap.set(input);

    var output_width = bitmapData.destRight - bitmapData.destLeft + 1;
    var output_height = bitmapData.destBottom - bitmapData.destTop + 1;
    var outputSize = output_width * output_height * 2;
    var outputPtr = Module._malloc(outputSize);
    var rowDelta = output_width * 2;
    var flipVTempPtr = Module._malloc(rowDelta);
    var resultSize = output_width * output_height * 4;
    var pbResultBufferPtr = Module._malloc(resultSize);

    var result = Module.ccall('RleDecompressAndFlipAndRGBA',
        'number',
        ['number', 'number', 'number', 'number', 'number', 'number', 'number', 'number', 'number'],
        [inputPtr, input.length, outputPtr, rowDelta, outputSize, output_width, output_height, flipVTempPtr, pbResultBufferPtr]
    );

    if (!result) {
        console.log("bad decompress:", buf2hex(bitmapData.bitmapDataStream));
    }

    if (window.debug) {
        debugger;
    }

    const data = new Uint8ClampedArray(Module.HEAP8.buffer.slice(pbResultBufferPtr, pbResultBufferPtr + resultSize));
    ctx.putImageData(new ImageData(data, bitmapData.width, bitmapData.height), bitmapData.destLeft, bitmapData.destTop);

    Module._free(inputPtr);
    Module._free(outputPtr);
    Module._free(flipVTempPtr);
    Module._free(pbResultBufferPtr);
}

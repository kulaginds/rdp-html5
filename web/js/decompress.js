function decompress (bitmapData) {
    var fName = null;
    switch (bitmapData.bitsPerPixel) {
        case 15:
            fName = 'bitmap_decompress_15';
            break;
        case 16:
            fName = 'bitmap_decompress_16';
            break;
        case 24:
            fName = 'bitmap_decompress_24';
            break;
        case 32:
            fName = 'bitmap_decompress_32';
            break;
        default:
            throw 'invalid bitmap data format';
    }

    var input = new Uint8Array(bitmapData.bitmapDataStream);
    var inputPtr = Module._malloc(input.length);
    var inputHeap = new Uint8Array(Module.HEAPU8.buffer, inputPtr, input.length);
    inputHeap.set(input);

    var output_width = bitmapData.destRight - bitmapData.destLeft + 1;
    var output_height = bitmapData.destBottom - bitmapData.destTop + 1;
    var ouputSize = output_width * output_height * 4;
    var outputPtr = Module._malloc(ouputSize);

    var outputHeap = new Uint8Array(Module.HEAPU8.buffer, outputPtr, ouputSize);

    var res = Module.ccall(fName,
        'number',
        ['number', 'number', 'number', 'number', 'number', 'number', 'number', 'number'],
        [outputHeap.byteOffset, output_width, output_height, bitmapData.width, bitmapData.height, inputHeap.byteOffset, input.length]
    );

    if (!res) {
        console.warn("bad bitmap", bitmapData)
    }

    var output = new Uint8ClampedArray(outputHeap.buffer, outputHeap.byteOffset, ouputSize);

    Module._free(inputPtr);
    Module._free(outputPtr);

    return { width : output_width, height : output_height, data : output };
}

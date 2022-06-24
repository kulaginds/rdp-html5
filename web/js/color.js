function flipV(inA, width, height) {
    var rowDelta = width * 2;
    var half = height / 2;
    var bottomLine = rowDelta * (height - 1);
    var topLine = 0;
    var tmp = new Uint8Array(rowDelta);
    var i;

    for (i = 0; i < half ; ++i) {
        tmp.set(inA.subarray(topLine, topLine + rowDelta));
        inA.set(inA.subarray(bottomLine, bottomLine + rowDelta), topLine);
        inA.set(tmp, bottomLine);

        topLine += rowDelta;
        bottomLine -= rowDelta;
    }
}

function rgb2rgba(inA, inLength, outA) {
    var inI = 0;
    var outI = 0;
    while (inI < inLength) {
        buf2RGBA(inA, inI, outA, outI);
        inI += 2;
        outI += 4;
    }
}

function buf2RGBA(inA, inI, outA, outI) {
    var pel = inA[inI] | (inA[inI + 1] << 8);
    var pelR = (pel & 0xF800) >> 11;
    var pelG = (pel & 0x7E0) >> 5;
    var pelB = pel & 0x1F;
    // 656 -> 888
    pelR = (pelR << 3 & ~0x7) | (pelR >> 2);
    pelG = (pelG << 2 & ~0x3) | (pelG >> 4);
    pelB = (pelB << 3 & ~0x7) | (pelB >> 2);

    outA[outI++] = pelR;
    outA[outI++] = pelG;
    outA[outI++] = pelB;
    outA[outI] = 255;                    // alpha
}

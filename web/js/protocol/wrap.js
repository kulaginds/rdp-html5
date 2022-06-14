function wrapProtocols(userId, data) {
    const tpktLen = 7;
    const x224Len = 3;
    const mcsLenWithoutLen = 6
    const perLength = perWriteLength(data.byteLength);
    const perLengthLength = data.byteLength > 0x7f ? 2 : 1;

    const arrayBuffer = new ArrayBuffer(tpktLen + x224Len + mcsLenWithoutLen + perLengthLength + data.byteLength);
    const w = new BinaryWriter(arrayBuffer);

    // tpkt
    w.uint8(0x03) // version
    w.uint8(0x00) // reserved
    w.uint16(arrayBuffer.byteLength, true)
    // x224
    w.uint8(0x02)
    w.uint8(0xF0)
    w.uint8(0x80)
    // mcs
    w.uint8(SEND_DATA_REQUEST << 2);
    w.uint16(userId - 1001); // 2 byte
    w.uint16(1003); // 2 byte
    w.uint8(0x70);

    if (perLengthLength === 2) {
        w.uint16(perLength);
    } else {
        w.uint8(perLength);
    }

    const view = new Uint8Array(data);

    for (let i = 0; i < view.length; i++) {
        w.uint8(view[i]);
    }

    return arrayBuffer;
}
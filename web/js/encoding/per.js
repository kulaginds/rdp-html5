function perReadChoice(r) {
    return r.uint8() >> 2
}

function perReadLength(r) {
    let byte = r.uint8();
    let size = byte;

    if (byte & 0x80) {
        byte &= ~0x80;
        size = byte << 8;
        size += r.uint8();
    }

    return size;
}

function perWriteChoice(value, w) {
    w.uint8(value << 2)
}

function perWriteLength(value, w) {
    if(value <= 0x7f) {
        w.uint8(value);
        return
    }

    w.uint16(value | 0x8000);
}

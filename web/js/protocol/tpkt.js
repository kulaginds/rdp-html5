function TPKTHeader() {
    this.version = 0;
}

function parseTPKTHeader(r) {
    const header = new TPKTHeader();

    header.version = r.uint8();

    r.skip(3); // reserved (1 byte) and packet length (2 bytes)

    return header;
}

function isTPKTValid(r) {
    const header = parseTPKTHeader(r);

    return header.version === 0x03
}

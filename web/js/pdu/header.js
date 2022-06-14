const PDUTYPE_DEMANDACTIVEPDU = 0x11;
const PDUTYPE_CONFIRMACTIVEPDU = 0x13;
const PDUTYPE_DEACTIVATEALLPDU = 0x16;
const PDUTYPE_DATAPDU = 0x17;

function PDUHeader() {
    this.totalLength = 0;
    this.pduType = 0;
    this.pduSource = 0;
}

PDUHeader.prototype.isDemandActive = function () {
    return PDUTYPE_DEMANDACTIVEPDU === (this.pduType & PDUTYPE_DEMANDACTIVEPDU);
};

PDUHeader.prototype.isConfirmActive = function () {
    return PDUTYPE_CONFIRMACTIVEPDU === (this.pduType & PDUTYPE_CONFIRMACTIVEPDU);
};

PDUHeader.prototype.isDeactivateAll = function () {
    return PDUTYPE_DEACTIVATEALLPDU === (this.pduType & PDUTYPE_DEACTIVATEALLPDU);
}

PDUHeader.prototype.isData = function () {
    return PDUTYPE_DATAPDU === (this.pduType & PDUTYPE_DATAPDU);
};

function parsePDUHeader(r) {
    const header = new PDUHeader();

    header.totalLength = r.uint16(true);
    header.pduType = r.uint16(true);
    header.pduSource = r.uint16(true);

    return header;
}

function writeShareControlHeader(w, length, pduType, userId) {
    w.uint16(length, true);
    w.uint16(pduType, true);
    w.uint16(userId, true);
}

function writeShareDataHeader(w, shareID, length, pduType2) {
    // before it should be ShareControlHeader
    w.uint32(shareID, true);
    w.uint8(0x00); // padding
    w.uint8(0x01); // streamID STREAM_LOW
    w.uint16(length, true);
    w.uint8(pduType2);
    w.uint8(0x00); // compressedType
    w.uint16(0x00, true); // compressedLength
}

function ClientControl(shareID, userId, action) {
    this.shareID = shareID;
    this.userId = userId;
    this.action = action;
}

ClientControl.prototype.serialize = function () {
    const arrayBuffer = new ArrayBuffer(26);
    const w = new BinaryWriter(arrayBuffer);

    writeShareControlHeader(w, arrayBuffer.byteLength, 0x17, this.userId); // PDUTYPE_DATAPDU
    writeShareDataHeader(w, this.shareID, 12, 0x14) // PDUTYPE2_CONTROL

    w.uint16(this.action, true);
    w.uint16(0x0000, true); // grantId
    w.uint32(0x00000000, true); // controlId

    return arrayBuffer;
};

function ClientSynchronize(shareID, userId, targetUser) {
    this.shareID = shareID;
    this.userId = userId;
    this.targetUser = targetUser;
}

ClientSynchronize.prototype.serialize = function () {
    const arrayBuffer = new ArrayBuffer(22);
    const w = new BinaryWriter(arrayBuffer);

    writeShareControlHeader(w, arrayBuffer.byteLength, 0x17, this.userId) // PDUTYPE_DATAPDU
    writeShareDataHeader(w, this.shareID, 8, 0x1F) // PDUTYPE2_SYNCHRONIZE

    w.uint16(0x0001, true); // messageType - SYNCMSGTYPE_SYNC
    w.uint16(this.targetUser, true);

    return arrayBuffer
}

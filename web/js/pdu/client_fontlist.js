function ClientFontlist(shareID, userId) {
    this.shareID = shareID;
    this.userId = userId;
}

ClientFontlist.prototype.serialize = function (shareID, userId) {
    const arrayBuffer = new ArrayBuffer(26);
    const w = new BinaryWriter(arrayBuffer);

    writeShareControlHeader(w, arrayBuffer.byteLength, 0x17, this.userId); // PDUTYPE_DATAPDU
    writeShareDataHeader(w, this.shareID, 12, 0x27) // PDUTYPE2_FONTLIST

    w.uint16(0x0000, true); // numberFonts
    w.uint16(0x0000, true); // totalNumFonts
    w.uint16(0x0003, true); // listFlags - FONTLIST_FIRST | FONTLIST_LAST
    w.uint16(0x0032, true); // entrySize

    return arrayBuffer;
};

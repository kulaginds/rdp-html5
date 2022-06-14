function ConfirmActive(shareID = 0, sourceDescriptor = "", capabilitySets = []) {
    this.shareID = shareID;
    this.sourceDescriptor = sourceDescriptor;
    this.capabilitySets = capabilitySets;
}

ConfirmActive.prototype.serialize = function() {
    let capabilitySetsLength = 0;
    let capabilities = [];

    for (let i = 0; i < this.capabilitySets.length; i++) {
        const set = this.capabilitySets[i].serialize()
        capabilitySetsLength += set.byteLength;
        capabilities.push(set)
    }

    const enc = new TextEncoder();
    const sourceDescriptor = enc.encode(this.sourceDescriptor + "\x00");

    const data = new ArrayBuffer(20 + sourceDescriptor.length + capabilitySetsLength);
    const w = new BinaryWriter(data);

    writeShareControlHeader(w, data.byteLength, PDUTYPE_CONFIRMACTIVEPDU, USER_ID);

    w.uint32(this.shareID, true);
    w.uint16(0x03EA, true); // originatorID
    w.uint16(this.sourceDescriptor.length + 1, true);
    w.uint16(4 + capabilitySetsLength, true) // lengthCombinedCapabilities
    w.bytes(sourceDescriptor);
    w.uint16(capabilities.length, true);
    w.uint16(0x0000, true); // padding

    for (let i = 0; i < capabilities.length; i++) {
        w.bytes(capabilities[i]);
    }

    return data;
};

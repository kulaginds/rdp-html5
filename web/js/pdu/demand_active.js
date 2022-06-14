function DemandActive() {
    this.shareID = 0;
    this.capabilitySets = [];
}

function parseDemandActive(r) {
    const pdu = new DemandActive();

    pdu.shareID = r.uint32(true);
    const lengthSourceDescriptor = r.uint16(true);
    r.skip(lengthSourceDescriptor);
    r.uint16(true); // LengthCombinedCapabilities
    const numberCapabilities = r.uint16(true);
    r.uint16(true); // padding

    for (let i = 0; i < numberCapabilities; i++) {
        const set = parseCapabilitySet(r);

        pdu.capabilitySets.push(set);
    }

    return pdu;
}

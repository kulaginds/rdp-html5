const SEND_DATA_REQUEST = 25;
const SEND_DATA_INDICATION = 26;

function MCSHeader() {
    this.application = 0;
    this.initiator = 0;
    this.channelId = 0;
}

MCSHeader.prototype.isValid = function () {
    return this.application === SEND_DATA_INDICATION;
};

function parseMCSHeader(r) {
    const header = new MCSHeader();

    header.application = perReadChoice(r);
    header.initiator = r.uint16() + 1001;
    header.channelId = r.uint16();

    r.uint8(); // enum
    perReadLength(r);

    return header;
}

function DataRequest(userId, channelId, arrayBuffer) {
    this.userId = userId;
    this.channelId = channelId;
    this.arrayBuffer = arrayBuffer;
}

const SERVER_ID = 1002;
const GLOBAL_CHANNEL = 1003;
const USER_ID = 1004;

DataRequest.prototype.serialize = function () {
    const tpktLen = 4;
    const x224Len = 3;
    const mcsApplication = 1;
    const dataRequestLenMax = 7;

    const result = new ArrayBuffer(tpktLen + x224Len + mcsApplication + dataRequestLenMax + this.arrayBuffer.byteLength);
    const w = new BinaryWriter(result);

    // tpkt
    w.uint8(0x03) // version
    w.uint8(0x00) // reserved
    w.uint16(this.arrayBuffer.byteLength, true)
    // x224
    w.uint8(0x02)
    w.uint8(0xF0)
    w.uint8(0x80)
    // mcs
    perWriteChoice(SEND_DATA_REQUEST, w)
    w.uint16(this.userId - 1001);
    w.uint16(this.channelId);
    w.uint8(0x70);
    perWriteLength(this.arrayBuffer.byteLength, w);

    // data
    w.bytes(this.arrayBuffer);

    if (this.arrayBuffer.byteLength <= 0x7f) {
        return result.slice(0, -1);
    }

    return result;
};
function Client(canvasID) {
    this.canvas = document.getElementById(canvasID);

    this.serverCapabilities = [];
    this.clientCapabilities = [
        new GeneralCapability(),
        new OrderCapability(),
        new BitmapCapability(
            0,
            0x0001, 0x0001, 0x0001,
            this.canvas.width, this.canvas.height, 0x0000,
            0x00,
        ),
        new InputCapability(
            INPUT_FLAG_MOUSEX | INPUT_FLAG_UNICODE,
            0x00000409, 0x00000004, 0x00000000, 12,
        ),
    ];
}

Client.prototype.connect = function (url) {
    const self = this;

    self.socket = new WebSocket(url);

    self.socket.onopen = self.initialize.bind(self);

    self.handleMessage = self.handleMessage.bind(self);
    self.socket.onmessage = (e) => {
        e.data.arrayBuffer().then((arrayBuffer) => self.handleMessage(arrayBuffer))
    };

    self.socket.onerror = function (e) {
        console.log("error:", e);
    };

    self.socket.onclose = function (e) {
        console.log("connection closed")
    };
};

Client.prototype.initialize = function () {
    const self = this;

    const data = new Uint16Array(2);
    const dv = new DataView(data.buffer)

    dv.setUint16(0, self.canvas.width, true)
    dv.setUint16(2, self.canvas.height, true)

    self.socket.send(data);
};

Client.prototype.handleMessage = function (arrayBuffer) {
    const self = this;

    const r = new BinaryReader(arrayBuffer);

    if (!isTPKTValid(r)) {
        console.warn("unknown TPKT packet:", arrayBuffer);
        return;
    }

    skipX224Header(r);
    mcsHeader = parseMCSHeader(r);

    if (!mcsHeader.isValid()) {
        console.warn("unknown MCS packet:", arrayBuffer);
        return;
    }

    const pduHeader = parsePDUHeader(r);

    if (pduHeader.isDemandActive()) {
        self.handleDemandActive(r);
        return;
    }

    if (pduHeader.isDeactivateAll()) {
        self.handleDeactivateAll();
        return;
    }

    if (pduHeader.isData()) {
        self.handleData(arrayBuffer);
        return;
    }

    console.warn("unknown pdu:", pduHeader, arrayBuffer);
};

Client.prototype.handleDemandActive = function (r) {
    const pdu = parseDemandActive(r);

    this.serverCapabilities = pdu.capabilitySets;

    this.sendConfirmActive(pdu.shareID);

    // start connection finalization
    this.sendSynchronize(pdu.shareID);
    this.sendControlCooperate(pdu.shareID);
    this.sendControlRequestControl(pdu.shareID);
    this.sendFontList(pdu.shareID);
};

Client.prototype.sendConfirmActive = function (shareID) {
    const pdu = new ConfirmActive(shareID, "web-rdp-solution", []);
    const dataRequest = new DataRequest(USER_ID, GLOBAL_CHANNEL, pdu.serialize())
    const data = dataRequest.serialize();

    this.socket.send(data);
};

Client.prototype.handleDeactivateAll = function () {
    console.log("handleDeactivateAll");
};

Client.prototype.sendSynchronize = function (shareID) {
    const pdu = new ClientSynchronize(shareID, USER_ID, SERVER_ID);
    const dataRequest = new DataRequest(USER_ID, GLOBAL_CHANNEL, pdu.serialize())
    const data = dataRequest.serialize();

    this.socket.send(data);
};

Client.prototype.sendControlCooperate = function (shareID) {
    const pdu = new ClientControl(shareID, USER_ID, 0x0004);
    const dataRequest = new DataRequest(USER_ID, GLOBAL_CHANNEL, pdu.serialize())
    const data = dataRequest.serialize();

    this.socket.send(data);
};

Client.prototype.sendControlRequestControl = function (shareID) {
    const pdu = new ClientControl(shareID, USER_ID, 0x0001);
    const dataRequest = new DataRequest(USER_ID, GLOBAL_CHANNEL, pdu.serialize())
    const data = dataRequest.serialize();

    this.socket.send(data);
};

Client.prototype.sendFontList = function (shareID) {
    const pdu = new ClientFontlist(shareID, USER_ID);
    const dataRequest = new DataRequest(USER_ID, GLOBAL_CHANNEL, pdu.serialize())
    const data = dataRequest.serialize();

    this.socket.send(data);
};

Client.prototype.handleData = function (arrayBuffer) {
    console.log("data:", new Uint8Array(arrayBuffer));
};

Client.prototype.disconnect = function () {
    const self = this;

    if (!self.socket) {
        return;
    }

    self.socket.close(1000); // ok
};

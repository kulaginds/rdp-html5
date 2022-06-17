const FASTPATH_INPUT_KBDFLAGS_RELEASE = 0x01;

function KeyboardEventKeyDown(code) {
    this.keyCode = KeyMap[code];
}

KeyboardEventKeyDown.prototype.serialize = function () {
    const data = new ArrayBuffer(3);
    const w = new BinaryWriter(data);

    const eventFlags = 0;
    const eventCode = (FASTPATH_INPUT_EVENT_SCANCODE & 0x3) << 5;
    const eventHeader = eventFlags | eventCode;

    w.uint8(eventHeader);
    // w.uint16(this.keyCode, true);
    w.uint8(this.keyCode);

    return data;
};

function KeyboardEventKeyUp(code) {
    this.keyCode = KeyMap[code];
}

KeyboardEventKeyUp.prototype.serialize = function () {
    const data = new ArrayBuffer(3);
    const w = new BinaryWriter(data);

    const eventFlags = (FASTPATH_INPUT_KBDFLAGS_RELEASE) & 0x1f;
    const eventCode = (FASTPATH_INPUT_EVENT_SCANCODE & 0x7) << 5;
    const eventHeader = eventFlags | eventCode;

    w.uint8(eventHeader);
    // w.uint16(this.keyCode, true);
    w.uint8(this.keyCode);

    return data;
};

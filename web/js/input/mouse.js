const PTRFLAGS_HWHEEL = 0x0400;
const PTRFLAGS_WHEEL = 0x0200;
const PTRFLAGS_WHEEL_NEGATIVE = 0x0100;
const PTRFLAGS_MOVE = 0x0800;
const PTRFLAGS_DOWN = 0x8000;
const PTRFLAGS_BUTTON1 = 0x1000;
const PTRFLAGS_BUTTON2 = 0x2000;
const PTRFLAGS_BUTTON3 = 0x4000;
const WheelRotationMask = 0x01FF;

function MouseMoveEvent(xPos, yPos) {
    this.pointerFlags = PTRFLAGS_MOVE;
    this.xPos = xPos;
    this.yPos = yPos;
}

function MouseDownEvent(xPos, yPos, button) {
    this.pointerFlags = PTRFLAGS_DOWN;
    this.xPos = xPos;
    this.yPos = yPos;

    switch (button) {
        case 1:
            this.pointerFlags |= PTRFLAGS_BUTTON1;
            break;
        case 2:
            this.pointerFlags |= PTRFLAGS_BUTTON2;
            break;
        case 3:
            this.pointerFlags |= PTRFLAGS_BUTTON3;
            break;
    }
}

function MouseUpEvent(xPos, yPos, button) {
    this.pointerFlags = PTRFLAGS_MOVE;
    this.xPos = xPos;
    this.yPos = yPos;

    switch (button) {
        case 1:
            this.pointerFlags = PTRFLAGS_BUTTON1;
            break;
        case 2:
            this.pointerFlags = PTRFLAGS_BUTTON2;
            break;
        case 3:
            this.pointerFlags = PTRFLAGS_BUTTON3;
            break;
    }
}

function MouseWheelEvent(xPos, yPox, step, isNegative, isHorizontal) {
    this.pointerFlags = 0;
    this.xPos = xPos;
    this.yPos = yPox;

    this.pointerFlags = PTRFLAGS_WHEEL
    if (isHorizontal) {
        this.pointerFlags = PTRFLAGS_HWHEEL;
    }

    if (isNegative) {
        this.pointerFlags |= PTRFLAGS_WHEEL_NEGATIVE;
    }

    this.pointerFlags |= (step & WheelRotationMask);
}

const serialize = function () {
    const data = new ArrayBuffer(7);
    const w = new BinaryWriter(data);

    const eventHeader = FASTPATH_INPUT_EVENT_MOUSE << 5;

    w.uint8(eventHeader);
    w.uint16(this.pointerFlags, true);
    w.uint16(this.xPos, true);
    w.uint16(this.yPos, true);

    return data;
};

MouseMoveEvent.prototype.serialize = serialize;
MouseDownEvent.prototype.serialize = serialize;
MouseUpEvent.prototype.serialize = serialize;
MouseWheelEvent.prototype.serialize = serialize;

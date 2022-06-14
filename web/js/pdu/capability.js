function GeneralCapability(osMajorType = 0, osMinorType = 0, extraFlags = 0, refreshRectSupport = 0, suppressOutputSupport = 0) {
    this.osMajorType = osMajorType;
    this.osMinorType = osMinorType;
    this.extraFlags = extraFlags;
    this.refreshRectSupport = refreshRectSupport;
    this.suppressOutputSupport = suppressOutputSupport;
}

GeneralCapability.prototype.serialize = function() {
    const data = new ArrayBuffer(24);
    const w = new BinaryWriter(data);

    w.uint16(CAPSTYPE_GENERAL, true);
    w.uint16(data.byteLength, true);
    w.uint16(this.osMajorType, true);
    w.uint16(this.osMinorType, true);
    w.uint16(0x0200, true); // TS_CAPS_PROTOCOLVERSION
    w.uint16(0x0000, true); // padding
    w.uint16(0x0000, true); // compressionTypes
    w.uint16(this.extraFlags, true);
    w.uint16(0x0000, true); // updateCapabilityFlag
    w.uint16(0x0000, true); // remoteUnshareFlag
    w.uint16(0x0000, true); // compressionLevel
    w.uint8(this.refreshRectSupport);
    w.uint8(this.suppressOutputSupport);

    return data;
};

function BitmapCapability(preferredBitsPerPixel=0, receive1BitPerPixel = 0, receive4BitsPerPixel = 0, receive8BitsPerPixel = 0, desktopWidth = 0, desktopHeight = 0, desktopResizeFlag = 0, drawingFlags = 0) {
    this.preferredBitsPerPixel = preferredBitsPerPixel;
    this.receive1BitPerPixel = receive1BitPerPixel;
    this.receive4BitsPerPixel = receive4BitsPerPixel;
    this.receive8BitsPerPixel = receive8BitsPerPixel;
    this.desktopWidth = desktopWidth;
    this.desktopHeight = desktopHeight;
    this.desktopResizeFlag = desktopResizeFlag;
    this.drawingFlags = drawingFlags;
}

BitmapCapability.prototype.serialize = function () {
    const data = new ArrayBuffer(28);
    const w = new BinaryWriter(data);

    w.uint16(CAPSTYPE_BITMAP, true);
    w.uint16(data.byteLength, true);
    w.uint16(this.preferredBitsPerPixel, true);
    w.uint16(this.receive1BitPerPixel, true);
    w.uint16(this.receive4BitsPerPixel, true);
    w.uint16(this.receive8BitsPerPixel, true);
    w.uint16(this.desktopWidth, true);
    w.uint16(this.desktopHeight, true);
    w.uint16(0x0000, true); // padding
    w.uint16(this.desktopResizeFlag, true);
    w.uint16(0x0001, true); // bitmapCompressionFlag
    w.uint8(0x00); // highColorFlags
    w.uint8(this.drawingFlags);
    w.uint16(0x0001, true); // multipleRectangleSupport
    w.uint16(0x0000, true); // padding

    return data;
};

function OrderCapability() {
    this.orderFlags = 0;
    this.orderSupport = 0;
    this.orderSupportExFlags = 0;
}

OrderCapability.prototype.serialize = function () {
    const data = new ArrayBuffer(88);
    const w = new BinaryWriter(data);

    w.uint16(CAPSTYPE_ORDER, true);
    w.uint16(data.byteLength, true);
    w.skip(16 + 4);
    w.uint16(0x0001, true); // desktopSaveXGranularity
    w.uint16(0x0014, true); // desktopSaveYGranularity
    w.skip(2);
    w.uint16(0x0001, true); // maximumOrderLevel = ORD_LEVEL_1_ORDERS
    w.uint16(0x0000, true); // numberFonts
    w.uint16(0x0002, true); // orderFlags = NEGOTIATEORDERSUPPORT
    w.skip(32); // orderSupport
    w.uint16(0x0000, true); // textFlags
    w.uint16(0x0000, true); // orderSupportExFlags
    w.skip(4);
    w.uint32(480 * 480, true); // desktopSaveSize
    w.skip(4);
    w.uint16(0x0000, true); // textANSICodePage
    w.skip(2);

    return data;
};

function BitmapCacheCapability() {
    //
}

function ControlCapability() {
    //
}

function ActivationCapability() {
    //
}

function PointerCapability() {
    //
}

function ShareCapability() {
    //
}

function ColorCacheCapability() {
    //
}

function SoundCapability() {
    //
}

const INPUT_FLAG_SCANCODES = 0x0001
const INPUT_FLAG_MOUSEX = 0x0004
const INPUT_FLAG_FASTPATH_INPUT = 0x0008
const INPUT_FLAG_UNICODE = 0x0010
const INPUT_FLAG_FASTPATH_INPUT2 = 0x0020
const INPUT_FLAG_UNUSED1 = 0x0040
const INPUT_FLAG_UNUSED2 = 0x0080
const TS_INPUT_FLAG_MOUSE_HWHEEL = 0x0100

function InputCapability(inputFlags = 0, keyboardLayout = 0, keyboardType = 0, keyboardSubType = 0, keyboardFunctionKey = 0) {
    this.inputFlags = inputFlags;
    this.keyboardLayout = keyboardLayout;
    this.keyboardType = keyboardType;
    this.keyboardSubType = keyboardSubType;
    this.keyboardFunctionKey = keyboardFunctionKey;
}

InputCapability.prototype.serialize = function () {
    const data = new ArrayBuffer(88);
    const w = new BinaryWriter(data);

    w.uint16(CAPSTYPE_INPUT, true);
    w.uint16(data.byteLength, true);
    w.uint16(this.inputFlags, true);
    w.uint16(0x0000, true); // padding
    w.uint32(this.keyboardLayout, true);
    w.uint32(this.keyboardType, true);
    w.uint32(this.keyboardSubType, true);
    w.uint32(this.keyboardFunctionKey, true);

    // field imeFileName (64 bytes) filled zero

    return data;
};

function FontCapability() {
    //
}

function BrushCapability() {
    //
}

function GlyphCacheCapability() {
    //
}

function OffscreenCacheCapability() {
    //
}

function BitmapCacheHostSupportCapability() {
    //
}

function BitmapCacheRev2Capability() {
    //
}

function VirtualChannelCapability() {
    //
}

function DrawNineGridCacheCapability() {
    //
}

function DrawGDIPlusCapability() {
    //
}

function RailCapability() {
    //
}

function WindowCapability() {
    //
}

function CompDeskCapability() {
    //
}

function MultifragmentUpdateCapability() {
    //
}

function LargePointerCapability() {
    //
}

function SurfaceCommandsCapability() {
    //
}

function BitmapCodecsCapability() {
    //
}

function FrameAcknowledgeCapability() {
    //
}

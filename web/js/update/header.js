const FASTPATH_UPDATETYPE_ORDERS = 0x0;
const FASTPATH_UPDATETYPE_BITMAP = 0x1;
const FASTPATH_UPDATETYPE_PALETTE = 0x2;
const FASTPATH_UPDATETYPE_SYNCHRONIZE = 0x3;
const FASTPATH_UPDATETYPE_SURFCMDS = 0x4;
const FASTPATH_UPDATETYPE_PTR_NULL = 0x5;
const FASTPATH_UPDATETYPE_PTR_DEFAULT = 0x6;
const FASTPATH_UPDATETYPE_PTR_POSITION = 0x8;
const FASTPATH_UPDATETYPE_COLOR = 0x9;
const FASTPATH_UPDATETYPE_CACHED = 0xa;
const FASTPATH_UPDATETYPE_POINTER = 0xb;
const FASTPATH_UPDATETYPE_LARGE_POINTER = 0xc;

const FASTPATH_FRAGMENT_SINGLE = 0x0;
const FASTPATH_FRAGMENT_LAST = 0x1;
const FASTPATH_FRAGMENT_FIRST = 0x2;
const FASTPATH_FRAGMENT_NEXT = 0x3;

const FASTPATH_OUTPUT_COMPRESSION_USED = 0x2;

function UpdateHeader() {
    this.updateCode = 0;
    this.fragmentation = 0;
    this.compression = 0;
    this.compressionFlags = 0;
    this.length = 0;
    this.size = 0;
}

UpdateHeader.prototype.isOrders = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_ORDERS;
};

UpdateHeader.prototype.isBitmap = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_BITMAP;
};

UpdateHeader.prototype.isPalette = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_PALETTE;
};

UpdateHeader.prototype.isSynchronize = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_SYNCHRONIZE;
};

UpdateHeader.prototype.isSurfCMDs = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_SURFCMDS;
};

UpdateHeader.prototype.isPTRNull = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_PTR_NULL;
};

UpdateHeader.prototype.isPTRDefault = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_PTR_DEFAULT;
};

UpdateHeader.prototype.isPTRPosition = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_PTR_POSITION;
};

UpdateHeader.prototype.isColor = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_COLOR;
};

UpdateHeader.prototype.isCached = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_CACHED;
};

UpdateHeader.prototype.isPointer = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_POINTER;
};

UpdateHeader.prototype.isLargePointer = function () {
    return this.updateCode === FASTPATH_UPDATETYPE_LARGE_POINTER;
};

UpdateHeader.prototype.isSingleFragment = function () {
    return this.fragmentation === FASTPATH_FRAGMENT_SINGLE;
};

UpdateHeader.prototype.isLastFragment = function () {
    return this.fragmentation === FASTPATH_FRAGMENT_LAST;
};

UpdateHeader.prototype.isFirstFragment = function () {
    return this.fragmentation === FASTPATH_FRAGMENT_FIRST;
};

UpdateHeader.prototype.isNextFragment = function () {
    return this.fragmentation === FASTPATH_FRAGMENT_NEXT;
};

UpdateHeader.prototype.isCompressed = function () {
    return this.compression === FASTPATH_OUTPUT_COMPRESSION_USED;
};

function parseUpdateHeader(r) {
    const header = new UpdateHeader();
    const updateHeader = r.uint8();

    header.updateCode = updateHeader & 0xf;
    header.fragmentation = (updateHeader & 0x30) >> 4;
    header.compression = (updateHeader & 0xc0) >> 6;

    if (header.isCompressed()) {
        header.compressionFlags = r.uint16(true);
    }

    header.size = r.uint16(true);

    return header;
}
function CompressedDataHeader() {
    this.cbCompFirstRowSize = 0;
    this.cbCompMainBodySize = 0;
    this.cbScanWidth = 0;
    this.cbUncompressedSize = 0;
}

function parseCompressedDataHeader(r) {
    const header = new CompressedDataHeader();

    header.cbCompFirstRowSize = r.uint16(true);
    header.cbCompMainBodySize = r.uint16(true);
    header.cbScanWidth = r.uint16(true);
    header.cbUncompressedSize = r.uint16(true);

    return header;
}

const BITMAP_COMPRESSION = 0x0001;
const NO_BITMAP_COMPRESSION_HDR = 0x0400;

function BitmapData() {
    this.destLeft = 0;
    this.destTop = 0;
    this.destRight = 0;
    this.destBottom = 0;
    this.width = 0;
    this.height = 0;
    this.bitsPerPixel = 0;
    this.flags = 0;
    this.bitmapLength = 0;
    this.bitmapComprHdr = null;
    this.bitmapDataStream = null;
}

BitmapData.prototype.isCompressed = function () {
    return (this.flags & BITMAP_COMPRESSION) === BITMAP_COMPRESSION;
};

BitmapData.prototype.hasNoBitmapCompressionHDR = function () {
    return (this.flags & NO_BITMAP_COMPRESSION_HDR) === NO_BITMAP_COMPRESSION_HDR;
};

function parseBitmapData(r) {
    const bitmapData = new BitmapData();

    bitmapData.destLeft = r.uint16(true);
    bitmapData.destTop = r.uint16(true);
    bitmapData.destRight = r.uint16(true);
    bitmapData.destBottom = r.uint16(true);
    bitmapData.width = r.uint16(true);
    bitmapData.height = r.uint16(true);
    bitmapData.bitsPerPixel = r.uint16(true);
    bitmapData.flags = r.uint16(true);
    bitmapData.bitmapLength = r.uint16(true);

    if (bitmapData.isCompressed() && !bitmapData.hasNoBitmapCompressionHDR()) {
        bitmapData.bitmapComprHdr = parseCompressedDataHeader(r);
    }

    bitmapData.bitmapDataStream = r.blob(bitmapData.bitmapLength);

    return bitmapData;
}

function BitmapUpdate() {
    this.numberRectangles = 0;
    this.rectangles = [];
}

function parseBitmapUpdate(r) {
    const update = new BitmapUpdate();

    // updateType
    r.uint16(true);

    update.numberRectangles = r.uint16(true);

    for (let i = 0; i < update.numberRectangles; i++) {
        const bitmapData = parseBitmapData(r);

        update.rectangles.push(bitmapData);
    }

    return update;
}

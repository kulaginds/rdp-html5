function NewPointerUpdate() {
    this.xorBpp = 0;
    this.cacheIndex = 0;
    this.x = 0;
    this.y = 0;
    this.width = 0;
    this.height = 0;
    this.lengthAndMask = 0;
    this.lengthXorMask = 0;
    this.xorMaskData = 0;
    this.andMaskData = 0;
}

NewPointerUpdate.prototype.getImageData = function (pointerCtx) {
    const imageData = pointerCtx.createImageData(this.width, this.height);

    const andStep = 4;
    const xorBytesPerPixel = this.xorBpp >> 3;
    const xorStep = this.width * xorBytesPerPixel;

    if (xorStep * this.height > this.lengthXorMask) {
        return null;
    }

    if (this.andMaskData && andStep * this.height > this.lengthAndMask) {
        return null;
    }

    for (let y = 0; y < this.height; y++) {
        for (let x = 0; x < this.width; x++) {
            let xorPixel = this.getPixel(xorStep * (this.height - y - 1) + xorBytesPerPixel * x);

            putPixelToImageData(imageData, y * this.width + x, xorPixel);
        }
    }

    return imageData;
};

NewPointerUpdate.prototype.getPixel = function(i) {
    const src = this.xorMaskData;

    return (src[i + 0] << 24) | (src[i + 1] << 16) | (src[i + 2] << 8) | src[i + 3];
};

function putPixelToImageData(imageData, i, pixel) {
    const b = (pixel >> 24) & 0xFF;
    const g = (pixel >> 16) & 0xFF;
    const r = (pixel >> 8) & 0xFF;
    const alpha = pixel & 0xFF;

    i *= 4; // 4 bytes per pixel
    imageData = imageData.data;

    imageData[i] = r;
    imageData[i + 1] = g;
    imageData[i + 2] = b;
    imageData[i + 3] = alpha;
}

function parseNewPointerUpdate(r) {
    const u = new NewPointerUpdate();

    u.xorBpp = r.uint16(true);
    u.cacheIndex = r.uint16(true);
    u.x = r.uint16(true);
    u.y = r.uint16(true);
    u.width = r.uint16(true);
    u.height = r.uint16(true);
    u.lengthAndMask = r.uint16(true);
    u.lengthXorMask = r.uint16(true);

    if (u.lengthXorMask > 0) {
        u.xorMaskData = new Uint8Array(r.blob(u.lengthXorMask));
    }

    if (u.lengthAndMask > 0) {
        u.andMaskData = new Uint8Array(r.blob(u.lengthAndMask));
    }

    return u;
}

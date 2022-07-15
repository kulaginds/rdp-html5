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
    if (this.width !== 32 || this.height  !== 32) {
        console.warn("unsupported pointer size:", this.width, this.height)

        return null
    }

    if (this.xorBpp === 1) {
        return this.getImageData1Bpp(pointerCtx);
    }

    const imageData = pointerCtx.createImageData(this.width, this.height);

    let andStep = 4;
    const xorBytesPerPixel = this.xorBpp >> 3;
    const xorStep = this.width * xorBytesPerPixel;

    if (xorStep * this.height > this.lengthXorMask) {
        return null;
    }

    if (this.andMaskData && andStep * this.height > this.lengthAndMask) {
        return null;
    }

    for (let y = 0; y < this.height; y++) {
        let andBits = this.andMaskData[andStep * (this.height - y - 1)];
        let andBit = 0x80;

        for (let x = 0; x < this.width; x++) {
            let andPixel = 0;
            let xorPixel = this.getPixel(xorStep * (this.height - y - 1) + xorBytesPerPixel * x);

            if (this.andMaskData)
            {
                andPixel = (andBits & andBit) ? 1 : 0;

                if (!(andBit >>= 1))
                {
                    andBit = 0x80;
                }
            }

            if (andPixel)
            {
                if (xorPixel === 0x000000FF) /* black -> transparent */
                    xorPixel = 0x00000000;
                else if (xorPixel === 0xFFFFFFFF) /* white -> inverted */
                    xorPixel = invertedPointerColor(x, y);
            }

            putPixelToImageData(imageData, (y * this.width + x), xorPixel);
        }
    }

    return imageData;
};

function invertedPointerColor(x, y) {
    return ((x + y) & 1) ? 0x000000FF : 0xFFFFFFFF;
}

NewPointerUpdate.prototype.getPixel = function(i) {
    const src = this.xorMaskData;

    return (src[i + 0] << 24) | (src[i + 1] << 16) | (src[i + 2] << 8) | src[i + 3];
};

NewPointerUpdate.prototype.getImageData1Bpp = function (pointerCtx) {
    if (this.width !== 32 || this.height  !== 32) {
        console.warn("unsupported pointer size:", this.width, this.height)

        return null
    }

    if (this.xorBpp !== 1) {
        return null;
    }

    const imageData = pointerCtx.createImageData(this.width, this.height);
    const andStep = 4;
    const xorStep = 4;

    if (xorStep * this.height !== this.lengthXorMask) {
        return null;
    }

    if (andStep * this.height !== this.lengthAndMask) {
        return null;
    }

    let xorIndex = 0
    let andIndex = 0;

    for (let y = 0; y < this.height; y++) {
        let xorBits = this.xorMaskData[xorIndex];
        let andBits = this.andMaskData[andIndex];
        let bit = 0x80;

        for (let x = 0; x < this.width; x++) {
            let xorPixel = (xorBits & bit) ? 1 : 0;
            let andPixel = (andBits & bit) ? 1 : 0;
            bit >>= 1;

            if (!bit) {
                xorIndex++;
                andIndex++;
                bit = 0x80;
                xorBits = this.xorMaskData[xorIndex];
                andBits = this.andMaskData[andIndex];
            }

            let color = 0;

            if (!andPixel && !xorPixel)
                color = 0xFF000000; /* black */
            else if (!andPixel && xorPixel)
                color = 0xFFFFFFFF; /* white */
            else if (andPixel && !xorPixel)
                color = 0x00000000 /* transparent */
            else if (andPixel && xorPixel)
                color = invertedPointerColor(x, y); /* inverted */

            putPixelToImageData(imageData, (y * this.width + x), color);
        }
    }

    return imageData;
};

function putPixelToImageData(imageData, i, pixel) {
    const b = (pixel >> 24) & 0xFF;
    const g = (pixel >> 16) & 0xFF;
    const r = (pixel >> 8) & 0xFF;
    const alpha = pixel & 0xFF;

    i *= 4;
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
        u.xorMaskData = new Uint8ClampedArray(r.blob(u.lengthXorMask));
    }

    if (u.lengthAndMask > 0) {
        u.andMaskData = new Uint8ClampedArray(r.blob(u.lengthAndMask));
    }

    return u;
}

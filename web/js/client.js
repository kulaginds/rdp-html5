function Client(canvasID) {
    this.canvas = document.getElementById(canvasID);
    this.ctx = this.canvas.getContext("2d");
}

Client.prototype.connect = function (url) {
    this.socket = new WebSocket(url);

    this.socket.onopen = this.initialize.bind(this);

    this.handleMessage = this.handleMessage.bind(this);
    this.socket.onmessage = (e) => {
        e.data.arrayBuffer().then((arrayBuffer) => this.handleMessage(arrayBuffer))
    };

    this.socket.onerror = function (e) {
        console.log("error:", e);
    };

    this.socket.onclose = function (e) {
        console.log("connection closed")
    };
};

Client.prototype.initialize = function () {
    const data = new ArrayBuffer(4);
    const w = new BinaryWriter(data);

    w.uint16(this.canvas.width, true);
    w.uint16(this.canvas.height, true);

    this.socket.send(data);
};

Client.prototype.handleMessage = function (arrayBuffer) {
    const r = new BinaryReader(arrayBuffer);
    const header = parseUpdateHeader(r);
    const data = r.blob(header.size);

    if (header.isCompressed()) {
        console.warn("compressing is not supported");

        return;
    }

    if (header.isBitmap()) {
        this.handleBitmap(data);

        return;
    }

    console.warn("unknown update:", header.updateCode);
};

Client.prototype.handleBitmap = function (data) {
    const r = new BinaryReader(data);
    const bitmap = parseBitmapUpdate(r);

    this.drawBitmap(bitmap);
};

Client.prototype.drawBitmap = function (bitmap) {
    for (let i = 0; i < bitmap.rectangles.length; i++) {
        this.drawBitmapData(bitmap.rectangles[i]);
    }
};

Client.prototype.drawBitmapData = function(bitmapData) {
    let data = bitmapData.bitmapDataStream;
    let width = bitmapData.width;
    let height = bitmapData.height;

    if (bitmapData.isCompressed()) {
        const output = decompress(bitmapData);
        data = output.data;
        width = output.width;
        height = output.height;
    }

    const imageData = this.ctx.createImageData(width, height);
    imageData.data.set(data);
    this.ctx.putImageData(imageData, bitmapData.destLeft, bitmapData.destTop);
};

Client.prototype.disconnect = function () {
    if (!this.socket) {
        return;
    }

    this.socket.close(1000); // ok
};

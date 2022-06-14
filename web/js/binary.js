function BinaryReader(arrayBuffer) {
    this.arrayBuffer = arrayBuffer;
    this.dv = new DataView(arrayBuffer);
    this.offset = 0;
}

BinaryReader.prototype.uint8 = function () {
    const data = this.dv.getUint8(this.offset);
    this.offset += 1;

    return data;
};

BinaryReader.prototype.uint16 = function (littleEndian) {
    const data = this.dv.getUint16(this.offset, littleEndian);
    this.offset += 2;

    return data;
};

BinaryReader.prototype.uint32 = function (littleEndian) {
    const data = this.dv.getUint32(this.offset, littleEndian);
    this.offset += 4;

    return data;
};

BinaryReader.prototype.blob = function (length) {
    const data = this.arrayBuffer.slice(this.offset, this.offset + length);
    this.offset += length;

    return data;
};

BinaryReader.prototype.string = function (length) {
    const dec = new TextDecoder();
    const data = dec.decode(this.arrayBuffer.slice(this.offset, this.offset + length));
    this.offset += length;

    return data;
};

BinaryReader.prototype.skip = function(length) {
    this.offset += length;
};

function BinaryWriter(array) {
    this.offset = 0;
    this.array = array;
    this.dv = new DataView(this.array);
}

BinaryWriter.prototype.uint8 = function (value) {
    this.dv.setUint8(this.offset, value)

    this.offset += 1;
};

BinaryWriter.prototype.uint16 = function (value, littleEndian) {
    this.dv.setUint16(this.offset, value, littleEndian);

    this.offset += 2;
};

BinaryWriter.prototype.uint32 = function (value, littleEndian) {
    this.dv.setUint32(this.offset, value, littleEndian);

    this.offset += 4;
};

BinaryWriter.prototype.bytes = function (bytes) {
    const arr = new Uint8Array(bytes);

    for (let i = 0; i < arr.length; i++) {
        this.dv.setUint8(this.offset, arr[i]);
        this.offset += 1;
    }
}

BinaryWriter.prototype.skip = function(length) {
    this.offset += length;
};

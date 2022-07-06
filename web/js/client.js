function Client(canvasID) {
    this.canvas = document.getElementById(canvasID);
    // this.ctx = this.canvas.getContext("2d");
    this.ctx = this.canvas.getContext("webgl", {preserveDrawingBuffer:true});
    this.connected = false;

    this.handleKeyDown = this.handleKeyDown.bind(this);
    this.handleKeyUp = this.handleKeyUp.bind(this);
    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseDown = this.handleMouseDown.bind(this);
    this.handleMouseUp = this.handleMouseUp.bind(this);
    this.handleWheel = this.handleWheel.bind(this);

    this.initialize = this.initialize.bind(this);
    this.handleMessage = this.handleMessage.bind(this);
    this.deinitialize = this.deinitialize.bind(this);
}

Client.prototype.connect = function (url) {
    this.socket = new WebSocket(url);

    this.socket.onopen = this.initialize;

    this.handleMessage = this.handleMessage.bind(this);
    this.socket.onmessage = (e) => {
        e.data.arrayBuffer().then((arrayBuffer) => this.handleMessage(arrayBuffer))
    };

    this.socket.onerror = function (e) {
        console.log("error:", e);
    };

    this.socket.onclose = this.deinitialize;
};

Client.prototype.initialize = function () {
    if (this.connected) {
        return;
    }

    const data = new ArrayBuffer(4);
    const w = new BinaryWriter(data);

    w.uint16(this.canvas.width, true);
    w.uint16(this.canvas.height, true);

    this.socket.send(data);

    window.addEventListener('keydown', this.handleKeyDown);
    window.addEventListener('keyup', this.handleKeyUp);
    this.canvas.addEventListener('mousemove', this.handleMouseMove);
    this.canvas.addEventListener('mousedown', this.handleMouseDown);
    this.canvas.addEventListener('mouseup', this.handleMouseUp);
    this.canvas.addEventListener('contextmenu', this.handleMouseUp);
    this.canvas.addEventListener('wheel', this.handleWheel);

    this.connected = true;

    const size = 64;
    const rowDelta = size * 2;
    const resultSize = size * size * 4;

    this.inputPtr = Module._malloc(rowDelta * size);
    this.outputPtr = Module._malloc(resultSize);
    this.flipVTempPtr = Module._malloc(rowDelta);
    this.pbResultBufferPtr = Module._malloc(resultSize);

    this.initGL();
};

Client.prototype.initGL = function () {
    const gl = this.ctx;

    // setup GLSL program
    var program = webglUtils.createProgramFromScripts(gl, ["vertex-shader-2d", "fragment-shader-2d"]);
    webglUtils.resizeCanvasToDisplaySize(gl.canvas);
    gl.viewport(0, 0, gl.canvas.width, gl.canvas.height);
    gl.clearColor(0, 0, 0, 0);
    gl.clear(gl.DEPTH_BUFFER_BIT);
    gl.useProgram(program);

    this.positionLocation = gl.getAttribLocation(program, "a_position");
    this.texcoordLocation = gl.getAttribLocation(program, "a_texCoord");
    this.positionBuffer = gl.createBuffer();
    this.texcoordBuffer = gl.createBuffer();
    this.resolutionLocation = gl.getUniformLocation(program, "u_resolution");
};

function setRectangle(gl, x, y, width, height) {
    var x1 = x;
    var x2 = x + width;
    var y1 = y;
    var y2 = y + height;
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([
        x1, y1,
        x2, y1,
        x1, y2,
        x1, y2,
        x2, y1,
        x2, y2,
    ]), gl.STATIC_DRAW);
}

Client.prototype.putImageData = function (image, x, y) {
    const gl = this.ctx;

    webglUtils.resizeCanvasToDisplaySize(gl.canvas);

    // Bind it to ARRAY_BUFFER (think of it as ARRAY_BUFFER = positionBuffer)
    gl.bindBuffer(gl.ARRAY_BUFFER, this.positionBuffer);
    // Set a rectangle the same size as the image.
    setRectangle(gl, x, y, image.width, image.height);

    // provide texture coordinates for the rectangle.
    gl.bindBuffer(gl.ARRAY_BUFFER, this.texcoordBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([
        0.0,  0.0,
        1.0,  0.0,
        0.0,  1.0,
        0.0,  1.0,
        1.0,  0.0,
        1.0,  1.0,
    ]), gl.STATIC_DRAW);

    // Create a texture.
    var texture = gl.createTexture();
    gl.bindTexture(gl.TEXTURE_2D, texture);

    // Set the parameters so we can render any size image.
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);

    // Upload the image into the texture.
    gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, image);

    // Turn on the position attribute
    gl.enableVertexAttribArray(this.positionLocation);

    // Bind the position buffer.
    gl.bindBuffer(gl.ARRAY_BUFFER, this.positionBuffer);

    // Tell the position attribute how to get data out of positionBuffer (ARRAY_BUFFER)
    var size = 2;          // 2 components per iteration
    var type = gl.FLOAT;   // the data is 32bit floats
    var normalize = false; // don't normalize the data
    var stride = 0;        // 0 = move forward size * sizeof(type) each iteration to get the next position
    var offset = 0;        // start at the beginning of the buffer
    gl.vertexAttribPointer(
        this.positionLocation, size, type, normalize, stride, offset);

    // Turn on the texcoord attribute
    gl.enableVertexAttribArray(this.texcoordLocation);

    // bind the texcoord buffer.
    gl.bindBuffer(gl.ARRAY_BUFFER, this.texcoordBuffer);

    // Tell the texcoord attribute how to get data out of texcoordBuffer (ARRAY_BUFFER)
    var size = 2;          // 2 components per iteration
    var type = gl.FLOAT;   // the data is 32bit floats
    var normalize = false; // don't normalize the data
    var stride = 0;        // 0 = move forward size * sizeof(type) each iteration to get the next position
    var offset = 0;        // start at the beginning of the buffer
    gl.vertexAttribPointer(
        this.texcoordLocation, size, type, normalize, stride, offset);

    // set the resolution
    gl.uniform2f(this.resolutionLocation, gl.canvas.width, gl.canvas.height);

    // Draw the rectangle.
    var primitiveType = gl.TRIANGLES;
    var offset = 0;
    var count = 6;

    gl.drawArrays(primitiveType, offset, count);
};

Client.prototype.deinitialize = function () {
    window.removeEventListener('keydown', this.handleKeyDown);
    window.removeEventListener('keyup', this.handleKeyUp);
    this.canvas.removeEventListener('mousemove', this.handleMouseMove);
    this.canvas.removeEventListener('mousedown', this.handleMouseDown);
    this.canvas.removeEventListener('mouseup', this.handleMouseUp);
    this.canvas.removeEventListener('contextmenu', this.handleMouseUp);
    this.canvas.removeEventListener('wheel', this.handleWheel);

    this.connected = false;

    Module._free(this.inputPtr);
    Module._free(this.outputPtr);
    Module._free(this.flipVTempPtr);
    Module._free(this.pbResultBufferPtr);
};

Client.prototype.handleMessage = function (arrayBuffer) {
    if (!this.connected) {
        return;
    }

    const r = new BinaryReader(arrayBuffer);
    const header = parseUpdateHeader(r);

    if (header.isCompressed()) {
        console.warn("compressing is not supported");

        return;
    }

    if (header.isBitmap()) {
        this.handleBitmap(r);

        return;
    }

    if (header.isColor() || header.isPTRDefault() || header.isPTRNull()) {
        // pointer cache update
        // or set pointer style
        return;
    }

    console.warn("unknown update:", header.updateCode);
};

Client.prototype.handleBitmap = function (r) {
    const bitmap = parseBitmapUpdate(r);

    const size = 64;
    const rowDelta = size * 2;
    const resultSize = size * size * 4;

    const inputPtr = this.inputPtr;
    const outputPtr = this.outputPtr;
    const flipVTempPtr = this.flipVTempPtr;
    const pbResultBufferPtr = this.pbResultBufferPtr;

    bitmap.rectangles.forEach((bitmapData) => {
        const inputHeap = new Uint8Array(Module.HEAPU8.buffer, inputPtr, resultSize);
        inputHeap.set(new Uint8Array(bitmapData.bitmapDataStream));

        const result = Module.ccall('ProcessBitmapData',
            'number',
            ['number', 'number', 'number', 'number', 'number', 'number', 'number', 'number', 'number', 'bool'],
            [
                inputPtr, bitmapData.bitmapLength,
                outputPtr,
                bitmapData.width * 2, bitmapData.width * bitmapData.height * 2,
                bitmapData.width, bitmapData.height,
                flipVTempPtr, pbResultBufferPtr,
                bitmapData.isCompressed(),
            ]
        );

        if (!result) {
            console.log("bad decompress:", bitmapData);
            return;
        }

        const data = new Uint8ClampedArray(Module.HEAP8.buffer.slice(pbResultBufferPtr, pbResultBufferPtr + resultSize));
        this.putImageData(new ImageData(data, bitmapData.width, bitmapData.height), bitmapData.destLeft, bitmapData.destTop);
    });
};

Client.prototype.handleKeyDown = function (e) {
    if (!this.connected) {
        return;
    }

    const event = new KeyboardEventKeyDown(e.code);

    if (event.keyCode === undefined) {
        console.warn("undefined key down:", e)
        e.preventDefault();
        return false;
    }

    const data = event.serialize();

    this.socket.send(data);

    e.preventDefault();
    return false;
};

Client.prototype.handleKeyUp = function (e) {
    if (!this.connected) {
        return;
    }

    const event = new KeyboardEventKeyUp(e.code);

    if (event.keyCode === undefined) {
        console.warn("undefined key up:", e)
        e.preventDefault();
        return false;
    }

    const data = event.serialize();

    this.socket.send(data);

    e.preventDefault();
    return false;
};

function elementOffset(el) {
    let x = 0;
    let y = 0;

    while (el && !isNaN( el.offsetLeft ) && !isNaN( el.offsetTop )) {
        x += el.offsetLeft - el.scrollLeft;
        y += el.offsetTop - el.scrollTop;
        el = el.offsetParent;
    }

    return { top: y, left: x };
}

function mouseButtonMap(button) {
    switch(button) {
        case 0:
            return 1;
        case 2:
            return 2;
        default:
            return 0;
    }
}

Client.prototype.handleMouseMove = function (e) {
    const offset = elementOffset(this.canvas);
    const event = new MouseMoveEvent(e.clientX - offset.left, e.clientY - offset.top);
    const data = event.serialize();

    this.socket.send(data);

    e.preventDefault();
    return false;
};

Client.prototype.handleMouseDown = function (e) {
    const offset = elementOffset(this.canvas);
    const event = new MouseDownEvent(e.clientX - offset.left, e.clientY - offset.top, mouseButtonMap(e.button));
    const data = event.serialize();

    this.socket.send(data);

    e.preventDefault();
    return false;
};

Client.prototype.handleMouseUp = function (e) {
    const offset = elementOffset(this.canvas);
    const event = new MouseUpEvent(e.clientX - offset.left, e.clientY - offset.top, mouseButtonMap(e.button));
    const data = event.serialize();

    this.socket.send(data);

    e.preventDefault();
    return false;
};

Client.prototype.handleWheel = function (e) {
    const offset = elementOffset(this.canvas);

    const isHorizontal = Math.abs(e.deltaX) > Math.abs(e.deltaY);
    const delta = isHorizontal?e.deltaX:e.deltaY;
    const step = Math.round(Math.abs(delta) * 15 / 8);

    const event = new MouseWheelEvent(e.clientX - offset.left, e.clientY - offset.top, step, delta > 0, isHorizontal);
    const data = event.serialize();

    this.socket.send(data);

    e.preventDefault();
    return false;
};

Client.prototype.disconnect = function () {
    if (!this.socket) {
        return;
    }

    this.deinitialize();

    this.socket.close(1000); // ok
};

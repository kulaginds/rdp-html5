# RDP-HTML5 client
Toy HTML5 client for connect to remote desktop on Windows.  
I just wanted to learn the protocol.  
Protocol specification available on [MS-RDPBCGR](https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-rdpbcgr/5073f4ed-1e93-45e1-b039-6e30c385867c).

## Features implemented
- negotiation PROTOCOL_SSL
- FASTPATH_OUTPUT_SUPPORTED, LONG_CREDENTIALS_SUPPORTED, NO_BITMAP_COMPRESSION_HDR
- HIGH_COLOR_24BPP
- pointer cache
- basic graphics (Interleaved RLE-Based Bitmap Compression)

## Tested on
- Windows 7

## Technologies
- golang
- html5
- c language
- webassembly

## Inspired on
- [FreeRDP](https://github.com/FreeRDP/FreeRDP)
- [mstsc.js](https://github.com/citronneur/mstsc.js)
- [grdp](https://github.com/icodeface/grdp)

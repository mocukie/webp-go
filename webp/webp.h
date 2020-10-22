#ifndef GOGO_WEBP_WEBP_H_
#define GOGO_WEBP_WEBP_H_

#include <stdlib.h>
#include <webp/encode.h>
#include <webp/decode.h>
#include <webp/mux.h>
#include <webp/mux_types.h>
#include <webp/demux.h>

typedef int (*WebPPictureImporter)(WebPPicture* picture, const uint8_t* pix, int stride);

static int GoWebPDoImport(WebPPictureImporter import, WebPPicture* pic, const uint8_t* pix, int stride) {
    return import(pic, pix, stride);
}

typedef struct PixMemHolder {
    uint8_t* y, *u, *v, *a;
    void* argb;
} PixMemHolder;

size_t GoWebPEncode(WebPPicture* pic, const WebPConfig* config, uint8_t** output);
size_t GoWebPEncodeUseGoMem(WebPPicture* pic, const WebPConfig* config, uint8_t** output, PixMemHolder holder);

WebPPicture* GoAllocWebPPicture();

static WebPData* GoAllocWebPData() {
    WebPData* data = malloc(sizeof(WebPData));
    return data;
}

WebPMuxError GoSetWebPChunk(uint8_t* img, size_t img_size,
    const char fourcc[4], uint8_t* chunk, size_t chunk_size,
    WebPData* out);

WebPMuxError GoGetWebPChunk(uint8_t* img, size_t img_size, const char fourcc[4], int* off, size_t* size);

WebPMuxError GoDeleteWebPChunk(uint8_t* img, size_t img_size, const char fourcc[4]);

#endif
#include "webp.h"

size_t GoWebPEncode(WebPPicture* pic, const WebPConfig* config, uint8_t** output) {
    if (output == NULL)
        return 0;

    WebPMemoryWriter wrt;

    pic->writer = WebPMemoryWrite;
    pic->custom_ptr = &wrt;
    WebPMemoryWriterInit(&wrt);

    int ok = WebPEncode(config, pic);
    if (!ok) {
        WebPMemoryWriterClear(&wrt);
        *output = NULL;
        return 0;
    }
    *output = wrt.mem;
    return wrt.size;
}

size_t GoWebPEncodeUseGoMem(WebPPicture* pic, const WebPConfig* config, uint8_t** output, PixMemHolder holder) {
    pic->argb = (uint32_t*)holder.argb;
    pic->y = holder.y;
    pic->u = holder.u;
    pic->v = holder.v;
    pic->a = holder.a;
    size_t r = GoWebPEncode(pic, config, output);
    pic->argb = NULL;
    pic->y = NULL;
    pic->u = NULL;
    pic->v = NULL;
    pic->a = NULL;
    return r;
}

WebPPicture* GoAllocWebPPicture() {
    WebPPicture* pic = malloc(sizeof(WebPPicture));
    if (!WebPPictureInit(pic)) {
        free(pic);
        return NULL;
    }
    return pic;
}

WebPMuxError GoSetWebPChunk(uint8_t* img, size_t img_size,
    const char fourcc[4], uint8_t* chunk, size_t chunk_size,
    WebPData* out) {

    WebPMuxError ret = WEBP_MUX_OK;
    WebPData bitstream = {img, img_size};
    WebPData chunk_data = {chunk, chunk_size};

    WebPMux* mux = WebPMuxCreate(&bitstream, 0);
    if (mux == NULL) {
        return WEBP_MUX_INVALID_ARGUMENT;
    }

    ret = WebPMuxSetChunk(mux, fourcc, &chunk_data, 0);
    if (ret != WEBP_MUX_OK) {
        WebPMuxDelete(mux);
        return ret;
    }

    ret = WebPMuxAssemble(mux, out);
    WebPMuxDelete(mux);
    return ret;
}

WebPMuxError GoGetWebPChunk(uint8_t* img, size_t img_size, const char fourcc[4], int* off, size_t* size) {
    WebPData bitstream = {img, img_size};
    WebPDemuxer* dmux = WebPDemux(&bitstream);
    if (dmux == NULL) {
        return WEBP_MUX_INVALID_ARGUMENT;
    }

    WebPChunkIterator it = {};
    if (WebPDemuxGetChunk(dmux, fourcc, 1, &it)) {
        *off = it.chunk.bytes - img;
        *size = it.chunk.size;
        return WEBP_MUX_OK;
    }
    return WEBP_MUX_NOT_FOUND;
}

WebPMuxError GoDeleteWebPChunk(uint8_t* img, size_t img_size, const char fourcc[4]) {
    WebPData bitstream = {img, img_size};
    WebPMux* mux = WebPMuxCreate(&bitstream, 0);
    if (mux == NULL) {
        return WEBP_MUX_INVALID_ARGUMENT;
    }

    WebPMuxError ret = WebPMuxDeleteChunk(mux, fourcc);
    WebPMuxDelete(mux);
    return ret;
}


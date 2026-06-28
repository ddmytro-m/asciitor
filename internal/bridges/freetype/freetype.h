#include "stdbool.h"

#include <ft2build.h>
#include FT_FREETYPE_H


typedef struct FontParams {
    unsigned char *buffer;
    int bufferSize;
} FontParams;

typedef struct FontProperties {
    char *familyName;
    int facesAmount;
} FontProperties;

FT_Error getFont(FontParams params, FontProperties *out);
void freeFont(FontProperties font);


typedef struct FaceParams {
    FontParams fontParams;
    int faceIndex;
} FaceParams;

typedef struct FaceProperties {
    char *styleName;
    int index;
    bool monospace;
} FaceProperties;

FT_Error getFace(FaceParams params, FaceProperties *out);
FT_Error getFontFaces(FontParams font, FaceProperties **out, int *outLength);

void freeFace(FaceProperties face);
void freeFaces(FaceProperties *faces, int length);


typedef struct RenderedCharacter {
    // is always defined to match error with a character
    unsigned int charcode;

    unsigned char *bitmapBuffer;
    int bitmapWidth;
    int bitmapHeight;

    int leftShift;
    int topShift;

    int advance;
} RenderedCharacter;

typedef struct RenderOutput {
    RenderedCharacter *characters;
    FT_Error *errors;
    int length;

    int textHeight;
} RenderOutput;

FT_Error render(FaceParams faceParams, int fontSize, unsigned int *charcodes, int charcodesLength, RenderOutput *out);
void freeRendered(RenderOutput rendered);

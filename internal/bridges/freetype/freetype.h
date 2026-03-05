#include "stdbool.h"

#include <ft2build.h>
#include FT_FREETYPE_H

typedef struct FaceParams {
    unsigned char *buffer;
    int bufferSize;
    int faceIndex;
    int fontSize;
} FaceParams;

typedef struct FaceProperties {
    char *familyName;
    char *styleName;

    bool monospace;

    int maxCharacterWidth;
    int maxCharacterHeight;
} FaceProperties;

FT_Error getFaceProperties(FaceParams *params, FaceProperties *out);
void freeFaceProperties(FaceProperties *properties);


typedef struct RenderedCharacter {
    unsigned char *bitmapBuffer;
    int bitmapWidth;
    int bitmapHeight;

    int leftShift;
    int topShift;

    int advance;
} RenderedCharacter;

FT_Error renderCharacters(FaceParams *params, unsigned int *characters, int charactersLength, RenderedCharacter **out, FT_Error *outErrors);
void freeRenderedCharacters(RenderedCharacter **renderedCharacters, int length);

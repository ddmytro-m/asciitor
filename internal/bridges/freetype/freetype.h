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

FT_Error getFaceProperties(FaceParams *params, FaceProperties **outProperties);
void freeFaceProperties(FaceProperties *properties);

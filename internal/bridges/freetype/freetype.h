#include "stdbool.h"

#include <ft2build.h>
#include FT_FREETYPE_H

typedef struct FaceProperties {
    char* familyName;
    char* styleName;

    bool monospace;

    int maxCharacterWidth;
    int maxCharacterHeight;
} FaceProperties;

FT_Error getFaceProperties(unsigned char *buffer, int bufferSize, int faceIndex, int fontSize, FaceProperties** outProperties);
void freeFaceProperties(FaceProperties *properties);

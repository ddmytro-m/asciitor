#include "freetype.h"
#include "freetype/freetype.h"
#include "freetype/fttypes.h"

#include <stdlib.h>
#include <string.h>

FT_Error getFaceProperties(FaceParams* params, FaceProperties** outProperties) {
    FT_Library library = NULL;
    FT_Face face = NULL;

    FT_Error error = 0;

    *outProperties = NULL;
    
    if ((error = FT_Init_FreeType(&library))) return error;

    if ((error = FT_New_Memory_Face(library, params->buffer, params->bufferSize, params->faceIndex, &face))) goto end;

    if ((error = FT_Set_Pixel_Sizes(face, 0, params->fontSize))) goto end;

    FaceProperties *properties = malloc(sizeof(FaceProperties));
    if (!properties) {
        error = FT_Err_Out_Of_Memory;
        goto end;
    }

    properties->familyName = strdup(face->family_name ? face->family_name : "");
    if (!properties->familyName) {
        error = FT_Err_Out_Of_Memory;
        goto cleanup_props;
    }

    properties->styleName = strdup(face->style_name ? face->style_name : "");
    if (!properties->styleName) {
        error = FT_Err_Out_Of_Memory;
        goto cleanup_family;
    }

    properties->monospace = (bool) FT_IS_FIXED_WIDTH(face);
    properties->maxCharacterWidth = face->size->metrics.max_advance >> 6;
    properties->maxCharacterHeight = face->size->metrics.height >> 6;

    *outProperties = properties;

    FT_Done_Face(face);
    FT_Done_FreeType(library);

    return 0;

cleanup_family:
    free(properties->familyName);

cleanup_props:
    free(properties);

end:
    if (face) FT_Done_Face(face);
    if (library) FT_Done_FreeType(library);

    return error;
}

void freeFaceProperties(FaceProperties *properties) {
    if (properties->familyName) free(properties->familyName);
    if (properties->styleName) free(properties->styleName);
    free(properties);
}

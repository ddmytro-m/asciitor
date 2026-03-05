#include "freetype.h"
#include "freetype/freetype.h"
#include "freetype/fttypes.h"

#include <stdlib.h>
#include <string.h>

FT_Error initFaceWithParams(FT_Library *library, FT_Face *face, FaceParams *params) {
    FT_Error error;

    if ((error = FT_Init_FreeType(library))) return error;

    if ((error = FT_New_Memory_Face(*library, params->buffer, params->bufferSize, params->faceIndex, face))) {
        FT_Done_FreeType(*library);
        return error;
    }

    if ((error = FT_Set_Pixel_Sizes(*face, 0, params->fontSize))) {
        FT_Done_Face(*face);
        FT_Done_FreeType(*library);
        return error;
    }

    return 0;
}

FT_Error getFaceProperties(FaceParams *params, FaceProperties *out) {
    if (!out) return FT_Err_Bad_Argument;

    FT_Library library = NULL;
    FT_Face face = NULL;

    FT_Error error = 0;

    if ((error = initFaceWithParams(&library, &face, params))) return error;

    out->familyName = strdup(face->family_name ? face->family_name : "");
    out->styleName = strdup(face->style_name ? face->style_name : "");
    if (!out->familyName || !out->styleName) {
        error = FT_Err_Out_Of_Memory;
        goto end;
    }

    out->monospace = (bool) FT_IS_FIXED_WIDTH(face);
    out->maxCharacterWidth = face->size->metrics.max_advance >> 6;
    out->maxCharacterHeight = face->size->metrics.height >> 6;

end:
    FT_Done_Face(face);
    FT_Done_FreeType(library);

    return error;
}

void freeFaceProperties(FaceProperties *properties) {
    if (properties->familyName) free(properties->familyName);
    if (properties->styleName) free(properties->styleName);
}

FT_Error renderCharacter(FT_Face face, unsigned int character, RenderedCharacter **out) {
    RenderedCharacter *renderedChar = malloc(sizeof(RenderedCharacter));
    if (!renderedChar) return FT_Err_Out_Of_Memory;
    
    FT_Error error;
    if ((error = FT_Load_Char(face, character, FT_LOAD_RENDER))) goto error;

    FT_GlyphSlot slot = face->glyph;
    int width = slot->bitmap.width;
    int height = slot->bitmap.rows;

    renderedChar->bitmapWidth = width;
    renderedChar->bitmapHeight = height;
    renderedChar->advance = (slot->advance.x >> 6); // @COMBAK: check for vertical fonts

    int baseline = (face->size->metrics.ascender >> 6);
    renderedChar->leftShift = slot->bitmap_left;
    renderedChar->topShift = baseline - slot->bitmap_top;

    if (width > 0 && height > 0) {
        unsigned char *buf = malloc(sizeof(unsigned char) * width * height);
        if (!buf) {
            error = FT_Err_Out_Of_Memory;
            goto error;
        }

        for (int i = 0; i < height; i++) {
            memcpy(
                buf + i * width,
                slot->bitmap.buffer + (i * slot->bitmap.pitch),
                width
            );
        }
        renderedChar->bitmapBuffer = buf;
    } else {
        renderedChar->bitmapBuffer = NULL;
    }

    *out = renderedChar;

    return 0;

error:
    free(renderedChar);
    *out = NULL;

    return error;
}

FT_Error renderCharacters(FaceParams *params, unsigned int *characters, int length, RenderedCharacter **out, FT_Error *outErrors) {
    FT_Library library = NULL;
    FT_Face face = NULL;

    FT_Error error = 0;

    if ((error = initFaceWithParams(&library, &face, params))) return error;

    for (int i = 0; i < length; i++) {
        outErrors[i] = renderCharacter(face, characters[i], &out[i]);
    }

    FT_Done_Face(face);
    FT_Done_FreeType(library);

    return 0;
}

void freeRenderedCharacter(RenderedCharacter *renderedCharacter) {
    if (!renderedCharacter) return;
    if (renderedCharacter->bitmapBuffer) free(renderedCharacter->bitmapBuffer);
    free(renderedCharacter);
}

void freeRenderedCharacters(RenderedCharacter **renderedCharacters, int length) {
    if (!renderedCharacters) return;

    for (int i = 0; i < length; i++) {
        if (!renderedCharacters[i]) continue;

        freeRenderedCharacter(renderedCharacters[i]);
    }
}

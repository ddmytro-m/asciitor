#include "freetype.h"
#include "freetype/freetype.h"
#include "freetype/fttypes.h"

#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

static FT_Error initFaceWithParams(FT_Library *library, FT_Face *face, FaceParams params) {
    FT_Error error;

    error = FT_New_Memory_Face(
        *library, 
        params.fontParams.buffer, 
        params.fontParams.bufferSize, 
        params.faceIndex, 
        face
    );
    if (error) return error;

    return 0;
}

FT_Error getFont(FontParams params, FontProperties *out) {
    if (!out) return FT_Err_Bad_Argument;

    FT_Error error;

    FT_Library library = NULL;
    FT_Face face = NULL;

    FaceParams faceParams = { params, 0 };
    
    if ((error = FT_Init_FreeType(&library))) return error;
    if ((error = initFaceWithParams(&library, &face, faceParams))) {
        FT_Done_FreeType(library);
        return error;
    }

    out->facesAmount = face->num_faces;
    out->familyName = strdup(face->family_name ? face->family_name : "");
    if (!out->familyName) {
        error = FT_Err_Out_Of_Memory;
        goto end;
    }

end:
    FT_Done_Face(face);
    FT_Done_FreeType(library);

    return error;
}

void freeFont(FontProperties font) {
    if (font.familyName) free(font.familyName);
}


FT_Error getFace(FaceParams params, FaceProperties *out) {
    if (!out) return FT_Err_Bad_Argument;

    FT_Library library = NULL;
    FT_Face face = NULL;

    FT_Error error = 0;

    if ((error = FT_Init_FreeType(&library))) return error;
    if ((error = initFaceWithParams(&library, &face, params))) {
        FT_Done_FreeType(library);
        return error;
    }

    out->index = params.faceIndex;
    out->styleName = strdup(face->style_name ? face->style_name : "");
    if (!out->styleName) {
        error = FT_Err_Out_Of_Memory;
        goto end;
    }

    out->monospace = (bool) FT_IS_FIXED_WIDTH(face);

end:
    FT_Done_Face(face);
    FT_Done_FreeType(library);

    return error;
}

void freeFace(FaceProperties face) {
    if (face.styleName) free(face.styleName);
}

FT_Error getFontFaces(FontParams fontParams, FaceProperties **out, int *outLength) {
    if (!out || !outLength) return FT_Err_Bad_Argument;

    FT_Library library = NULL;
    FT_Face face = NULL;

    FT_Error error = 0;
    if ((error = FT_Init_FreeType(&library))) return error;

    int numFaces = 0, index = 0;
    FaceProperties properties;

    FaceParams faceParams = { fontParams };
    do {
        faceParams.faceIndex = index;
        if ((error = initFaceWithParams(&library, &face, faceParams))) goto cleanup;

        if (index == 0) {
            numFaces = face->num_faces;

            if (numFaces <= 0) {
                FT_Done_Face(face);
                goto cleanup;
            }

            *out = (FaceProperties*) calloc(numFaces, sizeof(FaceProperties));
            if (!*out) {
                FT_Done_Face(face);
                error = FT_Err_Out_Of_Memory;
                goto cleanup;
            }
        }

        (*out)[index].index = index;
        (*out)[index].styleName = strdup(face->style_name ? face->style_name : "");
        
        if (!(*out)[index].styleName) {
            FT_Done_Face(face);
            error = FT_Err_Out_Of_Memory;
            goto cleanup;
        }

        (*out)[index].monospace = (bool)FT_IS_FIXED_WIDTH(face);

        FT_Done_Face(face);
        face = NULL; // Reset so cleanup doesn't double-free
        index++;
    } while (index < numFaces);

    *outLength = numFaces;
    FT_Done_FreeType(library);
    return FT_Err_Ok;

cleanup:
    if (library) FT_Done_FreeType(library);

    if (*out) {
        freeFaces(*out, index); 
        *out = NULL;
    }

    *outLength = 0;

    return error;
}

void freeFaces(FaceProperties *faces, int length) {
    for (int i = 0; i < length; i++) {
        if (faces[i].styleName) free(faces[i].styleName);
    }
    free(faces);
}


static FT_Error renderGlyph(FT_Face face, FT_UInt glyph, RenderedCharacter *out) {    
    FT_Error error;
    if ((error = FT_Load_Glyph(face, glyph, FT_LOAD_RENDER))) return error;

    FT_GlyphSlot slot = face->glyph;
    int width = slot->bitmap.width;
    int height = slot->bitmap.rows;

    out->bitmapWidth = width;
    out->bitmapHeight = height;
    out->advance = (slot->advance.x >> 6); // @COMBAK: check for vertical fonts

    int baseline = (face->size->metrics.ascender >> 6);
    out->leftShift = slot->bitmap_left;
    out->topShift = baseline - slot->bitmap_top;

    if (width > 0 && height > 0) {
        unsigned char *buf = malloc(sizeof(unsigned char) * width * height);
        if (!buf) return FT_Err_Out_Of_Memory;

        for (int i = 0; i < height; i++) {
            memcpy(
                buf + i * width,
                slot->bitmap.buffer + (i * slot->bitmap.pitch),
                width
            );
        }
        out->bitmapBuffer = buf;
    } else {
        out->bitmapBuffer = NULL;
    }

    return 0;
}

FT_Error render(FaceParams faceParams, int fontSize, unsigned int *charcodes, int charcodesLength, RenderOutput *out) {
    if (fontSize <= 0) return FT_Err_Bad_Argument;

    FT_Library library = NULL;
    FT_Face face = NULL;
    FT_Error error;

    if ((error = FT_Init_FreeType(&library))) return error;

    if ((error = initFaceWithParams(&library, &face, faceParams))) goto cleanup;
    if ((error = FT_Set_Pixel_Sizes(face, 0, fontSize))) goto cleanup;

    int textHeight = (face->size->metrics.ascender - face->size->metrics.descender) >> 6;

    for (int i = 0; i < charcodesLength; i++) {
        unsigned int charcode = charcodes[i];
        out->characters[i].charcode = charcode; // define charcode for error matching

        FT_UInt glyph;
        glyph = FT_Get_Char_Index(face, charcode);
        if (glyph == 0) {
            out->errors[i] = FT_Err_Invalid_Glyph_Index;
            continue;
        }

        FT_Error err = renderGlyph(face, glyph, &out->characters[i]);
        if (err) {
            out->errors[i] = err;
            continue;
        }
    }

    out->length = charcodesLength;
    out->textHeight = textHeight;

cleanup:
    if (face) FT_Done_Face(face);
    if (library) FT_Done_FreeType(library);

    return error;
}

static void freeRenderedCharacter(RenderedCharacter renderedCharacter) {
    if (renderedCharacter.bitmapBuffer != NULL) {
        free(renderedCharacter.bitmapBuffer);
        renderedCharacter.bitmapBuffer = NULL; 
    }
}

void freeRendered(RenderOutput rendered) {
    for (int i = 0; i < rendered.length; i++) {
        freeRenderedCharacter(rendered.characters[i]);
    }
}

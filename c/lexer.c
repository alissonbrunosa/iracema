#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lexer.h"

static ir_token_t* new_token(enum ID id, char* literal, int length, int line, int column);
static int is_letter(char c);
static int is_digit(char c);
static enum ID keyword_lookup(char* literal, int length);
static void token_free(ir_token_t* tok);
static char lexer_peek(ir_lexer_t* lex);
static char lexer_peek_next(ir_lexer_t* lex);
static char lexer_advance(ir_lexer_t* lex);
static void lexer_skip_whitespace(ir_lexer_t* lex);
static ir_token_t* lexer_read_ident(ir_lexer_t* lex);
static ir_token_t* lexer_read_number(ir_lexer_t* lex);
static ir_token_t* lexer_read_string(ir_lexer_t* lex);

ir_lexer_t* new_lexer(char* source) {
    ir_lexer_t* lex = (ir_lexer_t *) malloc(sizeof(ir_lexer_t));
    if (lex == NULL) {
        printf("Memory allocation failed\n");
        exit(1);
    }

    lex->source = source;
    lex->current = source;
    lex->line = 1;
    lex->column = 1;

    return lex;
}

void lexer_free(ir_lexer_t* lex) {
    if (lex == NULL) {
        return;
    }

    if (lex->source != NULL) {
        free(lex->source);
    }

    free(lex);
}

ir_token_t* lexer_next(ir_lexer_t* lex){
    lexer_skip_whitespace(lex);

    char* current = lex->current;

    if (is_letter(*current)) {
        return lexer_read_ident(lex);
    }

    if (is_digit(*current)) {
        return lexer_read_number(lex);
    }

    int line = lex->line;
    int column = lex->column;

    switch(*current) {
        case '"':
            return lexer_read_string(lex);

        case '+':
            lexer_advance(lex);
            return new_token(TOKEN_PLUS, current, 1, line, column);

        case '-':
            lexer_advance(lex);
            if (lexer_peek(lex) == '>') {
                lexer_advance(lex);
                return new_token(TOKEN_ARROW, current, 2, line, column);
            }

            return new_token(TOKEN_MINUS, current, 1, line, column);

        case '*':
            lexer_advance(lex);
            return new_token(TOKEN_STAR, current, 1,line, column);

        case '/':
            lexer_advance(lex);
            return new_token(TOKEN_SLASH, current, 1, line, column);

        case '.':
            lexer_advance(lex);
            return new_token(TOKEN_DOT, current, 1, line, column);

        case ':':
            lexer_advance(lex);
            return new_token(TOKEN_COLON, current, 1, line, column);

        case ',':
            lexer_advance(lex);
            return new_token(TOKEN_COMMA, current, 1, line, column);

        case ';':
            lexer_advance(lex);
            return new_token(TOKEN_SEMICOLON, current, 1, line, column);

        case '(':
            lexer_advance(lex);
            return new_token(TOKEN_LEFT_PAREN, current, 1, line, column);

        case ')':
            lexer_advance(lex);
            return new_token(TOKEN_RIGHT_PAREN, current, 1, line, column);

        case '[':
            lexer_advance(lex);
            return new_token(TOKEN_LEFT_BRACKET, current, 1, line, column);

        case ']':
            lexer_advance(lex);
            return new_token(TOKEN_RIGHT_BRACKET, current, 1, line, column);

        case '{':
            lexer_advance(lex);
            return new_token(TOKEN_LEFT_BRACE, current, 1, line, column);

        case '}':
            lexer_advance(lex);
            return new_token(TOKEN_RIGHT_BRACE, current, 1, line, column);

        case '=':
            lexer_advance(lex);
            if (lexer_peek(lex) == '=') {
                lexer_advance(lex);
                return new_token(TOKEN_EQUAL, current, 2, line, column);
            }

            return new_token(TOKEN_ASSIGN, current, 1, line, column);

        case '!':
            lexer_advance(lex);
            if (lexer_peek(lex) == '=') {
                lexer_advance(lex);
                return new_token(TOKEN_NOT_EQUAL, current, 2, line, column);
            }

            return new_token(TOKEN_BANG, current, 1, line, column);

        case '<':
            lexer_advance(lex);
            if (lexer_peek(lex) == '=') {
                lexer_advance(lex);
                return new_token(TOKEN_LESS_EQUAL, current, 2, line, column);
            }

            return new_token(TOKEN_LESS, current, 1, line, column);

        case '>':
            lexer_advance(lex);
            if (lexer_peek(lex) == '=') {
                lexer_advance(lex);
                return new_token(TOKEN_GREAT_EQUAL, current, 2, line, column);
            }

            return new_token(TOKEN_GREAT, current, 1, line, column);

        case '\0':
            return new_token(TOKEN_EOF, current, 0, line, column);

        default:
            lexer_advance(lex);
            return new_token(TOKEN_UNKNOWN, current, 1, line, column);
    }
}

const char* token_names[] = {
    "TOKEN_EOF",
    "TOKEN_IDENT",
    "TOKEN_NOT",
    "TOKEN_DOT",
    "TOKEN_PLUS",
    "TOKEN_MINUS",
    "TOKEN_STAR",
    "TOKEN_SLASH",
    "TOKEN_COLON",
    "TOKEN_COMMA",
    "TOKEN_BANG",
    "TOKEN_ARROW",
    "TOKEN_ASSIGN",
    "TOKEN_EQUAL",
    "TOKEN_NOT_EQUAL",
    "TOKEN_LESS",
    "TOKEN_LESS_EQUAL",
    "TOKEN_GREAT",
    "TOKEN_GREAT_EQUAL",
    "TOKEN_SEMICOLON",
    "TOKEN_LEFT_PAREN",
    "TOKEN_RIGHT_PAREN",
    "TOKEN_LEFT_BRACE",
    "TOKEN_RIGHT_BRACE",
    "TOKEN_LEFT_BRACKET",
    "TOKEN_RIGHT_BRACKET",
    "TOKEN_STRING",
    "TOKEN_NUMBER",
    "TOKEN_OR",
    "TOKEN_AND",
    "TOKEN_IF",
    "TOKEN_ELSE",
    "TOKEN_SWITCH",
    "TOKEN_CASE",
    "TOKEN_DEFAULT",
    "TOKEN_FOR",
    "TOKEN_IN",
    "TOKEN_NONE",
    "TOKEN_OBJECT",
    "TOKEN_IS",
    "TOKEN_VAR",
    "TOKEN_FUN",
    "TOKEN_THIS",
    "TOKEN_RETURN",
    "TOKEN_SUPER",
    "TOKEN_TRUE",
    "TOKEN_FALSE",
    "TOKEN_WHILE",
    "TOKEN_STOP",
    "TOKEN_NEXT",
    "TOKEN_UNKNOWN"
};

static int is_letter(char c) {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_';
}

static int is_digit(char c) {
  return c >= '0' && c <= '9';
}

static int is_alphanumetic(char c) {
    return is_digit(c) || is_letter(c);
}

static enum ID keyword_lookup(char* literal, int length) {
    switch(length) {
        case 2:
            if(memcmp(literal, "or", length) == 0) return TOKEN_OR;
            if(memcmp(literal, "if", length) == 0) return TOKEN_IF;
            if(memcmp(literal, "in", length) == 0) return TOKEN_IN;
            if(memcmp(literal, "is", length) == 0) return TOKEN_IS;
            break;

        case 3:
            if(memcmp(literal, "and", length) == 0) return TOKEN_AND;
            if(memcmp(literal, "var", length) == 0) return TOKEN_VAR;
            if(memcmp(literal, "fun", length) == 0) return TOKEN_FUN;
            if(memcmp(literal, "for", length) == 0) return TOKEN_FOR;
            break;

        case 4:
            if(memcmp(literal, "else", length) == 0) return TOKEN_ELSE;
            if(memcmp(literal, "case", length) == 0) return TOKEN_CASE;
            if(memcmp(literal, "none", length) == 0) return TOKEN_NONE;
            if(memcmp(literal, "this", length) == 0) return TOKEN_THIS;
            if(memcmp(literal, "true", length) == 0) return TOKEN_TRUE;
            if(memcmp(literal, "stop", length) == 0) return TOKEN_STOP;
            if(memcmp(literal, "next", length) == 0) return TOKEN_NEXT;
            break;

        case 5:
            if(memcmp(literal, "super", length) == 0) return TOKEN_SUPER;
            if(memcmp(literal, "false", length) == 0) return TOKEN_FALSE;
            if(memcmp(literal, "while", length) == 0) return TOKEN_WHILE;
            break;

        case 6:
            if(memcmp(literal, "switch", length) == 0) return TOKEN_SWITCH;
            if(memcmp(literal, "object", length) == 0) return TOKEN_OBJECT;
            if(memcmp(literal, "return", length) == 0) return TOKEN_RETURN;
            break;

        case 7:
            if(memcmp(literal, "default", length) == 0) return TOKEN_DEFAULT;
            break;

        default:
            return TOKEN_IDENT;
    }

    return TOKEN_IDENT;
}

static ir_token_t* new_token(enum ID id, char* literal, int length, int line, int column) {
    ir_token_t* tok = (ir_token_t *) malloc(sizeof(ir_token_t));
    if (tok == NULL) {
        printf("Memory allocation failed\n");
        exit(1);
    }

    tok->id = id;
    tok->literal = literal;
    tok->length = length;
    tok->line = line;
    tok->column = column;

    return tok;
}

static void token_free(ir_token_t* tok) {
    if (tok == NULL) {
        return;
    }

    free(tok);
}

static char lexer_peek(ir_lexer_t* lex) {
    return *lex->current;
}

static char lexer_peek_next(ir_lexer_t* lex) {
    if (lexer_peek(lex) == '\0') {
        return '\0';
    }

    return lex->current[1];
}

static char lexer_advance(ir_lexer_t* lex) {
    char current = lexer_peek(lex);
    lex->current++;

    if (current == '\n') {
        lex->line++;
        lex->column = 1;
    } else {
        lex->column++;
    }

    return current;
}

static void lexer_skip_whitespace(ir_lexer_t* lex) {
    while (1) {
        char current = lexer_peek(lex);
        switch(current) {
            case ' ':
            case '\r':
            case '\t':
            case '\n':
                lexer_advance(lex);
                break;

            default:
                return;
        }
    }
}

static ir_token_t* lexer_read_ident(ir_lexer_t* lex) {
    int line = lex->line;
    int column = lex->column;
    char* start = lex->current;

    for (char chr = lexer_peek(lex); is_alphanumetic(chr); chr = lexer_peek(lex)) {
        lexer_advance(lex);
    }

    int length = (int) (lex->current - start);
    enum ID id = keyword_lookup(start, length);
    return new_token(id, start, length, line, column);
}

static ir_token_t* lexer_read_number(ir_lexer_t* lex) {
    int line = lex->line;
    int column = lex->column;
    char* start = lex->current;

    for (char chr = lexer_peek(lex); is_digit(chr); chr = lexer_peek(lex)) {
        lexer_advance(lex);
    }

    char next = lexer_peek_next(lex);
    if (lexer_peek(lex) == '.' && is_digit(next)) {
        lexer_advance(lex);
        for (char chr = lexer_peek(lex); is_digit(chr); chr = lexer_peek(lex)) {
            lexer_advance(lex);
        }
    }

    int length = (int) (lex->current - start);
    return new_token(TOKEN_NUMBER, start, length, line, column);
}

static ir_token_t* lexer_read_string(ir_lexer_t* lex) {
    int line = lex->line;
    int column = lex->column;
    char* start = lex->current;

    lexer_advance(lex);
    char chr = lexer_peek(lex);

    while (chr != '"' && chr != '\n' && chr != '\0') {
        lexer_advance(lex);
        chr = lexer_peek(lex);
    }

    if (chr == '\n' || chr == '\0') {
        printf("string not terminated\n");
        exit(1);
    }

    lexer_advance(lex);
    int length = (int) (lex->current - start);
    return new_token(TOKEN_STRING, start, length, line, column);
}

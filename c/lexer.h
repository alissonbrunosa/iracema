#ifndef LEXER_H
#define LEXER_H

#include <stdlib.h>

extern const char* token_names[];
enum ID {
    TOKEN_EOF,

    TOKEN_IDENT,
    TOKEN_NOT,
    TOKEN_DOT,
    TOKEN_PLUS,
    TOKEN_MINUS,
    TOKEN_STAR,
    TOKEN_SLASH,
    TOKEN_COLON,
    TOKEN_COMMA,
    TOKEN_BANG,
    TOKEN_ARROW,
    TOKEN_ASSIGN,
    TOKEN_EQUAL,
    TOKEN_NOT_EQUAL,
    TOKEN_LESS,
    TOKEN_LESS_EQUAL,
    TOKEN_GREAT,
    TOKEN_GREAT_EQUAL,
    TOKEN_SEMICOLON,
    TOKEN_LEFT_PAREN,
    TOKEN_RIGHT_PAREN,
    TOKEN_LEFT_BRACE,
    TOKEN_RIGHT_BRACE,
    TOKEN_LEFT_BRACKET,
    TOKEN_RIGHT_BRACKET,

    TOKEN_STRING,
    TOKEN_NUMBER,

    TOKEN_OR,
    TOKEN_AND,
    TOKEN_IF,
    TOKEN_ELSE,
    TOKEN_SWITCH,
    TOKEN_CASE,
    TOKEN_DEFAULT,
    TOKEN_FOR,
    TOKEN_IN,
    TOKEN_NONE,
    TOKEN_OBJECT,
    TOKEN_IS,
    TOKEN_VAR,
    TOKEN_FUN,
    TOKEN_THIS,
    TOKEN_RETURN,
    TOKEN_SUPER,
    TOKEN_TRUE,
    TOKEN_FALSE,
    TOKEN_WHILE,
    TOKEN_STOP,
    TOKEN_NEXT,
    TOKEN_UNKNOWN
};

struct token {
    enum ID id;
    char* literal;
    int length;
    int line;
    int column;
};

typedef struct token ir_token_t;

void token_free(ir_token_t* tok);

struct lexer {
    char* source;
    char* current;
    int line;
    int column;
};

typedef struct lexer ir_lexer_t;

char lexer_peek(ir_lexer_t* lex);

char lexer_advance(ir_lexer_t* lex);

void lexer_skip_whitespace(ir_lexer_t* lex);

ir_token_t* lexer_next(ir_lexer_t* lex);

ir_token_t* lexer_read_ident(ir_lexer_t* lex);

ir_token_t* lexer_read_number(ir_lexer_t* lex);

ir_token_t* lexer_read_string(ir_lexer_t* lex);

ir_lexer_t* new_lexer(char* source);

void lexer_free(ir_lexer_t* lex);

#endif

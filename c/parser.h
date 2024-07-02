#ifndef PARSER_H
#define PARSER_H

#include "lexer.h"

#define NOOP 0
#define UNARY_NOT 1
#define UNARY_PLUS 2
#define UNARY_MINUS 3
#define BINARY_OR 4
#define BINARY_AND 5
#define BINARY_PLUS 6
#define BINARY_MINUS 7
#define BINARY_SLASH 8
#define BINARY_STAR 9
#define BINARY_GREAT 10
#define BINARY_GREAT_EQUAL 11
#define BINARY_LESS 12
#define BINARY_LESS_EQUAL 13
#define BINARY_EQUAL 14
#define BINARY_NOT_EQUAL 15

enum node_type {
    OBJECT_NODE,
    FUNCTION_NODE,
    IF_NODE,
    FOR_NODE,
    SWITCH_NODE,
    CASE_NODE,
    ASSIGN_NODE,
    BOOL_NODE,
    THIS_NODE,
    STRING_NODE,
    NUMBER_NODE,
    NEXT_NODE,
    STOP_NODE,
    RETURN_NODE,
    BLOCK_NODE,
    PARAMETER_NODE,
    INDEX_ACCESS_NODE,
    MEMBER_ACCESS_NODE,

    CALL_NODE,
    IDENT_NODE,
    GROUP_NODE,
    UNARY_NODE,
    BINARY_NODE,
    STRING_EXPR,
    NUMBER_EXPR,
};

typedef struct node {
    enum node_type type;
    int line;
    int column;
} ir_node_t;

struct node_list {
    int size;
    int capacity;
    struct node** items;
};

struct ident_node {
    ir_node_t base;
    const char* value;
};

struct parameter {
    ir_node_t base;
    struct ident_node* name;
    struct node* value;
};

struct fun_node {
    ir_node_t base;
    struct ident_node* name;
    struct node_list* params;
    struct node_list* body;
};

struct call_node {
    ir_node_t base;
    struct node* func;
    struct node_list* arguments;
};

struct block_node {
    ir_node_t base;
    struct node_list* nodes;
};

struct object_node {
    ir_node_t base;
    struct ident_node* name;
    struct ident_node* parent;
    struct node_list* functions;
};

struct if_node {
    ir_node_t base;
    struct node* cond;
    struct node_list* then;
    struct node* consequent;
};

struct for_node {
    ir_node_t base;
    struct ident_node* element;
    struct node* collection;
    struct node_list* body;
};

struct case_clause_node {
    ir_node_t base;
    struct node_list* expr_list;
    struct node_list* body;
};

struct switch_node {
    ir_node_t base;
    struct node* tag;
    struct node_list* cases;
};

struct number_node {
    ir_node_t base;
    double value;
};

struct string_node {
    ir_node_t base;
    char* value;
};

struct bool_node {
    ir_node_t base;
    unsigned char value;
};

struct assign_node {
    ir_node_t base;
    struct node_list* lhs;
    struct node_list* rhs;
};

struct binary_node {
    ir_node_t base;
    int operator;
    struct node* lhs;
    struct node* rhs;
};

struct unary_node {
    ir_node_t base;
    int operator;
    ir_node_t* expr;
};

struct group_node {
    ir_node_t base;
    struct node* expr;
};

// node that only wrap an expression
// return 10;
// (10 + 10);
struct expr_wrapper_node {
    ir_node_t base;
    struct node* expr;
};

struct member_access_node {
    ir_node_t base;
    struct node* object;
    struct ident_node* member;
};

struct index_access_node {
    ir_node_t base;
    struct node* expr;
    struct node* index;
};

struct file {
    int size;
    int capacity;
    ir_node_t** nodes;
};

typedef struct parser {
    struct lexer* lex;
    struct file* file;
    ir_token_t* current;
} ir_parser_t;

struct file* new_file();
void file_free(struct file*);
void file_add_node(struct file*, ir_node_t*);

ir_parser_t* new_parser(char*);
void free_parser(ir_parser_t*);
struct file* start(ir_parser_t*);
void node_free(ir_node_t*);
#endif

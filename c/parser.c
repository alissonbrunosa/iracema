#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <assert.h>

#include "parser.h"

#define AS_NODE(stmt) (ir_node_t*) stmt
#define NEW_NODE(node_struct, node_type) (node_struct*) new_node(sizeof(node_struct), node_type)

static void advance(ir_parser_t*);
static int consume(ir_parser_t*, enum ID);
static int unary_operator(ir_token_t*);
static int binary_operator(ir_token_t*);
static int token_precedence(ir_token_t*);

static struct number_node* parse_number(ir_parser_t*);
static struct ident_node* parse_ident(ir_parser_t*);
static struct group_node* parse_group_expr(ir_parser_t*);
static struct bool_node* parse_bool(ir_parser_t*);
static struct string_node* parse_string(ir_parser_t*);
static struct node_list* parse_parameter_list(ir_parser_t*);
static struct fun_node* parse_function(ir_parser_t*);
static struct object_node* parse_object(ir_parser_t*);
static struct if_node* parse_if_node(ir_parser_t*);
static struct for_node* parse_for_node(ir_parser_t*);
static struct switch_node* parse_switch_node(ir_parser_t*);
static ir_node_t* parse_case_clause(ir_parser_t*);
static struct expr_wrapper_node* parse_return_node(ir_parser_t*);
static struct node_list* parse_block(ir_parser_t*);
static ir_node_t* parse_stmt(ir_parser_t*);
static struct node_list* parse_stmt_list(ir_parser_t*);
static ir_node_t* parse_simple_stmt(ir_parser_t*);
static ir_node_t* parse_assigment(ir_parser_t*, struct node_list*);
static struct node_list* parse_expr_list(ir_parser_t*);
static ir_node_t* parse_expr(ir_parser_t*);
static ir_node_t* parse_binary_expr(ir_parser_t*, int);
static ir_node_t* parse_unary_expr(ir_parser_t*);
static ir_node_t* parse_primary_expr(ir_parser_t*);
static ir_node_t* parse_operand(ir_parser_t*);
static ir_node_t* parse_member_access(ir_parser_t*, ir_node_t*);
static ir_node_t* parse_method_call(ir_parser_t*, ir_node_t*);
static ir_node_t* parse_index_access(ir_parser_t*, ir_node_t*);
static struct node_list* parse_argument_list(ir_parser_t*);

static void advance(ir_parser_t *p) {
    p->current = lexer_next(p->lex);
}

static int consume(ir_parser_t *p, enum ID id) {
    if (p->current->id == id) {
        advance(p);
        return 1;
    }

    return 0;
}

static int unary_operator(ir_token_t* t) {
    switch (t->id) {
        case TOKEN_BANG: return UNARY_NOT;
        case TOKEN_PLUS: return UNARY_PLUS;
        case TOKEN_MINUS: return UNARY_MINUS;
        default: return NOOP;
    }
}

static int binary_operator(ir_token_t* t) {
    switch (t->id) {
        case TOKEN_OR: return BINARY_OR;
        case TOKEN_AND: return BINARY_AND;
        case TOKEN_PLUS: return BINARY_PLUS;
        case TOKEN_MINUS: return BINARY_MINUS;
        case TOKEN_SLASH: return BINARY_SLASH;
        case TOKEN_STAR: return BINARY_STAR;
        case TOKEN_GREAT: return BINARY_GREAT;
        case TOKEN_GREAT_EQUAL: return BINARY_GREAT_EQUAL;
        case TOKEN_LESS: return BINARY_LESS;
        case TOKEN_LESS_EQUAL: return BINARY_LESS_EQUAL;
        case TOKEN_EQUAL: return BINARY_EQUAL;
        case TOKEN_NOT_EQUAL: return BINARY_NOT_EQUAL;
        default: return NOOP;
    }
}

static int token_precedence(ir_token_t* t) {
    switch (t->id) {
        case TOKEN_OR:
        case TOKEN_AND:
            return 2;
        case TOKEN_EQUAL:
        case TOKEN_NOT_EQUAL:
        case TOKEN_LESS:
        case TOKEN_LESS_EQUAL:
        case TOKEN_GREAT:
        case TOKEN_GREAT_EQUAL:
            return 3;
        case TOKEN_MINUS:
        case TOKEN_PLUS:
            return 4;
        case TOKEN_SLASH:
        case TOKEN_STAR:
            return 5;
        default:
            return 0;
    }
}

static ir_node_t* new_node(size_t size, enum node_type type) {
    ir_node_t* node = (ir_node_t*) malloc(size);
    if (node == NULL) {
        printf("memory allocation failed\n");
        exit(1);
    }

    node->type = type;
    return node;
}

static void node_list_free(struct node_list* list) {
    if (list->capacity > 0) {
        for (int i = 0; i < list->size; i++) {
            node_free(list->items[i]);
        }
    }

    free(list);
}

static struct node_list* create_node_list() {
    struct node_list* list = calloc(1, sizeof(struct node_list));
    if (list == NULL) {
        perror("Failed to allocate memory for the node_list");
        exit(1);
    }

    return list;
}

static void node_list_append(struct node_list* list, ir_node_t* node) {
    if (list->size == list->capacity) {
        list->capacity = list->capacity == 0 ? 4 : list->capacity * 2;
        list->items = (ir_node_t **) realloc(list->items, sizeof(ir_node_t *) * list->capacity);
    }
    list->items[list->size++] = node;
}

ir_token_t* expect(ir_parser_t *p, enum ID id) {
    ir_token_t* current = p->current;

    if (current->id != id) {
        printf("expected %s, got %s\n", token_names[id], token_names[current->id]);
        exit(1);
    }

    advance(p);
    return current;
}

static struct number_node* parse_number(ir_parser_t* p) {
    ir_token_t* t = expect(p, TOKEN_NUMBER);
    char *endptr;

    char* literal = strndup(t->literal, t->length);
    if (literal == NULL) {
        printf("failed to duplicated number literal");
        exit(1);
    }

    struct number_node* node = NEW_NODE(struct number_node, NUMBER_NODE);

    node->value = strtod(literal, &endptr);
    if (endptr == literal) {
        printf("failed to parser '%s' into double.\n", literal);
        exit(1);
    }

    return node;
}

static struct ident_node* parse_ident(ir_parser_t* p) {
    ir_token_t* t = expect(p, TOKEN_IDENT);

    struct ident_node* ident = NEW_NODE(struct ident_node, IDENT_NODE);
    ident->value = strndup(t->literal, t->length);
    return ident;
}

static struct group_node* parse_group_expr(ir_parser_t* p) {
    expect(p, TOKEN_LEFT_PAREN);
    struct group_node* node = NEW_NODE(struct group_node, GROUP_NODE);
    node->expr = parse_expr(p);
    expect(p, TOKEN_RIGHT_PAREN);

    return node;
}

static struct bool_node* parse_bool(ir_parser_t* p) {
    struct bool_node* node = NEW_NODE(struct bool_node, BOOL_NODE);
    node->value = p->current->id == TOKEN_TRUE;

    advance(p);
    return node;
}

static struct string_node* parse_string(ir_parser_t* p) {
    ir_token_t* t = expect(p, TOKEN_STRING);

    struct string_node* node = NEW_NODE(struct string_node, STRING_NODE);

    node->value = strndup(t->literal, t->length);
    if (node->value == NULL) {
        printf("failed to dup token literal into string node\n");
        exit(1);
    }

    return node;
}

static struct object_node* parse_object(ir_parser_t *p) {
    expect(p, TOKEN_OBJECT);

    struct object_node* node = NEW_NODE(struct object_node, OBJECT_NODE);

    node->name = parse_ident(p);
    if (consume(p, TOKEN_IS)) {
        node->parent = parse_ident(p);
    }

    node->functions = create_node_list();

    expect(p, TOKEN_LEFT_BRACE);
    while (p->current->id == TOKEN_FUN) {
        struct fun_node* func = parse_function(p);
        node_list_append(node->functions, AS_NODE(func));
    }

    expect(p, TOKEN_RIGHT_BRACE);
    return node;
}

static struct fun_node* parse_function(ir_parser_t* p) {
    expect(p, TOKEN_FUN);

    struct fun_node* node = NEW_NODE(struct fun_node, FUNCTION_NODE);
    node->name = parse_ident(p);
    node->params = parse_parameter_list(p);
    node->body = parse_block(p);
    return node;
}

static struct node_list* parse_parameter_list(ir_parser_t* p) {
    expect(p, TOKEN_LEFT_PAREN);

    struct node_list* list = create_node_list();
    while (p->current->id != TOKEN_RIGHT_PAREN && p->current->id != TOKEN_EOF) {
        struct parameter* node = NEW_NODE(struct parameter, PARAMETER_NODE);
        node->name = parse_ident(p);

        if (consume(p, TOKEN_ASSIGN)) {
            node->value = parse_expr(p);
        }

        consume(p, TOKEN_COMMA);
        node_list_append(list, AS_NODE(node));
    }

    expect(p, TOKEN_RIGHT_PAREN);
    return list;
}

ir_parser_t* new_parser(char* source) {
    ir_parser_t* p = (ir_parser_t*) malloc(sizeof(ir_parser_t));
    if (p == NULL) {
        printf("memory allocation failed\n");
        exit(1);
    }

    p->lex = new_lexer(source);
    p->file = new_file();
    advance(p);

    return p;
}

void free_parser(ir_parser_t *p) {
    if (p == NULL) {
        return;
    }

    lexer_free(p->lex);
    file_free(p->file);
    free(p);
}

struct file* start(ir_parser_t *p) {
    while(p->current->id != TOKEN_EOF) {
        file_add_node(p->file, parse_stmt(p));
    }

    expect(p, TOKEN_EOF);
    return p->file;
}

static struct node_list* parse_stmt_list(ir_parser_t* p) {
    struct node_list* list = create_node_list();
    while (p->current->id != TOKEN_RIGHT_BRACE && p->current->id != TOKEN_CASE && p->current->id != TOKEN_DEFAULT && p->current->id != TOKEN_EOF) {
        node_list_append(list, parse_stmt(p));
    }

    return list;
}

static ir_node_t* parse_stmt(ir_parser_t* p) {
    while(p->current->id != TOKEN_EOF) {
        switch (p->current->id) {
            case TOKEN_OBJECT:
                return (ir_node_t *) parse_object(p);

            case TOKEN_FUN:
                return (ir_node_t *) parse_function(p);

            case TOKEN_IF:
                return (ir_node_t *) parse_if_node(p);

            case TOKEN_FOR:
                return (ir_node_t *) parse_for_node(p);

            case TOKEN_SWITCH:
                return (ir_node_t *) parse_switch_node(p);

            case TOKEN_RETURN:
                return (ir_node_t *) parse_return_node(p);

            case TOKEN_STOP:
            case TOKEN_NEXT:
            case TOKEN_THIS:
                parse_object(p);
                break;

            default:
                return parse_simple_stmt(p);
        }
    }

    return NULL;
}

static struct if_node* parse_if_node(ir_parser_t* p) {
    expect(p, TOKEN_IF);

    struct if_node* nd = NEW_NODE(struct if_node, IF_NODE);

    nd->cond = parse_expr(p);
    nd->then = parse_block(p);

    if (consume(p, TOKEN_ELSE)) {
        switch (p->current->id) {
            case TOKEN_IF:
                struct if_node* elsif = parse_if_node(p);
                nd->consequent = AS_NODE(elsif);
                break;

            case TOKEN_LEFT_BRACE:
                struct node_list* block = parse_block(p);
                nd->consequent = AS_NODE(block);
                break;

            default:
                printf("expected '{' or an if statement\n");
                exit(1);
        }
    }

    return nd;
}

static struct for_node* parse_for_node(ir_parser_t* p) {
    expect(p, TOKEN_FOR);

    struct for_node* node = NEW_NODE(struct for_node, FOR_NODE);
    node->element = parse_ident(p);
    expect(p, TOKEN_IN);
    node->collection = parse_expr(p);
    node->body = parse_block(p);

    return node;
}

///
// switch_statement ::= "switch" expression "{" { switch_case } "}"
// switch_case ::= case_clause | default_clause
// case_clause ::= "case" expression_list ":" statement_list
// default_clause ::= "default" ":" statement_list
///
static struct switch_node* parse_switch_node(ir_parser_t* p) {
    expect(p, TOKEN_SWITCH);

    struct switch_node* node = (struct switch_node *) NEW_NODE(struct switch_node, SWITCH_NODE);
    node->tag = parse_expr(p);
    node->cases = create_node_list();

    expect(p, TOKEN_LEFT_BRACE);
    while (p->current->id != TOKEN_RIGHT_BRACE && p->current->id != TOKEN_EOF) {
        node_list_append(node->cases, parse_case_clause(p));
    }

    expect(p, TOKEN_RIGHT_BRACE);
    return node;
}

static ir_node_t* parse_case_clause(ir_parser_t* p) {
    struct case_clause_node* node = NEW_NODE(struct case_clause_node, CASE_NODE);

    if (consume(p, TOKEN_CASE)) {
        node->expr_list = parse_expr_list(p);
    } else if (consume(p, TOKEN_DEFAULT)) {
        node->expr_list = NULL;
    } else {
        printf("expecting 'case' or 'default', but got %s\n", token_names[p->current->id]);
        exit(1);
    }

    expect(p, TOKEN_COLON);
    node->body = parse_stmt_list(p);

    return (ir_node_t *) node;
}

static struct expr_wrapper_node* parse_return_node(ir_parser_t* p) {
    expect(p, TOKEN_RETURN);

    struct expr_wrapper_node* node = NEW_NODE(struct expr_wrapper_node, RETURN_NODE);
    node->expr = parse_expr(p);
    return node;
}

static struct node_list* parse_block(ir_parser_t* p) {
    expect(p, TOKEN_LEFT_BRACE);

    struct node_list* block = create_node_list();
    while (p->current->id != TOKEN_RIGHT_BRACE && p->current->id != TOKEN_EOF) {
        node_list_append(block, parse_stmt(p));
    }

    expect(p, TOKEN_RIGHT_BRACE);
    return block;
}

static ir_node_t* parse_simple_stmt(ir_parser_t* p) {
    struct node_list* lhs = parse_expr_list(p);

    switch (p->current->id) {
        case TOKEN_ASSIGN:
            return parse_assigment(p, lhs);

        default:
            ir_node_t* node = lhs->items[lhs->size-1];
            free(lhs);
            return node;
    }
}

static ir_node_t* parse_assigment(ir_parser_t* p, struct node_list* lhs) {
    assert(lhs != NULL);

    expect(p, TOKEN_ASSIGN);
    struct assign_node* node = NEW_NODE(struct assign_node, ASSIGN_NODE);
    node->lhs = lhs;
    node->rhs = parse_expr_list(p);
    return AS_NODE(node);
}

static struct node_list* parse_expr_list(ir_parser_t* p) {
    struct node_list* list = create_node_list();

    node_list_append(list, parse_expr(p));
    while (consume(p, TOKEN_COMMA)) {
        node_list_append(list, parse_expr(p));
    }

    return list;
}

static ir_node_t* parse_expr(ir_parser_t* p) {
    return parse_binary_expr(p, 0);
}

static ir_node_t* parse_binary_expr(ir_parser_t* p, int prec) {
    ir_node_t* expr = parse_unary_expr(p);

    while (1) {
        int current_prec = token_precedence(p->current);
        if (current_prec <= prec) {
            break;
        }

        ir_token_t* t = expect(p, p->current->id);

        struct binary_node* node = NEW_NODE(struct binary_node, BINARY_NODE);
        node->lhs = expr;
        node->operator = binary_operator(t);
        node->rhs = parse_binary_expr(p, current_prec);
        expr = (ir_node_t *) node;
    }

    return expr;
}

static ir_node_t* parse_unary_expr(ir_parser_t* p) {
    switch (p->current->id) {
        case TOKEN_NOT:
        case TOKEN_PLUS:
        case TOKEN_MINUS:
            struct unary_node* node = NEW_NODE(struct unary_node, UNARY_NODE);
            node->operator = unary_operator(p->current);
            advance(p);
            node->expr = parse_unary_expr(p);
            return AS_NODE(node);

        default:
            return parse_primary_expr(p);
    }
}

static ir_node_t* parse_primary_expr(ir_parser_t* p) {
    ir_node_t* operand = parse_operand(p);

    while (1) {
        switch (p->current->id) {
            case TOKEN_DOT:
                operand = parse_member_access(p, operand);
                break;

            case TOKEN_LEFT_PAREN:
                operand = parse_method_call(p, operand);
                break;

            case TOKEN_LEFT_BRACKET:
                operand = parse_index_access(p, operand);
                break;

            default:
                return operand;
        }
    }
}

static ir_node_t* parse_operand(ir_parser_t* p) {
    switch (p->current->id) {
        case TOKEN_NUMBER:
            return (ir_node_t *) parse_number(p);

        case TOKEN_STRING:
            return (ir_node_t *) parse_string(p);

        case TOKEN_TRUE:
        case TOKEN_FALSE:
            return (ir_node_t *) parse_bool(p);

        case TOKEN_IDENT:
            return (ir_node_t *) parse_ident(p);

        case TOKEN_LEFT_PAREN:
            return (ir_node_t *) parse_group_expr(p);

        default:
            printf("bad expr: %s\n", token_names[p->current->id]);
            exit(1);
    }
}

static ir_node_t* parse_member_access(ir_parser_t* p, ir_node_t* object) {
    assert(object != NULL);

    expect(p, TOKEN_DOT);
    struct member_access_node* node = NEW_NODE(struct member_access_node, MEMBER_ACCESS_NODE);
    node->object = object;
    node->member = parse_ident(p);
    return (ir_node_t *) node;
}

static ir_node_t* parse_method_call(ir_parser_t* p, ir_node_t* func) {
    struct call_node* node = NEW_NODE(struct call_node, CALL_NODE);
    node->func = func;
    node->arguments = parse_argument_list(p);

    return (ir_node_t *) node;
}

static ir_node_t* parse_index_access(ir_parser_t* p, ir_node_t* expr) {
    assert(expr != NULL);

    expect(p, TOKEN_LEFT_BRACKET);
    struct index_access_node* node = NEW_NODE(struct index_access_node, INDEX_ACCESS_NODE);
    node->expr = expr;
    node->index = parse_expr(p);
    expect(p, TOKEN_RIGHT_BRACKET);

    return (ir_node_t *) node;
}

static struct node_list* parse_argument_list(ir_parser_t *p) {
    expect(p, TOKEN_LEFT_PAREN);
    if (consume(p, TOKEN_RIGHT_PAREN)) {
        return NULL;
    }

    struct node_list* list = parse_expr_list(p);
    expect(p, TOKEN_RIGHT_PAREN);
    return list;
}

struct file* new_file() {
    struct file* f = (struct file*) malloc(sizeof(struct file));
    if (f == NULL) {
        printf("Memory allocation failed\n");
        exit(1);
    }
    f->size = 0;
    f->capacity = 0;
    f->nodes = NULL;

    return f;
}

void file_free(struct file* f) {
    if (f == NULL) {
        return;
    }

    if (f->nodes != NULL) {
        for(int i = 0; i < f->size; i++) {
            node_free(f->nodes[i]);
        }
    }

    free(f);
}

void file_add_node(struct file* f, ir_node_t* node) {
    if (f->nodes == NULL) {
        f->nodes = (ir_node_t**) malloc(3 * sizeof(ir_node_t*));
        f->capacity = 3;
    } else if (f->capacity < f->size + 1) {
        f->capacity *= 2;
        f->nodes = (ir_node_t**) realloc(f->nodes, f->capacity * sizeof(ir_node_t*));
    }

    f->nodes[f->size++] = node;
}

void node_free(ir_node_t* node) {
    switch(node->type) {
        case NUMBER_NODE:
            free(node);
            break;

        case STRING_NODE:
            struct string_node* s = (struct string_node*) node;
            if (s->value != NULL) {
                free(s->value);
            }
            free(s);
            break;

        default:
            free(node);
    }
}

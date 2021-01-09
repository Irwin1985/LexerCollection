from enum import Enum

# exception handling
class Error(Exception):
    def __init__(self, error_code=None, token=None, message=None):
        self.error_code = error_code
        self.token = token
        self.message = f'{self.__class__.__name__}: {message}'

class LexerError(Error):
    pass

class TokenType(Enum):
    # single character token types
    PLUS        = '+'
    MINUS       = '-'
    MUL         = '*'
    FLOAT_DIV   = '/'
    LPAREN      = '('
    RPAREN      = ')'
    LBRACE      = '{'
    RBRACE      = '}'
    SEMI        = ';'
    DOT         = '.'
    COLON       = ':'
    COMMA       = ','
    ASSIGN      = ':='
    LESS        = '<'
    LESS_EQUAL      = '<='
    GREATER         = '>'
    GREATER_EQUAL   = '>='
    BANG            = '!'
    EQUAL           = '='
    BANG_EQUAL      = '!='
    # block of reserved words
    PROGRAM     = 'PROGRAM'
    INTEGER     = 'INTEGER'
    REAL        = 'REAL'
    INTEGER_DIV = 'DIV'
    VAR         = 'VAR'
    PROCEDURE   = 'PROCEDURE'
    BEGIN       = 'BEGIN'
    END         = 'END'
    # misc
    IDENT           = 'IDENT'
    INTEGER_CONST   = 'INTEGER_CONST'
    REAL_CONST      = 'REAL_CONST'
    EOF             = 'EOF'

class Token(object):
    def __init__(self, type, value, line = 0, col = 0):
        self.type = type
        self.value = value
        self.line = line
        self.col = col

class Lexer(object):
    def __init__(self, text):
        self.text = text
        self.pos = 0
        self.current_char = self.text[self.pos]
        self.lineno = 1
        self.column = 1
        self.reserved_keywords = {
            'PROGRAM': TokenType.PROGRAM,
            'INTEGER': TokenType.INTEGER,
            'REAL': TokenType.REAL,
            'DIV': TokenType.INTEGER_DIV,
            'VAR': TokenType.VAR,
            'PROCEDURE': TokenType.PROCEDURE,
            'BEGIN': TokenType.BEGIN,
            'END': TokenType.END,
        }

    def error(self):
        s = "Lexer error on '{lexeme}' line: {lineno} column: {column}".format(
            lexeme=self.current_char,
            lineno=self.lineno,
            column=self.column,
        )
        raise Exception(s)

    def advance(self):
        self.pos += 1
        if self.pos > len(self.text) - 1:
            self.current_char = None
        else:
            self.current_char = self.text[self.pos]
        # new line
        if self.current_char == '\n':
            self.lineno += 1
            self.column = 1
    
    def peek(self):
        peek_pos = self.pos + 1
        if peek_pos > len(self.text) - 1:
            return None
        else:
            return self.text[peek_pos]
    
    def skip_whitespace(self):
        while self.current_char != None and self.current_char.isspace():
            self.advance()
    
    def skip_comments(self):
        while self.current_char != None and self.current_char != '}':
            self.advance()

        if self.current_char == None:
            self.error()
        self.advance() # eat '}'
    
    def number(self):
        result = ''
        while self.current_char != None and self.current_char.isdigit():
            result += self.current_char
            self.advance()
        # decimal number
        if self.current_char == '.' and self.peek().isdigit():
            result += self.current_char
            self.advance() # skip the '.'
            while self.current_char != None and self.current_char.isdigit():
                result += self.current_char
                self.advance()
            return self.new_token(TokenType.REAL_CONST, float(result))
        else:
            return self.new_token(TokenType.INTEGER_CONST, int(result))
    
    def isalpha(self, ch):
        return ch.isalnum() or ch == '_'

    def identifier(self):
        result = ''
        while self.current_char != None and self.current_char.isalpha():
            result += self.current_char
            self.advance()
        token_type = self.reserved_keywords.get(result, TokenType.IDENT)
        return self.new_token(token_type, result)
    
    def new_token(self, type = None, value = None):
        if value == None:
            value = type
        return Token(type, value, self.lineno, self.column)

    def get_next_token(self):
        while self.current_char != None:
            if self.current_char.isspace():
                self.skip_whitespace()
                continue

            if self.current_char == '{':
                self.skip_comments()
                continue

            if self.current_char.isdigit():
                return self.number()
            
            if self.isalpha(self.current_char):
                return self.identifier()

            if self.current_char == '+':
                self.advance()
                return self.new_token(TokenType.PLUS)

            if self.current_char == '-':
                self.advance()
                return self.new_token(TokenType.MINUS)

            if self.current_char == '*':
                self.advance()
                return self.new_token(TokenType.MUL)
            
            if self.current_char == '/':
                self.advance()
                return self.new_token(TokenType.FLOAT_DIV)
            
            if self.current_char == '(':
                self.advance()
                return self.new_token(TokenType.LPAREN)
            
            if self.current_char == ')':
                self.advance()
                return self.new_token(TokenType.RPAREN)

            if self.current_char == ';':
                self.advance()
                return self.new_token(TokenType.SEMI)

            if self.current_char == '.':
                self.advance()
                return self.new_token(TokenType.DOT)

            if self.current_char == ':':
                if self.peek() == '=':
                    self.advance() # eat ':'
                    self.advance() # eat '='
                    return self.new_token(TokenType.ASSIGN)
                else:
                    self.advance()
                    return self.new_token(TokenType.COLON)

            if self.current_char == ',':
                self.advance()
                return self.new_token(TokenType.COMMA)
            
            if self.current_char == '<':
                self.advance()
                if self.current_char == '=':
                    self.advance()
                    return self.new_token(TokenType.LESS_EQUAL)
                else:
                    return self.new_token(TokenType.LESS)
            
            if self.current_char == '>':
                self.advance()
                if self.current_char == '=':
                    self.advance()
                    return self.new_token(TokenType.GREATER_EQUAL)
                else:
                    return self.new_token(TokenType.GREATER)
            
            if self.current_char == '!':
                self.advance()
                if self.current_char == '=':
                    self.advance()
                    return self.new_token(TokenType.BANG_EQUAL)
                else:
                    return self.new_token(TokenType.BANG)

            self.error()

        return Token(TokenType.EOF, None)

def main():
    # while True:
    #     try:
    #         text = input('lexer> ')
    #     except EOFError:
    #         break
    #     if not text:
    #         continue
    text = """\
program Main;

procedure Alpha(a : integer; b : integer);
var x : integer;

   procedure Beta(a : integer; b : integer);
   var x : integer;
   begin
      x := a * 10 + b * 2;
   end;

begin
   x := (a + b ) * 2;

   Beta(5, 10);      { procedure call }
end;

begin { Main }

   Alpha(3 + 5, 7);  { procedure call }

end.  { Main }
        """
    lexer = Lexer(text)
    token = lexer.get_next_token()
    while token.type != TokenType.EOF:
        print_token(token)
        token = lexer.get_next_token()
    print_token(token)
        
def print_token(token):
    print('type: {type} value: {value}'.format(
        type = token.type,
        value = token.value
    ))

if __name__ == '__main__':
    main()
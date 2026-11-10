%{
package internal

import "github.com/nightmarlin/lambda"
%}

%union {
	Name string // Set by T_NAME

	Term        lambda.Term         // Set by <term> & <expression>
	Variable    *lambda.Variable    // Set by <variable>
	Abstraction *lambda.Abstraction // Set by <abstraction>
	Application *lambda.Application // Set by <application>
}

%token T_LEFT T_RIGHT
%token T_DOT T_COMMA
%token T_LAMBDA
%token <Name> T_NAME // UTF-8 letters, numbers, '-' or '_' only

%type <Term>        expression
%type <Term>        term
%type <Variable>    variable
%type <Abstraction> abstraction
%type <Application> application

%%

// entrypoint into parse tree
expression : term { llex.(*lLexerImpl).expression = $1 } ;

term : variable { $$ = $1 }
     | abstraction { $$ = $1 }
		 | application { $$ = $1 }
		 | T_LEFT term T_RIGHT { $$ = $2 }
		 ;

variable : T_NAME  { $$ = lambda.NewVariable($1) };

abstraction : T_LAMBDA variable T_DOT term { $$ = lambda.NewAbstraction($2, $4) } ;

application : T_LEFT term term T_RIGHT { $$ = lambda.NewApplication($2, $3) } ;

%%

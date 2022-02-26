package ast

type StatementBase struct {
	Span
	leadingComments         []CommentGroup
	trailingComments        []CommentGroup
	leadingDetachedComments []CommentGroup
}

type Proto struct {
	statements []TopLevelStatement
}

type TopLevelStatement struct{}

type TopLevelDef struct {
	TopLevelStatement
}

type Syntax struct {
	TopLevelStatement
	StatementBase
	keyword    Keyword
	eq         Token
	quoteOpen  Token
	syntax     Token
	quoteClose Token
	semi       Semi
}

type Import struct {
	TopLevelStatement
	StatementBase
	keyword      Keyword
	weakOrPublic Token
	strLit       StrLit
	semi         Semi
}

type Package struct {
	TopLevelStatement
	StatementBase
	keyword   Keyword
	fullIdent FullIdent
	semi      Semi
}

type Option struct {
	TopLevelStatement
	StatementBase
	keyword    Keyword
	optionName OptionName
	eq         Token
	constant   Constant
	semi       Semi
}

type OptionName struct {
	Span
	optionNameSegmentOrDots []OptionNameSegmentOrDot
}

type OptionNameSegmentOrDot struct{}

type OptionNameSegment struct {
	Span
	bracketOpen  Token
	name         FullIdent
	bracketClose Token
}

/*
begin core/parser/recursive-descent-parser.ts
*/

type Span struct {
	start int
	end   int
}

type Token struct {
	Span
}

/*
end core/parser/recursive-descent-parser.ts
*/

/*
begin core/ast/lexical-elements.ts
*/

type CommentGroup struct {
	Span
	comments []Comment
}

type Comment struct {
	Token
}

type SinglelineComment struct {
	Comment
}

type MultilineComment struct {
	Token
}

type Keyword struct {
	Token
}

type Type struct {
	Span
	identOrDots []IdentOrDot
}

type IdentOrDot struct{ Token }

type FullIdent struct {
	Constant
	Span
	identOrDots []IdentOrDot
}

type ident struct{ Token }

type dot struct{ Token }

type comma struct{ Token }

type Semi struct{ Token }

type IntLit struct {
	Token
}

type SignedIntLit struct {
	Constant
	Span
	sign  Token
	value IntLit
}

type FloatLit struct {
	Token
}

type SignedFloatLit struct {
	Constant
	Span
	sign  Token
	value FloatLit
}

type BoolLit struct {
	Constant
	Token
}

type StrLit struct {
	Constant
	Span
	tokens []Token
}

type Aggregate struct {
	Constant
	Span
}

type Empty struct {
	TopLevelStatement
	StatementBase
	semi Semi
}

type Constant struct{}

/*
end core/ast/lexical-elements.ts
*/

/*
begin core/ast/top-level-definitions.ts
*/

type Enum struct {
	StatementBase
	keyword  Keyword
	enumName Token
	enumBody EnumBody
}

type EnumBody struct {
	Span
	bracketOpen  Token
	statements   []EnumBodyStatement
	bracketClose Token
}

type EnumBodyStatement struct{}

// type EnumBodyStatement =
//   | Option
//   | Reserved
//   | EnumField
//   | Empty

type EnumField struct {
	StatementBase
	fieldName    Token
	eq           Token
	fieldNumber  SignedIntLit
	fieldOptions FieldOptions
	semi         Semi
}

type Message struct {
	StatementBase
	keyword     Keyword
	messageName Token
	messageBody MessageBody
}

type MessageBody struct {
	Span
	bracketOpen  Token
	statements   []MessageBodyStatement
	bracketClose Token
}

type MessageBodyStatement struct{}

// type MessageBodyStatement =
//   | Field
//   | MalformedField
//   | Enum
//   | Message
//   | Extend
//   | Extensions
//   | Group
//   | Option
//   | Oneof
//   | MapField
//   | Reserved
//   | Empty

type Extend struct {
	StatementBase
	keyword    Keyword
	extendBody ExtendBody
}

type ExtendBody struct {
	Span
	acketOpen    Token
	nts          []ExtendBodyStatement
	bracketClose Token
}

type ExtendBodyStatement struct{}

// type ExtendBodyStatement =
//   | Field
//   | MalformedField
//   | Group
//   | Empty

type Service struct {
	StatementBase
	keyword     Keyword
	serviceName Token
	serviceBody ServiceBody
}

type ServiceBody struct {
	Span
	bracketOpen  Token
	statements   []ServiceBodyStatement
	bracketClose Token
}

type ServiceBodyStatement struct{}

// type ServiceBodyStatement = Option | Rpc | Empty

type Rpc struct {
	StatementBase
	keyword Keyword
	rpcName Token
	reqType RpcType
	returns Token
	resType RpcType
	//   semiOrRpcBody Semi | RpcBody
}

type RpcBody struct {
	Span
	bracketOpen  Token
	statements   []RpcBodyStatement
	bracketClose Token
}

type RpcBodyStatement struct{}

// type RpcBodyStatement =
//   | Option
//   | Empty

type RpcType struct {
	Span
	bracketOpen  Token
	stream       Keyword
	messageType  Type
	bracketClose Token
}

/*
end core/ast/top-level-definitions.ts
*/

/*
begin core/ast/fields.ts
*/

// type Node =
//   | Field
//   | MalformedField
//   | FieldOptions
//   | FieldOption
//   | Group
//   | Oneof
//   | OneofBody
//   | OneofField
//   | OneofGroup
//   | MapField

type Field struct {
	StatementBase
	fieldLabel   Keyword
	fieldType    Type
	fieldName    Token
	eq           Token
	fieldNumber  IntLit
	fieldOptions FieldOptions
	semi         Semi
}

// type MalformedField = MalformedBase<
//   Field,
//   "malformed-field",
//   | "fieldLabel"
//   | "fieldType"
// >

type FieldOptions struct {
	Span
	bracketOpen Token
	//   fieldOptionOrCommas [](FieldOption | Comma)
	bracketClose Token
}

type FieldOption struct {
	Span
	optionName OptionName
	eq         Token
	constant   Constant
}

type Group struct {
	StatementBase
	groupLabel   Keyword
	keyword      Keyword
	groupName    Token
	eq           Token
	fieldNumber  IntLit
	fieldOptions FieldOptions
	messageBody  MessageBody
}

type Oneof struct {
	StatementBase
	keyword   Keyword
	oneofName Token
	oneofBody OneofBody
}

type OneofBody struct {
	Span
	bracketOpen  Token
	statements   []OneofBodyStatement
	bracketClose Token
}

type OneofBodyStatement struct{}

// type OneofBodyStatement =
//   | Option
//   | OneofField
//   | OneofGroup
//   | Empty

type OneofField struct {
	StatementBase
	fieldType    Type
	fieldName    Token
	eq           Token
	fieldNumber  IntLit
	fieldOptions FieldOptions
	semi         Semi
}

type OneofGroup struct {
	StatementBase
	keyword     Keyword
	groupName   Token
	eq          Token
	fieldNumber IntLit
	messageBody MessageBody
}

type MapField struct {
	StatementBase
	keyword          Keyword
	typeBracketOpen  Token
	keyType          Type
	typeSep          Token
	valueType        Type
	typeBracketClose Token
	mapName          Token
	eq               Token
	fieldNumber      IntLit
	fieldOptions     FieldOptions
	semi             Semi
}

/*
end core/ast/fields.ts
*/

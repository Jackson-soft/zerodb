//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:2
package sqlparser

import __yyfmt__ "fmt"

//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:2

import "bytes"

func SetParseTree(yylex interface{}, stmt Statement) {
	yylex.(*Tokenizer).ParseTree = stmt
}

func SetAllowComments(yylex interface{}, allow bool) {
	yylex.(*Tokenizer).AllowComments = allow
}

func ForceEOF(yylex interface{}) {
	yylex.(*Tokenizer).ForceEOF = true
}

var (
	SHARE        = []byte("share")
	MODE         = []byte("mode")
	IF_BYTES     = []byte("if")
	VALUES_BYTES = []byte("values")
)

//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:27
type yySymType struct {
	yys         int
	empty       struct{}
	statement   Statement
	selStmt     SelectStatement
	byt         byte
	bytes       []byte
	bytes2      [][]byte
	str         string
	selectExprs SelectExprs
	selectExpr  SelectExpr
	columns     Columns
	colName     *ColName
	tableExprs  TableExprs
	tableExpr   TableExpr
	smTableExpr SimpleTableExpr
	tableName   *TableName
	indexHints  *IndexHints
	expr        Expr
	boolExpr    BoolExpr
	valExpr     ValExpr
	tuple       Tuple
	valExprs    ValExprs
	values      Values
	subquery    *Subquery
	caseExpr    *CaseExpr
	whens       []*When
	when        *When
	orderBy     OrderBy
	order       *Order
	limit       *Limit
	insRows     InsertRows
	updateExprs UpdateExprs
	updateExpr  *UpdateExpr
}

const LEX_ERROR = 57346
const SELECT = 57347
const INSERT = 57348
const UPDATE = 57349
const DELETE = 57350
const FROM = 57351
const WHERE = 57352
const GROUP = 57353
const HAVING = 57354
const ORDER = 57355
const BY = 57356
const LIMIT = 57357
const FOR = 57358
const ALL = 57359
const DISTINCT = 57360
const AS = 57361
const EXISTS = 57362
const NULL = 57363
const ASC = 57364
const DESC = 57365
const VALUES = 57366
const INTO = 57367
const DUPLICATE = 57368
const KEY = 57369
const DEFAULT = 57370
const SET = 57371
const LOCK = 57372
const ID = 57373
const STRING = 57374
const NUMBER = 57375
const VALUE_ARG = 57376
const COMMENT = 57377
const UNION = 57378
const MINUS = 57379
const EXCEPT = 57380
const INTERSECT = 57381
const JOIN = 57382
const STRAIGHT_JOIN = 57383
const LEFT = 57384
const RIGHT = 57385
const INNER = 57386
const OUTER = 57387
const CROSS = 57388
const NATURAL = 57389
const USE = 57390
const FORCE = 57391
const ON = 57392
const OR = 57393
const AND = 57394
const NOT = 57395
const BETWEEN = 57396
const CASE = 57397
const WHEN = 57398
const THEN = 57399
const ELSE = 57400
const LE = 57401
const GE = 57402
const NE = 57403
const NULL_SAFE_EQUAL = 57404
const IS = 57405
const LIKE = 57406
const IN = 57407
const UNARY = 57408
const END = 57409
const BEGIN = 57410
const START = 57411
const TRANSACTION = 57412
const COMMIT = 57413
const ROLLBACK = 57414
const NAMES = 57415
const REPLACE = 57416
const ADMIN = 57417
const HELP = 57418
const OFFSET = 57419
const COLLATE = 57420
const CREATE = 57421
const ALTER = 57422
const DROP = 57423
const RENAME = 57424
const TABLE = 57425
const DATABASE = 57426
const INDEX = 57427
const VIEW = 57428
const TO = 57429
const IGNORE = 57430
const IF = 57431
const UNIQUE = 57432
const USING = 57433
const TRUNCATE = 57434
const SHOW = 57435
const TABLES = 57436
const DATABASES = 57437
const COBAR_CLUSTER = 57438
const ROUTE = 57439
const EXPLAIN = 57440

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"LEX_ERROR",
	"SELECT",
	"INSERT",
	"UPDATE",
	"DELETE",
	"FROM",
	"WHERE",
	"GROUP",
	"HAVING",
	"ORDER",
	"BY",
	"LIMIT",
	"FOR",
	"ALL",
	"DISTINCT",
	"AS",
	"EXISTS",
	"NULL",
	"ASC",
	"DESC",
	"VALUES",
	"INTO",
	"DUPLICATE",
	"KEY",
	"DEFAULT",
	"SET",
	"LOCK",
	"ID",
	"STRING",
	"NUMBER",
	"VALUE_ARG",
	"COMMENT",
	"'('",
	"'~'",
	"UNION",
	"MINUS",
	"EXCEPT",
	"INTERSECT",
	"','",
	"JOIN",
	"STRAIGHT_JOIN",
	"LEFT",
	"RIGHT",
	"INNER",
	"OUTER",
	"CROSS",
	"NATURAL",
	"USE",
	"FORCE",
	"ON",
	"OR",
	"AND",
	"NOT",
	"BETWEEN",
	"CASE",
	"WHEN",
	"THEN",
	"ELSE",
	"'='",
	"'<'",
	"'>'",
	"LE",
	"GE",
	"NE",
	"NULL_SAFE_EQUAL",
	"IS",
	"LIKE",
	"IN",
	"'|'",
	"'&'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"'^'",
	"'.'",
	"UNARY",
	"END",
	"BEGIN",
	"START",
	"TRANSACTION",
	"COMMIT",
	"ROLLBACK",
	"NAMES",
	"REPLACE",
	"ADMIN",
	"HELP",
	"OFFSET",
	"COLLATE",
	"CREATE",
	"ALTER",
	"DROP",
	"RENAME",
	"TABLE",
	"DATABASE",
	"INDEX",
	"VIEW",
	"TO",
	"IGNORE",
	"IF",
	"UNIQUE",
	"USING",
	"TRUNCATE",
	"SHOW",
	"TABLES",
	"DATABASES",
	"COBAR_CLUSTER",
	"ROUTE",
	"EXPLAIN",
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 229
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 693

var yyAct = [...]int{

	129, 126, 336, 373, 207, 406, 120, 367, 87, 250,
	291, 137, 89, 205, 208, 3, 168, 77, 286, 219,
	231, 116, 115, 107, 102, 68, 127, 53, 56, 283,
	55, 132, 136, 415, 57, 142, 415, 94, 71, 355,
	357, 415, 119, 133, 134, 135, 176, 124, 140, 176,
	60, 182, 181, 91, 90, 176, 96, 79, 300, 78,
	99, 43, 44, 45, 46, 104, 161, 123, 279, 143,
	384, 248, 74, 75, 76, 383, 85, 62, 63, 64,
	65, 59, 382, 60, 95, 138, 139, 117, 111, 92,
	98, 356, 112, 113, 114, 121, 156, 308, 309, 310,
	311, 312, 150, 313, 314, 417, 166, 61, 416, 91,
	172, 277, 147, 414, 157, 141, 66, 160, 363, 178,
	364, 328, 170, 319, 295, 281, 70, 326, 180, 153,
	204, 206, 167, 109, 191, 420, 174, 280, 209, 264,
	278, 152, 210, 247, 213, 287, 88, 91, 90, 91,
	90, 181, 263, 262, 226, 218, 155, 216, 101, 148,
	238, 221, 381, 182, 181, 223, 224, 227, 359, 366,
	368, 217, 228, 93, 379, 236, 232, 287, 239, 331,
	256, 226, 368, 388, 242, 296, 69, 254, 258, 259,
	121, 243, 194, 195, 196, 191, 246, 260, 255, 165,
	265, 266, 267, 269, 270, 271, 272, 273, 274, 275,
	276, 261, 257, 190, 189, 192, 193, 194, 195, 196,
	191, 380, 103, 58, 105, 121, 121, 293, 182, 181,
	97, 353, 297, 389, 282, 284, 294, 268, 290, 352,
	351, 288, 235, 237, 234, 91, 90, 349, 220, 91,
	304, 298, 350, 302, 347, 148, 279, 390, 170, 348,
	301, 303, 318, 220, 254, 43, 44, 45, 46, 305,
	400, 132, 136, 84, 253, 142, 399, 321, 322, 252,
	306, 398, 119, 133, 134, 135, 151, 124, 140, 320,
	91, 90, 214, 325, 332, 148, 339, 121, 334, 175,
	212, 335, 22, 170, 330, 333, 327, 123, 211, 143,
	192, 193, 194, 195, 196, 191, 345, 346, 254, 254,
	108, 144, 289, 317, 244, 138, 139, 117, 253, 108,
	342, 108, 176, 252, 179, 316, 72, 362, 370, 92,
	412, 360, 369, 358, 341, 365, 72, 22, 23, 24,
	25, 371, 374, 340, 413, 141, 241, 240, 375, 227,
	190, 189, 192, 193, 194, 195, 196, 191, 70, 173,
	164, 26, 162, 146, 159, 385, 158, 154, 100, 401,
	386, 222, 387, 396, 395, 145, 397, 394, 22, 106,
	82, 324, 223, 37, 404, 229, 163, 405, 80, 407,
	407, 407, 402, 403, 374, 408, 409, 169, 337, 378,
	338, 91, 90, 292, 377, 344, 421, 220, 86, 418,
	419, 422, 410, 423, 22, 31, 32, 48, 33, 34,
	21, 35, 36, 20, 19, 18, 27, 28, 30, 29,
	17, 22, 189, 192, 193, 194, 195, 196, 191, 38,
	39, 16, 392, 393, 40, 41, 132, 136, 361, 15,
	142, 14, 13, 12, 110, 230, 54, 92, 133, 134,
	135, 299, 124, 140, 233, 190, 189, 192, 193, 194,
	195, 196, 191, 171, 132, 136, 411, 391, 142, 372,
	376, 343, 123, 329, 143, 92, 133, 134, 135, 215,
	124, 140, 190, 189, 192, 193, 194, 195, 196, 191,
	138, 139, 285, 136, 131, 128, 142, 130, 245, 125,
	123, 183, 143, 92, 133, 134, 135, 122, 151, 140,
	22, 354, 251, 307, 249, 118, 315, 177, 138, 139,
	141, 81, 42, 83, 11, 225, 136, 10, 9, 142,
	143, 8, 7, 6, 5, 4, 92, 133, 134, 135,
	2, 151, 140, 1, 0, 0, 138, 139, 141, 136,
	323, 0, 142, 0, 0, 0, 149, 0, 0, 92,
	133, 134, 135, 143, 151, 140, 0, 190, 189, 192,
	193, 194, 195, 196, 191, 0, 141, 0, 0, 138,
	139, 136, 0, 0, 142, 0, 143, 0, 0, 0,
	0, 92, 133, 134, 135, 0, 151, 140, 0, 0,
	0, 0, 138, 139, 0, 0, 0, 0, 0, 141,
	190, 189, 192, 193, 194, 195, 196, 191, 143, 308,
	309, 310, 311, 312, 0, 313, 314, 0, 0, 0,
	185, 187, 141, 47, 138, 139, 197, 198, 199, 200,
	201, 202, 203, 188, 186, 184, 190, 189, 192, 193,
	194, 195, 196, 191, 0, 0, 0, 49, 50, 51,
	52, 0, 0, 0, 141, 0, 0, 0, 0, 67,
	0, 0, 73,
}
var yyPact = [...]int{

	342, -1000, -1000, 227, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -71, -20, 9,
	-21, -1000, 31, -1000, -1000, -1000, 95, 305, -1000, -37,
	-1000, -1000, 419, 381, -1000, -1000, -1000, 372, -1000, -53,
	337, 409, 58, -67, -16, 305, -67, -1000, -8, 305,
	-1000, 347, -80, -80, 305, -80, -1000, 364, 284, -1000,
	53, -1000, -1000, -10, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 251, -1000, 286, 360, 344, 337, 213, 548, -1000,
	79, -1000, 49, 346, 100, 305, -1000, 345, 343, -1000,
	-36, 341, 376, 339, 146, 305, 337, 383, 308, 338,
	337, -1000, -1000, -1000, -1000, 290, -1000, -1000, 315, 48,
	174, 594, -1000, 464, 436, -1000, -1000, -1000, 580, 272,
	264, -1000, 256, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 580, -1000, 337, 308, 407, 308, -1000,
	288, 525, 492, 328, -1000, 375, -86, -1000, -1000, 147,
	-1000, 326, -1000, -1000, -1000, 325, -1000, 295, -1000, 250,
	227, 29, -1000, -1000, -1000, 243, 251, -1000, -1000, 305,
	136, 464, 464, 580, 250, 82, 580, 580, 181, 580,
	580, 580, 580, 580, 580, 580, 580, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 594, -3, 26, 23, 594,
	-1000, 11, 251, -1000, 419, 86, 558, 293, 253, 400,
	464, -1000, 580, 558, 558, -1000, -1000, 44, -1000, -1000,
	132, 305, -1000, -1000, -44, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 383, 308, 219, -1000, -1000, 308, 238,
	596, 304, 297, 43, -1000, -1000, -1000, -1000, -1000, 96,
	558, -1000, 250, 580, 580, 558, 515, -1000, 370, 236,
	369, -1000, 116, 116, 55, 55, 55, -1000, -1000, 580,
	-1000, -1000, 13, 251, 7, 118, -1000, 464, 383, 308,
	400, 393, 396, 174, 558, 305, 322, -1000, -1000, 313,
	-1000, -1000, 213, 250, -1000, 404, 243, 243, -1000, -1000,
	211, 204, 197, 196, 188, -12, -1000, 312, 54, 310,
	-1000, 558, 403, 580, -1000, 558, -1000, 4, -1000, 38,
	-1000, 580, 109, 129, 117, 393, -1000, 580, 580, -1000,
	-1000, -1000, -1000, 402, 395, 596, 121, -1000, 178, -1000,
	119, -1000, -1000, -1000, -1000, -18, -25, -30, -1000, -1000,
	-1000, 580, 558, -1000, -1000, 558, 580, -1000, 356, -1000,
	-1000, 141, 215, -1000, 430, -1000, 400, 464, 580, 464,
	-1000, -1000, 245, 240, 234, 558, 558, 352, 580, 580,
	580, -1000, -1000, -1000, 393, 174, 214, 174, 305, 305,
	305, 415, 558, 558, -1000, 324, -1, -1000, -6, -9,
	308, -1000, 413, 64, -1000, 305, -1000, -1000, 213, -1000,
	305, -1000, 305, -1000,
}
var yyPgo = [...]int{

	0, 563, 560, 14, 555, 554, 553, 552, 551, 548,
	547, 544, 653, 543, 542, 541, 223, 22, 21, 537,
	536, 535, 534, 9, 533, 532, 25, 531, 5, 19,
	6, 527, 521, 16, 519, 13, 26, 4, 518, 517,
	11, 515, 1, 514, 512, 18, 499, 493, 491, 490,
	10, 489, 3, 487, 2, 486, 23, 483, 7, 8,
	12, 158, 173, 474, 471, 466, 465, 0, 17, 464,
	463, 462, 461, 459, 451, 440, 435, 434, 433, 430,
	427,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 3, 3, 3, 4, 4, 73, 73, 5, 6,
	7, 7, 7, 7, 70, 70, 71, 72, 74, 74,
	75, 76, 8, 8, 8, 8, 77, 77, 77, 78,
	79, 9, 9, 9, 10, 11, 11, 11, 11, 80,
	12, 13, 13, 14, 14, 14, 14, 14, 15, 15,
	17, 17, 18, 18, 18, 21, 21, 19, 19, 19,
	22, 22, 23, 23, 23, 23, 20, 20, 20, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 25, 25,
	25, 26, 26, 27, 27, 27, 27, 28, 28, 29,
	29, 30, 30, 30, 30, 30, 31, 31, 31, 31,
	31, 31, 31, 31, 31, 31, 32, 32, 32, 32,
	32, 32, 32, 33, 33, 38, 38, 36, 36, 40,
	37, 37, 35, 35, 35, 35, 35, 35, 35, 35,
	35, 35, 35, 35, 35, 35, 35, 35, 35, 39,
	39, 41, 41, 41, 43, 46, 46, 44, 44, 45,
	47, 47, 42, 42, 42, 34, 34, 34, 34, 48,
	48, 49, 49, 50, 50, 51, 51, 52, 53, 53,
	53, 54, 54, 54, 54, 55, 55, 55, 56, 56,
	57, 57, 58, 58, 59, 59, 60, 60, 61, 61,
	62, 62, 16, 16, 63, 63, 63, 63, 63, 64,
	64, 65, 65, 66, 66, 67, 68, 69, 69,
}
var yyR2 = [...]int{

	0, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 4, 12, 3, 8, 8, 6, 6, 8, 7,
	3, 4, 4, 6, 1, 2, 1, 1, 4, 2,
	2, 4, 5, 8, 4, 5, 3, 3, 3, 2,
	2, 6, 7, 4, 5, 4, 4, 5, 5, 0,
	2, 0, 2, 1, 2, 1, 1, 1, 0, 1,
	1, 3, 1, 2, 3, 1, 1, 0, 1, 2,
	1, 3, 3, 3, 3, 5, 0, 1, 2, 1,
	1, 2, 3, 2, 3, 2, 2, 2, 1, 3,
	1, 1, 3, 0, 5, 5, 5, 1, 3, 0,
	2, 1, 3, 3, 2, 3, 3, 3, 4, 3,
	4, 5, 6, 3, 4, 2, 1, 1, 1, 1,
	1, 1, 1, 2, 1, 1, 3, 3, 1, 3,
	1, 3, 1, 1, 1, 3, 3, 3, 3, 3,
	3, 3, 3, 2, 3, 4, 5, 4, 1, 1,
	1, 1, 1, 1, 5, 0, 1, 1, 2, 4,
	0, 2, 1, 3, 5, 1, 1, 1, 1, 0,
	3, 0, 2, 0, 3, 1, 3, 2, 0, 1,
	1, 0, 2, 4, 4, 0, 2, 4, 0, 3,
	1, 3, 0, 5, 1, 3, 3, 3, 0, 2,
	0, 3, 0, 1, 1, 1, 1, 1, 1, 0,
	1, 0, 1, 0, 2, 1, 0, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -8, -9,
	-10, -11, -70, -71, -72, -73, -74, -75, -76, -77,
	-78, -79, 5, 6, 7, 8, 29, 94, 95, 97,
	96, 83, 84, 86, 87, 89, 90, 51, 107, 108,
	112, 113, -14, 38, 39, 40, 41, -12, -80, -12,
	-12, -12, -12, 98, -65, 101, 99, 105, -16, 101,
	103, 98, 98, 99, 100, 101, 85, -12, -26, 91,
	31, -67, 31, -12, 109, 110, 111, -68, -68, -3,
	17, -15, 18, -13, -16, -26, 9, -59, 88, -60,
	-42, -67, 31, -62, 104, 100, -67, -62, 98, -67,
	31, -61, 104, -61, -67, -61, 25, -56, 36, 80,
	-69, 98, -68, -68, -68, -17, -18, 76, -21, 31,
	-30, -35, -31, 56, 36, -34, -42, -36, -41, -67,
	-39, -43, 20, 32, 33, 34, 21, -40, 74, 75,
	37, 104, 24, 58, 35, 25, 29, -26, 42, 28,
	-35, 36, 62, 80, 31, 56, -67, -68, 31, 31,
	-68, 102, 31, 20, 31, 53, -67, -26, -33, 24,
	-3, -57, -42, 31, -26, 9, 42, -19, -67, 19,
	80, 55, 54, -32, 71, 56, 70, 57, 69, 73,
	72, 79, 74, 75, 76, 77, 78, 62, 63, 64,
	65, 66, 67, 68, -30, -35, -30, -37, -3, -35,
	-35, 36, 36, -40, 36, -46, -35, -26, -59, -29,
	10, -60, 93, -35, -35, 53, -67, 31, -68, 20,
	-66, 106, -68, -63, 97, 95, 28, 96, 13, 31,
	31, 31, -68, -56, 29, -38, -36, 114, 42, -22,
	-23, -25, 36, 31, -40, -18, -67, 76, -30, -30,
	-35, -36, 71, 70, 57, -35, -35, 21, 56, -35,
	-35, -35, -35, -35, -35, -35, -35, 114, 114, 42,
	114, 114, -17, 18, -17, -44, -45, 59, -56, 29,
	-29, -50, 13, -30, -35, 80, 53, -67, -68, -64,
	102, -33, -59, 42, -42, -29, 42, -24, 43, 44,
	45, 46, 47, 49, 50, -20, 31, 19, -23, 80,
	-36, -35, -35, 55, 21, -35, 114, -17, 114, -47,
	-45, 61, -30, -33, -59, -50, -54, 15, 14, -67,
	31, 31, -36, -48, 11, -23, -23, 43, 48, 43,
	48, 43, 43, 43, -27, 51, 103, 52, 31, 114,
	31, 55, -35, 114, 82, -35, 60, -58, 53, -58,
	-54, -35, -51, -52, -35, -68, -49, 12, 14, 53,
	43, 43, 100, 100, 100, -35, -35, 26, 42, 92,
	42, -53, 22, 23, -50, -30, -37, -30, 36, 36,
	36, 27, -35, -35, -52, -54, -28, -67, -28, -28,
	7, -55, 16, 30, 114, 42, 114, 114, -59, 7,
	71, -67, -67, -67,
}
var yyDef = [...]int{

	0, -2, 1, 2, 3, 4, 5, 6, 7, 8,
	9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
	19, 20, 59, 59, 59, 59, 59, 221, 212, 0,
	0, 34, 0, 36, 37, 59, 0, 0, 59, 0,
	226, 226, 0, 63, 65, 66, 67, 68, 61, 212,
	0, 0, 0, 210, 0, 0, 210, 222, 0, 0,
	213, 0, 208, 208, 0, 208, 35, 0, 198, 39,
	101, 40, 225, 227, 226, 226, 226, 49, 50, 23,
	64, 0, 69, 60, 0, 0, 0, 30, 0, 204,
	0, 172, 225, 0, 0, 0, 226, 0, 0, 226,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 228, 46, 47, 48, 21, 70, 72, 77, 225,
	75, 76, 111, 0, 0, 142, 143, 144, 0, 172,
	0, 158, 0, 175, 176, 177, 178, 138, 161, 162,
	163, 159, 160, 165, 62, 0, 0, 109, 0, 31,
	32, 0, 0, 0, 226, 0, 223, 44, 226, 0,
	53, 0, 55, 209, 56, 0, 226, 198, 38, 0,
	134, 0, 200, 102, 41, 0, 0, 73, 78, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 126, 127, 128,
	129, 130, 131, 132, 114, 0, 0, 0, 0, 140,
	153, 0, 0, 125, 0, 0, 166, 198, 109, 183,
	0, 205, 0, 140, 206, 207, 173, 225, 42, 211,
	0, 0, 45, 226, 219, 214, 215, 216, 217, 218,
	54, 57, 58, 0, 0, 133, 135, 199, 0, 109,
	80, 86, 0, 98, 100, 71, 79, 74, 112, 113,
	116, 117, 0, 0, 0, 119, 0, 123, 0, 145,
	146, 147, 148, 149, 150, 151, 152, 115, 137, 0,
	139, 154, 0, 0, 0, 170, 167, 0, 0, 0,
	183, 191, 0, 110, 33, 0, 0, 224, 51, 0,
	220, 26, 27, 0, 201, 179, 0, 0, 89, 90,
	0, 0, 0, 0, 0, 103, 87, 0, 0, 0,
	118, 120, 0, 0, 124, 141, 155, 0, 157, 0,
	168, 0, 0, 202, 202, 191, 29, 0, 0, 174,
	226, 52, 136, 181, 0, 81, 84, 91, 0, 93,
	0, 95, 96, 97, 82, 0, 0, 0, 88, 83,
	99, 0, 121, 156, 164, 171, 0, 24, 0, 25,
	28, 192, 184, 185, 188, 43, 183, 0, 0, 0,
	92, 94, 0, 0, 0, 122, 169, 0, 0, 0,
	0, 187, 189, 190, 191, 182, 180, 85, 0, 0,
	0, 0, 193, 194, 186, 195, 0, 107, 0, 0,
	0, 22, 0, 0, 104, 0, 105, 106, 203, 196,
	0, 108, 0, 197,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 78, 73, 3,
	36, 114, 76, 74, 42, 75, 80, 77, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	63, 62, 64, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 79, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 72, 3, 37,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 38, 39, 40, 41, 43, 44,
	45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
	55, 56, 57, 58, 59, 60, 61, 65, 66, 67,
	68, 69, 70, 71, 81, 82, 83, 84, 85, 86,
	87, 88, 89, 90, 91, 92, 93, 94, 95, 96,
	97, 98, 99, 100, 101, 102, 103, 104, 105, 106,
	107, 108, 109, 110, 111, 112, 113,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*    parser for yacc output    */

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:188
		{
			SetParseTree(yylex, yyDollar[1].statement)
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:194
		{
			yyVAL.statement = yyDollar[1].selStmt
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:218
		{
			yyVAL.selStmt = &SimpleSelect{Comments: Comments(yyDollar[2].bytes2), Distinct: yyDollar[3].str, SelectExprs: yyDollar[4].selectExprs}
		}
	case 22:
		yyDollar = yyS[yypt-12 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:222
		{
			yyVAL.selStmt = &Select{Comments: Comments(yyDollar[2].bytes2), Distinct: yyDollar[3].str, SelectExprs: yyDollar[4].selectExprs, From: yyDollar[6].tableExprs, Where: NewWhere(AST_WHERE, yyDollar[7].boolExpr), GroupBy: GroupBy(yyDollar[8].valExprs), Having: NewWhere(AST_HAVING, yyDollar[9].boolExpr), OrderBy: yyDollar[10].orderBy, Limit: yyDollar[11].limit, Lock: yyDollar[12].str}
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:226
		{
			yyVAL.selStmt = &Union{Type: yyDollar[2].str, Left: yyDollar[1].selStmt, Right: yyDollar[3].selStmt}
		}
	case 24:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:233
		{
			yyVAL.statement = &Insert{Comments: Comments(yyDollar[2].bytes2), Ignore: yyDollar[3].str, Table: yyDollar[5].tableName, Columns: yyDollar[6].columns, Rows: yyDollar[7].insRows, OnDup: OnDup(yyDollar[8].updateExprs)}
		}
	case 25:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:237
		{
			cols := make(Columns, 0, len(yyDollar[7].updateExprs))
			vals := make(ValTuple, 0, len(yyDollar[7].updateExprs))
			for _, col := range yyDollar[7].updateExprs {
				cols = append(cols, &NonStarExpr{Expr: col.Name})
				vals = append(vals, col.Expr)
			}
			yyVAL.statement = &Insert{Comments: Comments(yyDollar[2].bytes2), Ignore: yyDollar[3].str, Table: yyDollar[5].tableName, Columns: cols, Rows: Values{vals}, OnDup: OnDup(yyDollar[8].updateExprs)}
		}
	case 26:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:249
		{
			yyVAL.statement = &Replace{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Columns: yyDollar[5].columns, Rows: yyDollar[6].insRows}
		}
	case 27:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:253
		{
			cols := make(Columns, 0, len(yyDollar[6].updateExprs))
			vals := make(ValTuple, 0, len(yyDollar[6].updateExprs))
			for _, col := range yyDollar[6].updateExprs {
				cols = append(cols, &NonStarExpr{Expr: col.Name})
				vals = append(vals, col.Expr)
			}
			yyVAL.statement = &Replace{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Columns: cols, Rows: Values{vals}}
		}
	case 28:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:266
		{
			yyVAL.statement = &Update{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[3].tableName, Exprs: yyDollar[5].updateExprs, Where: NewWhere(AST_WHERE, yyDollar[6].boolExpr), OrderBy: yyDollar[7].orderBy, Limit: yyDollar[8].limit}
		}
	case 29:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:272
		{
			yyVAL.statement = &Delete{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Where: NewWhere(AST_WHERE, yyDollar[5].boolExpr), OrderBy: yyDollar[6].orderBy, Limit: yyDollar[7].limit}
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:278
		{
			yyVAL.statement = &Set{Comments: Comments(yyDollar[2].bytes2), Exprs: yyDollar[3].updateExprs}
		}
	case 31:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:282
		{
			yyVAL.statement = &Set{Comments: Comments(yyDollar[2].bytes2), Exprs: UpdateExprs{&UpdateExpr{Name: &ColName{Name: []byte("names")}, Expr: StrVal("default")}}}
		}
	case 32:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:286
		{
			yyVAL.statement = &Set{Comments: Comments(yyDollar[2].bytes2), Exprs: UpdateExprs{&UpdateExpr{Name: &ColName{Name: []byte("names")}, Expr: yyDollar[4].valExpr}}}
		}
	case 33:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:290
		{
			yyVAL.statement = &Set{
				Comments: Comments(yyDollar[2].bytes2),
				Exprs: UpdateExprs{
					&UpdateExpr{
						Name: &ColName{Name: []byte("names")}, Expr: yyDollar[4].valExpr,
					},
					&UpdateExpr{
						Name: &ColName{Name: []byte("collate")}, Expr: yyDollar[6].valExpr,
					},
				},
			}
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:306
		{
			yyVAL.statement = &Begin{}
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:310
		{
			yyVAL.statement = &Begin{}
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:317
		{
			yyVAL.statement = &Commit{}
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:323
		{
			yyVAL.statement = &Rollback{}
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:329
		{
			yyVAL.statement = &Admin{Region: yyDollar[2].tableName, Columns: yyDollar[3].columns, Rows: yyDollar[4].insRows}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:333
		{
			yyVAL.statement = &AdminHelp{}
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:339
		{
			yyVAL.statement = &UseDB{DB: string(yyDollar[2].bytes)}
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:345
		{
			yyVAL.statement = &Truncate{Comments: Comments(yyDollar[2].bytes2), TableOpt: yyDollar[3].str, Table: yyDollar[4].tableName}
		}
	case 42:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:351
		{
			yyVAL.statement = &DDL{Action: AST_TABLE_CREATE, Table: yyDollar[4].bytes, NewName: yyDollar[4].bytes}
		}
	case 43:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:355
		{
			yyVAL.statement = &DDL{Action: AST_INDEX_CREATE, Table: yyDollar[7].bytes, NewName: yyDollar[4].bytes}
		}
	case 44:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:359
		{
			yyVAL.statement = &DDL{Action: AST_TABLE_CREATE, NewName: yyDollar[3].bytes}
		}
	case 45:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:363
		{
			yyVAL.statement = &DDL{Action: AST_DATABASE_CREATE, Database: yyDollar[4].bytes, NewName: yyDollar[4].bytes}
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:369
		{
			yyVAL.statement = &Show{Type: AST_SHOW_TYPE_TABLES}
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:373
		{
			yyVAL.statement = &Show{Type: AST_SHOW_TYPE_DATABASES}
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:377
		{
			yyVAL.statement = &Show{Type: AST_SHOW_TYPE_COBAR_CLUSTER}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:383
		{
			yyVAL.statement = &Route{}
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:389
		{
			yyVAL.statement = &Explain{}
		}
	case 51:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:395
		{
			yyVAL.statement = &DDL{Action: AST_ALTER, Ignore: yyDollar[2].str, Table: yyDollar[4].bytes, NewName: yyDollar[4].bytes}
		}
	case 52:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:399
		{
			// Change this to a rename statement
			yyVAL.statement = &DDL{Action: AST_RENAME, Ignore: yyDollar[2].str, Table: yyDollar[4].bytes, NewName: yyDollar[7].bytes}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:404
		{
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[3].bytes, NewName: yyDollar[3].bytes}
		}
	case 54:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:410
		{
			yyVAL.statement = &DDL{Action: AST_RENAME, Table: yyDollar[3].bytes, NewName: yyDollar[5].bytes}
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:416
		{
			yyVAL.statement = &DDL{Action: AST_TABLE_DROP, Table: yyDollar[4].bytes}
		}
	case 56:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:420
		{
			yyVAL.statement = &DDL{Action: AST_DATABASE_DROP, Database: yyDollar[4].bytes}
		}
	case 57:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:424
		{
			// Change this to an alter statement
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[5].bytes, NewName: yyDollar[5].bytes}
		}
	case 58:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:429
		{
			yyVAL.statement = &DDL{Action: AST_TABLE_DROP, Table: yyDollar[4].bytes}
		}
	case 59:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:434
		{
			SetAllowComments(yylex, true)
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:438
		{
			yyVAL.bytes2 = yyDollar[2].bytes2
			SetAllowComments(yylex, false)
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:444
		{
			yyVAL.bytes2 = nil
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:448
		{
			yyVAL.bytes2 = append(yyDollar[1].bytes2, yyDollar[2].bytes)
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:454
		{
			yyVAL.str = AST_UNION
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:458
		{
			yyVAL.str = AST_UNION_ALL
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:462
		{
			yyVAL.str = AST_SET_MINUS
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:466
		{
			yyVAL.str = AST_EXCEPT
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:470
		{
			yyVAL.str = AST_INTERSECT
		}
	case 68:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:475
		{
			yyVAL.str = ""
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:479
		{
			yyVAL.str = AST_DISTINCT
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:485
		{
			yyVAL.selectExprs = SelectExprs{yyDollar[1].selectExpr}
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:489
		{
			yyVAL.selectExprs = append(yyVAL.selectExprs, yyDollar[3].selectExpr)
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:495
		{
			yyVAL.selectExpr = &StarExpr{}
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:499
		{
			yyVAL.selectExpr = &NonStarExpr{Expr: yyDollar[1].expr, As: yyDollar[2].bytes}
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:503
		{
			yyVAL.selectExpr = &StarExpr{TableName: yyDollar[1].bytes}
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:509
		{
			yyVAL.expr = yyDollar[1].boolExpr
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:513
		{
			yyVAL.expr = yyDollar[1].valExpr
		}
	case 77:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:518
		{
			yyVAL.bytes = nil
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:522
		{
			yyVAL.bytes = yyDollar[1].bytes
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:526
		{
			yyVAL.bytes = yyDollar[2].bytes
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:532
		{
			yyVAL.tableExprs = TableExprs{yyDollar[1].tableExpr}
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:536
		{
			yyVAL.tableExprs = append(yyVAL.tableExprs, yyDollar[3].tableExpr)
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:542
		{
			yyVAL.tableExpr = &AliasedTableExpr{Expr: yyDollar[1].smTableExpr, As: yyDollar[2].bytes, Hints: yyDollar[3].indexHints}
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:546
		{
			yyVAL.tableExpr = &ParenTableExpr{Expr: yyDollar[2].tableExpr}
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:550
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr}
		}
	case 85:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:554
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr, On: yyDollar[5].boolExpr}
		}
	case 86:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:559
		{
			yyVAL.bytes = nil
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:563
		{
			yyVAL.bytes = yyDollar[1].bytes
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:567
		{
			yyVAL.bytes = yyDollar[2].bytes
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:573
		{
			yyVAL.str = AST_JOIN
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:577
		{
			yyVAL.str = AST_STRAIGHT_JOIN
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:581
		{
			yyVAL.str = AST_LEFT_JOIN
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:585
		{
			yyVAL.str = AST_LEFT_JOIN
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:589
		{
			yyVAL.str = AST_RIGHT_JOIN
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:593
		{
			yyVAL.str = AST_RIGHT_JOIN
		}
	case 95:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:597
		{
			yyVAL.str = AST_JOIN
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:601
		{
			yyVAL.str = AST_CROSS_JOIN
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:605
		{
			yyVAL.str = AST_NATURAL_JOIN
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:611
		{
			yyVAL.smTableExpr = &TableName{Name: yyDollar[1].bytes}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:615
		{
			yyVAL.smTableExpr = &TableName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:619
		{
			yyVAL.smTableExpr = yyDollar[1].subquery
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:625
		{
			yyVAL.tableName = &TableName{Name: yyDollar[1].bytes}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:629
		{
			yyVAL.tableName = &TableName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 103:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:634
		{
			yyVAL.indexHints = nil
		}
	case 104:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:638
		{
			yyVAL.indexHints = &IndexHints{Type: AST_USE, Indexes: yyDollar[4].bytes2}
		}
	case 105:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:642
		{
			yyVAL.indexHints = &IndexHints{Type: AST_IGNORE, Indexes: yyDollar[4].bytes2}
		}
	case 106:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:646
		{
			yyVAL.indexHints = &IndexHints{Type: AST_FORCE, Indexes: yyDollar[4].bytes2}
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:652
		{
			yyVAL.bytes2 = [][]byte{yyDollar[1].bytes}
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:656
		{
			yyVAL.bytes2 = append(yyDollar[1].bytes2, yyDollar[3].bytes)
		}
	case 109:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:661
		{
			yyVAL.boolExpr = nil
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:665
		{
			yyVAL.boolExpr = yyDollar[2].boolExpr
		}
	case 112:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:672
		{
			yyVAL.boolExpr = &AndExpr{Left: yyDollar[1].boolExpr, Right: yyDollar[3].boolExpr}
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:676
		{
			yyVAL.boolExpr = &OrExpr{Left: yyDollar[1].boolExpr, Right: yyDollar[3].boolExpr}
		}
	case 114:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:680
		{
			yyVAL.boolExpr = &NotExpr{Expr: yyDollar[2].boolExpr}
		}
	case 115:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:684
		{
			yyVAL.boolExpr = &ParenBoolExpr{Expr: yyDollar[2].boolExpr}
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:690
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: yyDollar[2].str, Right: yyDollar[3].valExpr}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:694
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_IN, Right: yyDollar[3].tuple}
		}
	case 118:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:698
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_NOT_IN, Right: yyDollar[4].tuple}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:702
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_LIKE, Right: yyDollar[3].valExpr}
		}
	case 120:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:706
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_NOT_LIKE, Right: yyDollar[4].valExpr}
		}
	case 121:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:710
		{
			yyVAL.boolExpr = &RangeCond{Left: yyDollar[1].valExpr, Operator: AST_BETWEEN, From: yyDollar[3].valExpr, To: yyDollar[5].valExpr}
		}
	case 122:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:714
		{
			yyVAL.boolExpr = &RangeCond{Left: yyDollar[1].valExpr, Operator: AST_NOT_BETWEEN, From: yyDollar[4].valExpr, To: yyDollar[6].valExpr}
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:718
		{
			yyVAL.boolExpr = &NullCheck{Operator: AST_IS_NULL, Expr: yyDollar[1].valExpr}
		}
	case 124:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:722
		{
			yyVAL.boolExpr = &NullCheck{Operator: AST_IS_NOT_NULL, Expr: yyDollar[1].valExpr}
		}
	case 125:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:726
		{
			yyVAL.boolExpr = &ExistsExpr{Subquery: yyDollar[2].subquery}
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:732
		{
			yyVAL.str = AST_EQ
		}
	case 127:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:736
		{
			yyVAL.str = AST_LT
		}
	case 128:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:740
		{
			yyVAL.str = AST_GT
		}
	case 129:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:744
		{
			yyVAL.str = AST_LE
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:748
		{
			yyVAL.str = AST_GE
		}
	case 131:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:752
		{
			yyVAL.str = AST_NE
		}
	case 132:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:756
		{
			yyVAL.str = AST_NSE
		}
	case 133:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:762
		{
			yyVAL.insRows = yyDollar[2].values
		}
	case 134:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:766
		{
			yyVAL.insRows = yyDollar[1].selStmt
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:772
		{
			yyVAL.values = Values{yyDollar[1].tuple}
		}
	case 136:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:776
		{
			yyVAL.values = append(yyDollar[1].values, yyDollar[3].tuple)
		}
	case 137:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:782
		{
			yyVAL.tuple = ValTuple(yyDollar[2].valExprs)
		}
	case 138:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:786
		{
			yyVAL.tuple = yyDollar[1].subquery
		}
	case 139:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:792
		{
			yyVAL.subquery = &Subquery{yyDollar[2].selStmt}
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:798
		{
			yyVAL.valExprs = ValExprs{yyDollar[1].valExpr}
		}
	case 141:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:802
		{
			yyVAL.valExprs = append(yyDollar[1].valExprs, yyDollar[3].valExpr)
		}
	case 142:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:808
		{
			yyVAL.valExpr = yyDollar[1].valExpr
		}
	case 143:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:812
		{
			yyVAL.valExpr = yyDollar[1].colName
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:816
		{
			yyVAL.valExpr = yyDollar[1].tuple
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:820
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITAND, Right: yyDollar[3].valExpr}
		}
	case 146:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:824
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITOR, Right: yyDollar[3].valExpr}
		}
	case 147:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:828
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITXOR, Right: yyDollar[3].valExpr}
		}
	case 148:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:832
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_PLUS, Right: yyDollar[3].valExpr}
		}
	case 149:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:836
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MINUS, Right: yyDollar[3].valExpr}
		}
	case 150:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:840
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MULT, Right: yyDollar[3].valExpr}
		}
	case 151:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:844
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_DIV, Right: yyDollar[3].valExpr}
		}
	case 152:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:848
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MOD, Right: yyDollar[3].valExpr}
		}
	case 153:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:852
		{
			if num, ok := yyDollar[2].valExpr.(NumVal); ok {
				switch yyDollar[1].byt {
				case '-':
					yyVAL.valExpr = append(NumVal("-"), num...)
				case '+':
					yyVAL.valExpr = num
				default:
					yyVAL.valExpr = &UnaryExpr{Operator: yyDollar[1].byt, Expr: yyDollar[2].valExpr}
				}
			} else {
				yyVAL.valExpr = &UnaryExpr{Operator: yyDollar[1].byt, Expr: yyDollar[2].valExpr}
			}
		}
	case 154:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:867
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes}
		}
	case 155:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:871
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Exprs: yyDollar[3].selectExprs}
		}
	case 156:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:875
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Distinct: true, Exprs: yyDollar[4].selectExprs}
		}
	case 157:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:879
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Exprs: yyDollar[3].selectExprs}
		}
	case 158:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:883
		{
			yyVAL.valExpr = yyDollar[1].caseExpr
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:889
		{
			yyVAL.bytes = IF_BYTES
		}
	case 160:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:893
		{
			yyVAL.bytes = VALUES_BYTES
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:899
		{
			yyVAL.byt = AST_UPLUS
		}
	case 162:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:903
		{
			yyVAL.byt = AST_UMINUS
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:907
		{
			yyVAL.byt = AST_TILDA
		}
	case 164:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:913
		{
			yyVAL.caseExpr = &CaseExpr{Expr: yyDollar[2].valExpr, Whens: yyDollar[3].whens, Else: yyDollar[4].valExpr}
		}
	case 165:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:918
		{
			yyVAL.valExpr = nil
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:922
		{
			yyVAL.valExpr = yyDollar[1].valExpr
		}
	case 167:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:928
		{
			yyVAL.whens = []*When{yyDollar[1].when}
		}
	case 168:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:932
		{
			yyVAL.whens = append(yyDollar[1].whens, yyDollar[2].when)
		}
	case 169:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:938
		{
			yyVAL.when = &When{Cond: yyDollar[2].boolExpr, Val: yyDollar[4].valExpr}
		}
	case 170:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:943
		{
			yyVAL.valExpr = nil
		}
	case 171:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:947
		{
			yyVAL.valExpr = yyDollar[2].valExpr
		}
	case 172:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:953
		{
			yyVAL.colName = &ColName{Name: yyDollar[1].bytes}
		}
	case 173:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:957
		{
			yyVAL.colName = &ColName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 174:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:961
		{
			yyVAL.colName = &ColName{Qualifier: yyDollar[3].bytes, Name: yyDollar[5].bytes}
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:967
		{
			yyVAL.valExpr = StrVal(yyDollar[1].bytes)
		}
	case 176:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:971
		{
			yyVAL.valExpr = NumVal(yyDollar[1].bytes)
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:975
		{
			yyVAL.valExpr = ValArg(yyDollar[1].bytes)
		}
	case 178:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:979
		{
			yyVAL.valExpr = &NullVal{}
		}
	case 179:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:984
		{
			yyVAL.valExprs = nil
		}
	case 180:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:988
		{
			yyVAL.valExprs = yyDollar[3].valExprs
		}
	case 181:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:993
		{
			yyVAL.boolExpr = nil
		}
	case 182:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:997
		{
			yyVAL.boolExpr = yyDollar[2].boolExpr
		}
	case 183:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1002
		{
			yyVAL.orderBy = nil
		}
	case 184:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1006
		{
			yyVAL.orderBy = yyDollar[3].orderBy
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1012
		{
			yyVAL.orderBy = OrderBy{yyDollar[1].order}
		}
	case 186:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1016
		{
			yyVAL.orderBy = append(yyDollar[1].orderBy, yyDollar[3].order)
		}
	case 187:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1022
		{
			yyVAL.order = &Order{Expr: yyDollar[1].valExpr, Direction: yyDollar[2].str}
		}
	case 188:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1027
		{
			yyVAL.str = AST_ASC
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1031
		{
			yyVAL.str = AST_ASC
		}
	case 190:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1035
		{
			yyVAL.str = AST_DESC
		}
	case 191:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1040
		{
			yyVAL.limit = nil
		}
	case 192:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1044
		{
			yyVAL.limit = &Limit{Rowcount: yyDollar[2].valExpr}
		}
	case 193:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1048
		{
			yyVAL.limit = &Limit{Offset: yyDollar[2].valExpr, Rowcount: yyDollar[4].valExpr}
		}
	case 194:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1052
		{
			yyVAL.limit = &Limit{Offset: yyDollar[4].valExpr, Rowcount: yyDollar[2].valExpr}
		}
	case 195:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1057
		{
			yyVAL.str = ""
		}
	case 196:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1061
		{
			yyVAL.str = AST_FOR_UPDATE
		}
	case 197:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1065
		{
			if !bytes.Equal(yyDollar[3].bytes, SHARE) {
				yylex.Error("expecting share")
				return 1
			}
			if !bytes.Equal(yyDollar[4].bytes, MODE) {
				yylex.Error("expecting mode")
				return 1
			}
			yyVAL.str = AST_SHARE_MODE
		}
	case 198:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1078
		{
			yyVAL.columns = nil
		}
	case 199:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1082
		{
			yyVAL.columns = yyDollar[2].columns
		}
	case 200:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1088
		{
			yyVAL.columns = Columns{&NonStarExpr{Expr: yyDollar[1].colName}}
		}
	case 201:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1092
		{
			yyVAL.columns = append(yyVAL.columns, &NonStarExpr{Expr: yyDollar[3].colName})
		}
	case 202:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1097
		{
			yyVAL.updateExprs = nil
		}
	case 203:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1101
		{
			yyVAL.updateExprs = yyDollar[5].updateExprs
		}
	case 204:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1107
		{
			yyVAL.updateExprs = UpdateExprs{yyDollar[1].updateExpr}
		}
	case 205:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1111
		{
			yyVAL.updateExprs = append(yyDollar[1].updateExprs, yyDollar[3].updateExpr)
		}
	case 206:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1117
		{
			yyVAL.updateExpr = &UpdateExpr{Name: yyDollar[1].colName, Expr: yyDollar[3].valExpr}
		}
	case 207:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1121
		{
			yyVAL.updateExpr = &UpdateExpr{Name: yyDollar[1].colName, Expr: StrVal("ON")}
		}
	case 208:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1126
		{
			yyVAL.empty = struct{}{}
		}
	case 209:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1128
		{
			yyVAL.empty = struct{}{}
		}
	case 210:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1131
		{
			yyVAL.empty = struct{}{}
		}
	case 211:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1133
		{
			yyVAL.empty = struct{}{}
		}
	case 212:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1136
		{
			yyVAL.str = ""
		}
	case 213:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1138
		{
			yyVAL.str = AST_IGNORE
		}
	case 214:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1142
		{
			yyVAL.empty = struct{}{}
		}
	case 215:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1144
		{
			yyVAL.empty = struct{}{}
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1146
		{
			yyVAL.empty = struct{}{}
		}
	case 217:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1148
		{
			yyVAL.empty = struct{}{}
		}
	case 218:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1150
		{
			yyVAL.empty = struct{}{}
		}
	case 219:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1153
		{
			yyVAL.empty = struct{}{}
		}
	case 220:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1155
		{
			yyVAL.empty = struct{}{}
		}
	case 221:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1158
		{
			yyVAL.empty = struct{}{}
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1160
		{
			yyVAL.empty = struct{}{}
		}
	case 223:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1163
		{
			yyVAL.empty = struct{}{}
		}
	case 224:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1165
		{
			yyVAL.empty = struct{}{}
		}
	case 225:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1169
		{
			yyVAL.bytes = bytes.ToLower(yyDollar[1].bytes)
		}
	case 226:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1174
		{
			ForceEOF(yylex)
		}
	case 227:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1179
		{
			yyVAL.str = ""
		}
	case 228:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line /Users/jacksoncy/go/src/git.2dfire.net/zerodb/proxy/proxy/sqlparser/sql.y:1183
		{
			yyVAL.str = AST_TABLE
		}
	}
	goto yystack /* stack new state and value */
}

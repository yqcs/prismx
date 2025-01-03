package temp

import (
	"bytes"
	"encoding/base64"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"math/rand"
	"net/url"
	"prismx_cli/utils/cryptoPlus"
	"prismx_cli/utils/randomUtils"
	"regexp"
	"strings"
)

type CustomLib struct {
	envOptions     []cel.EnvOption
	programOptions []cel.ProgramOption
}

func Evaluate(env *cel.Env, expression string, params map[string]any) (ref.Val, error) {
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}

	out, _, err := prg.Eval(params)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewEnv(c *CustomLib) (*cel.Env, error) {
	return cel.NewEnv(cel.Lib(c))
}

// NewEnvOption 名称限定禁止为exploit，防止影响到exploit模块
func NewEnvOption() CustomLib {
	c := CustomLib{}
	reg := types.NewEmptyRegistry()

	c.envOptions = []cel.EnvOption{
		cel.CustomTypeAdapter(reg),
		cel.CustomTypeProvider(reg),
		cel.Declarations(

			//字符串包含
			decls.NewFunction("noContains",
				decls.NewInstanceOverload("contains_string_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool)),

			//字符串包含 忽略大小写
			decls.NewFunction("iContains",
				decls.NewInstanceOverload("iContains_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool)),

			//字符串截取
			decls.NewFunction("subString",
				decls.NewOverload("substr_string_int_int",
					[]*exprpb.Type{decls.String, decls.Int, decls.Int},
					decls.String)),

			//字符串以开头
			decls.NewFunction("hasPrefix",
				decls.NewOverload("hasPrefix_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool)),

			//字符串以结尾
			decls.NewFunction("hasSuffix",
				decls.NewOverload("hasSuffix_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool)),

			//字符串以结尾
			decls.NewFunction("replace",
				decls.NewOverload("replaceAll",
					[]*exprpb.Type{decls.String, decls.String, decls.String},
					decls.String)),
			//判断是否包含
			decls.NewFunction("bContains",
				decls.NewInstanceOverload("bContains_bytes",
					[]*exprpb.Type{decls.Bytes, decls.Bytes},
					decls.Bool)),
			//判断是否包含
			decls.NewFunction("ibContains",
				decls.NewInstanceOverload("ibContains_bytes",
					[]*exprpb.Type{decls.Bytes, decls.Bytes},
					decls.Bool)),
			//md5转换
			decls.NewFunction("md5",
				decls.NewOverload("md5String",
					[]*exprpb.Type{decls.String},
					decls.String)),
			//base64
			decls.NewFunction("base64String",
				decls.NewOverload("base64_string",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("base64Bytes",
				decls.NewOverload("base64_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),
			decls.NewFunction("base64StringDecode",
				decls.NewOverload("base64Decode_string",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("base64BytesDecode",
				decls.NewOverload("base64DecodeBytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),
			decls.NewFunction("urlEncode",
				decls.NewOverload("urlEncodeString",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("UrlEncode",
				decls.NewOverload("urlEncode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),
			decls.NewFunction("urlDecode",
				decls.NewOverload("urlDecodeString",
					[]*exprpb.Type{decls.String},
					decls.String)),
			decls.NewFunction("UrlDecode",
				decls.NewOverload("urlDecode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String)),

			decls.NewFunction("randomInt",
				decls.NewOverload("randomInt",
					[]*exprpb.Type{decls.Int, decls.Int},
					decls.Int)),
			decls.NewFunction("randomLowercase",
				decls.NewOverload("randomLowercaseInt",
					[]*exprpb.Type{decls.Int},
					decls.String)),
			decls.NewFunction("sMatches",
				decls.NewInstanceOverload("string_bMatches_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool)),
			decls.NewFunction("bMatches",
				decls.NewInstanceOverload("string_bMatches_bytes",
					[]*exprpb.Type{decls.String, decls.Bytes},
					decls.Bool)),
		),
	}

	c.programOptions = []cel.ProgramOption{
		cel.Functions(
			//字符串包含
			&functions.Overload{
				Operator: "contains_string_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to contains", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to contains", rhs.Type())
					}
					// 不区分大小写包含
					return types.Bool(!strings.Contains(strings.ToLower(string(v1)), strings.ToLower(string(v2))))
				},
			},

			&functions.Overload{
				Operator: "iContains_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to iContains", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to iContains", rhs.Type())
					}
					// 不区分大小写包含
					return types.Bool(strings.Contains(string(v1), string(v2)))
				},
			},
			//是否以开头
			&functions.Overload{
				Operator: "hasPrefix_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to hasPrefix", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to hasPrefix", rhs.Type())
					}
					return types.Bool(strings.HasPrefix(string(v1), string(v2)))
				},
			},
			//是否以结尾
			&functions.Overload{
				Operator: "hasSuffix_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to hasSuffix", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to hasSuffix", rhs.Type())
					}
					return types.Bool(strings.HasSuffix(string(v1), string(v2)))
				},
			},
			//替换字符串
			&functions.Overload{
				Operator: "replaceAll",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to replace", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to replace", rhs.Type())
					}
					v3, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to replace", rhs.Type())
					}
					return types.String(strings.ReplaceAll(string(v1), string(v2), string(v3)))
				},
			},

			//判断是否包含
			&functions.Overload{
				Operator: "bContains_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bContains", lhs.Type())
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bContains", rhs.Type())
					}
					return types.Bool(bytes.Contains(v1, v2))
				},
			},
			//转小写判断是否包含
			&functions.Overload{
				Operator: "ibContains_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to ibContains", lhs.Type())
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to ibContains", rhs.Type())
					}
					return types.Bool(bytes.Contains(bytes.ToLower(v1), bytes.ToLower(v2)))
				},
			},
			&functions.Overload{
				Operator: "string_bMatches_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to matches", lhs.Type())
					}
					v2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to matches", rhs.Type())
					}
					res, err := regexp.MatchString(string(v1), string(v2))
					if res && err == nil {
						return types.Bool(true)
					}
					return types.Bool(false)
				},
			},
			&functions.Overload{
				Operator: "string_bMatches_bytes",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					v1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bMatches", lhs.Type())
					}
					v2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bMatches", rhs.Type())
					}
					res, err := regexp.Match(string(v1), v2)
					if res && err == nil {
						return types.Bool(true)
					}
					return types.Bool(false)
				},
			},
			&functions.Overload{
				Operator: "md5String",
				Unary: func(value ref.Val) ref.Val {
					v1, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to md5 String", value.Type())
					}
					return types.String(cryptoPlus.ToMD5(string(v1)))
				},
			},
			&functions.Overload{
				Operator: "randomInt",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					from, ok := lhs.(types.Int)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to randomInt", lhs.Type())
					}
					to, ok := rhs.(types.Int)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to randomInt", rhs.Type())
					}
					min, max := int(from), int(to)
					return types.Int(rand.Intn(max-min) + min)
				},
			},
			&functions.Overload{
				Operator: "randomLowercaseInt",
				Unary: func(value ref.Val) ref.Val {
					n, ok := value.(types.Int)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to randomLowercase", value.Type())
					}
					return types.String(strings.ToLower(randomUtils.RandomString(int(n))))
				},
			},
			&functions.Overload{
				Operator: "base64_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64_string", value.Type())
					}
					return types.String(base64.StdEncoding.EncodeToString([]byte(v)))
				},
			},
			&functions.Overload{
				Operator: "base64_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64_bytes", value.Type())
					}
					return types.String(base64.StdEncoding.EncodeToString(v))
				},
			},
			&functions.Overload{
				Operator: "base64DecodeBytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode Bytes", value.Type())
					}
					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "base64Decode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_string", value.Type())
					}
					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeBytes)
				},
			}, &functions.Overload{
				Operator: "urlEncodeString",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlEncode String", value.Type())
					}
					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urlEncode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_bytes", value.Type())
					}
					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urlDecodeString",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlDecode String", value.Type())
					}
					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "urlDecode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_bytes", value.Type())
					}
					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}
					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "substr_string_int_int",
				Function: func(values ...ref.Val) ref.Val {
					if len(values) == 3 {
						str, ok := values[0].(types.String)
						if !ok {
							return types.NewErr("invalid string to 'substr'")
						}
						start, ok := values[1].(types.Int)
						if !ok {
							return types.NewErr("invalid start to 'substr'")
						}
						end, ok := values[2].(types.Int)
						if !ok {
							return types.NewErr("invalid length to 'substr'")
						}
						runes := []rune(str)
						if start < 0 || end < 0 || int(start+end) > len(runes) {
							return types.NewErr("invalid start or length to 'substr'")
						}
						return types.String(runes[start : start+end])
					} else {
						return types.NewErr("too many arguments to 'substr'")
					}
				},
			},
		),
	}
	return c
}

func (c *CustomLib) CompileOptions() []cel.EnvOption {
	return c.envOptions
}
func (c *CustomLib) ProgramOptions() []cel.ProgramOption {
	return c.programOptions
}

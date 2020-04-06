// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package proto

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6f7b1b5bDecodeGithubComGoTrellisTrellisMessageProto(in *jlexer.Lexer, out *Service) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			out.Type = RegistryType(in.Int32())
		case "BaseService":
			if in.IsNull() {
				in.Skip()
				out.BaseService = nil
			} else {
				if out.BaseService == nil {
					out.BaseService = new(BaseService)
				}
				(*out.BaseService).UnmarshalEasyJSON(in)
			}
		case "Endpoint":
			out.Endpoint = string(in.String())
		case "metadata":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Metadata = make(map[string]string)
				} else {
					out.Metadata = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 string
					v1 = string(in.String())
					(out.Metadata)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6f7b1b5bEncodeGithubComGoTrellisTrellisMessageProto(out *jwriter.Writer, in Service) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Type != 0 {
		const prefix string = ",\"type\":"
		first = false
		out.RawString(prefix[1:])
		out.Int32(int32(in.Type))
	}
	if in.BaseService != nil {
		const prefix string = ",\"BaseService\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.BaseService).MarshalEasyJSON(out)
	}
	if in.Endpoint != "" {
		const prefix string = ",\"Endpoint\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Endpoint))
	}
	if len(in.Metadata) != 0 {
		const prefix string = ",\"metadata\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.Metadata {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				out.String(string(v2Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Service) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6f7b1b5bEncodeGithubComGoTrellisTrellisMessageProto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Service) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6f7b1b5bEncodeGithubComGoTrellisTrellisMessageProto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Service) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6f7b1b5bDecodeGithubComGoTrellisTrellisMessageProto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Service) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6f7b1b5bDecodeGithubComGoTrellisTrellisMessageProto(l, v)
}
func easyjson6f7b1b5bDecodeGithubComGoTrellisTrellisMessageProto1(in *jlexer.Lexer, out *BaseService) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "version":
			out.Version = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6f7b1b5bEncodeGithubComGoTrellisTrellisMessageProto1(out *jwriter.Writer, in BaseService) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if in.Version != "" {
		const prefix string = ",\"version\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Version))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BaseService) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6f7b1b5bEncodeGithubComGoTrellisTrellisMessageProto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BaseService) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6f7b1b5bEncodeGithubComGoTrellisTrellisMessageProto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BaseService) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6f7b1b5bDecodeGithubComGoTrellisTrellisMessageProto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BaseService) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6f7b1b5bDecodeGithubComGoTrellisTrellisMessageProto1(l, v)
}

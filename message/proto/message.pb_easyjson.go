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

func easyjson52bdc5a7DecodeGithubComGoTrellisTrellisMessageProto(in *jlexer.Lexer, out *Payload) {
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
		case "trace_id":
			out.TraceId = string(in.String())
		case "id":
			out.Id = string(in.String())
		case "service_name":
			out.ServiceName = string(in.String())
		case "service_version":
			out.ServiceVersion = string(in.String())
		case "topic":
			out.Topic = string(in.String())
		case "header":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Header = make(map[string]string)
				} else {
					out.Header = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 string
					v1 = string(in.String())
					(out.Header)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "content_type":
			out.ContentType = string(in.String())
		case "body":
			if in.IsNull() {
				in.Skip()
				out.Body = nil
			} else {
				out.Body = in.Bytes()
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
func easyjson52bdc5a7EncodeGithubComGoTrellisTrellisMessageProto(out *jwriter.Writer, in Payload) {
	out.RawByte('{')
	first := true
	_ = first
	if in.TraceId != "" {
		const prefix string = ",\"trace_id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.TraceId))
	}
	if in.Id != "" {
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Id))
	}
	if in.ServiceName != "" {
		const prefix string = ",\"service_name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ServiceName))
	}
	if in.ServiceVersion != "" {
		const prefix string = ",\"service_version\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ServiceVersion))
	}
	if in.Topic != "" {
		const prefix string = ",\"topic\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Topic))
	}
	if len(in.Header) != 0 {
		const prefix string = ",\"header\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v3First := true
			for v3Name, v3Value := range in.Header {
				if v3First {
					v3First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v3Name))
				out.RawByte(':')
				out.String(string(v3Value))
			}
			out.RawByte('}')
		}
	}
	if in.ContentType != "" {
		const prefix string = ",\"content_type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ContentType))
	}
	if len(in.Body) != 0 {
		const prefix string = ",\"body\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Base64Bytes(in.Body)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Payload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson52bdc5a7EncodeGithubComGoTrellisTrellisMessageProto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Payload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson52bdc5a7EncodeGithubComGoTrellisTrellisMessageProto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Payload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson52bdc5a7DecodeGithubComGoTrellisTrellisMessageProto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Payload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson52bdc5a7DecodeGithubComGoTrellisTrellisMessageProto(l, v)
}
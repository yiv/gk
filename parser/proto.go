package parser

import (
	"bytes"
	"fmt"
	"github.com/emicklei/proto"
)

type Proto struct {
	PackageName string
	ServiceName string
	Methods     []Method
	Messages    []Struct
}

func NewProto() *Proto {
	return &Proto{}
}

func (p *Proto) handleService() func(s *proto.Service) {
	return func(s *proto.Service) {
		p.ServiceName = s.Name
	}
}
func (p *Proto) handleMessage() func(s *proto.Message) {
	return func(s *proto.Message) {
		fields := []NamedTypeValue{}
		for _, v := range s.Elements {
			switch v.(type) {
			case *proto.NormalField:
				n := v.(*proto.NormalField)
				seq := fmt.Sprintf("%v", n.Sequence)
				if n.Repeated {
					fields = append(fields, NewNameTypeValue(n.Name, "repeated "+n.Type, seq))
				} else {
					fields = append(fields, NewNameTypeValue(n.Name, n.Type, seq))
				}
			case *proto.MapField:
				n := v.(*proto.MapField)
				seq := fmt.Sprintf("%v", n.Sequence)
				typeName := fmt.Sprintf("map<%v, %v>", n.KeyType, n.Type)
				fields = append(fields, NewNameTypeValue(n.Name, typeName, seq))
			}
		}
		message := Struct{Name: s.Name, Vars: fields}
		p.Messages = append(p.Messages, message)
	}
}

func (p *Proto) handleRPC() func(s *proto.RPC) {
	return func(s *proto.RPC) {
		reqParam := []NamedTypeValue{NewNameType("", s.RequestType)}
		resParam := []NamedTypeValue{NewNameType("", s.ReturnsType)}
		method := Method{Name: s.Name, Parameters: reqParam, Results: resParam}
		p.Methods = append(p.Methods, method)
	}
}

type ProtoParser struct{}

func NewProtoParser() *ProtoParser {
	return &ProtoParser{}
}

func (pp *ProtoParser) Parse(src []byte) (*Proto, error) {
	p := NewProto()
	reader := bytes.NewReader(src)
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	proto.Walk(definition,
		proto.WithService(p.handleService()),
		proto.WithMessage(p.handleMessage()),
		proto.WithRPC(p.handleRPC()),
	)
	return p, nil
}

package interceptor

import (
	"buf.build/go/protovalidate"
	protovalidatemw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

// UnaryProtovalidateInterceptor создаёт unary-интерцептор, который проверяет
// входящие сообщения по правилам Protovalidate из дескрипторов protobuf.
func UnaryProtovalidateInterceptor() (grpc.UnaryServerInterceptor, error) {
	v, err := protovalidate.New()
	if err != nil {
		return nil, err
	}
	return protovalidatemw.UnaryServerInterceptor(v), nil
}

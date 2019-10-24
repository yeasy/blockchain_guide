## ProtoBuf 与 gRPC
[ProtoBuf](https://github.com/google/protobuf) 是一套接口描述语言（[Interface Definition Language，IDL](https://en.wikipedia.org/wiki/Interface_description_language)），类似 Apache 的 [Thrift](https://thrift.apache.org/)。

相关处理工具主要是 [protoc](https://github.com/google/protobuf)，基于 C++ 语言实现。

用户写好 `.proto` 描述文件，之后便可以使用 protoc 自动编译生成众多计算机语言（C++、Java、Python、C#、Golang 等）的接口代码。这些代码可以支持 gRPC，也可以不支持。

[gRPC](https://github.com/grpc/grpc) 是 Google 开源的 RPC 框架和库，已支持主流计算机语言。底层通信采用 HTTP2 协议，比较适合互联网场景。gRPC 在设计上考虑了跟 ProtoBuf 的配合使用。

两者分别解决不同问题，可以配合使用，也可以分开单独使用。

典型的配合使用场景是，写好 `.proto` 描述文件定义 RPC 的接口，然后用 protoc（带 gRPC 插件）基于 `.proto` 模板自动生成客户端和服务端的接口代码。

### ProtoBuf

需要工具主要包括：

* 编译工具：[protoc](https://github.com/google/protobuf)，以及一些官方没有带的语言插件；
* 运行环境：各种语言的 protobuf 库，不同语言有不同的安装来源；

语法类似 C++ 语言，可以参考 ProtoBuf 语言规范：[https://developers.google.com/protocol-buffers/docs/proto](https://developers.google.com/protocol-buffers/docs/proto)。

比较核心的，`message` 是代表数据结构（里面可以包括不同类型的成员变量，包括字符串、数字、数组、字典……），`service` 代表 RPC 接口。变量后面的数字是代表进行二进制编码时候的提示信息，1~15 表示热变量，会用较少的字节来编码。另外，支持导入。

默认所有变量都是可选的（optional），repeated 则表示数组。主要 service rpc 接口接受单个 message 参数，返回单个 message。如下所示。

```protobuf
syntax = "proto3";
package hello;

message HelloRequest {
  string greeting = 1;
}

message HelloResponse {
  string reply = 1;
  repeated int32 number=4;
}

service HelloService {
  rpc SayHello(HelloRequest) returns (HelloResponse){}
}
```

编译最关键的参数是输出语言格式参数，例如，python 为 `--python_out=OUT_DIR`。

一些还没有官方支持的语言，可以通过安装 protoc 对应的 plugin 来支持。例如，对于 Go 语言，可以安装

```sh
$ go get -u github.com/golang/protobuf/{protoc-gen-go,proto} // 前者是 plugin；后者是 go 的依赖库
```

之后，正常使用 `protoc --go_out=./ hello.proto` 来生成 hello.pb.go，会自动调用 `protoc-gen-go` 插件。

ProtoBuf 提供了 `Marshal/Unmarshal` 方法来将数据结构进行序列化操作。所生成的二进制文件在存储效率上比 XML 高 3~10 倍，并且处理性能高 1~2 个数量级。

### gRPC

相关工具主要包括：

* 运行时库：各种不同语言有不同的[安装方法](https://github.com/grpc/grpc/blob/master/INSTALL.md)，主流语言的包管理器都已支持。
* protoc，以及 gRPC 插件和其它插件：采用 ProtoBuf 作为 IDL 时，对 .proto 文件进行编译处理。

类似其它 RPC 框架，gRPC 的库在服务端提供一个 gRPC Server，客户端的库是 gRPC Stub。典型的场景是客户端发送请求，同步或异步调用服务端的接口。客户端和服务端之间的通信协议是基于 HTTP2 的 [gRPC](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md) 协议，支持双工的流式保序消息，性能比较好，同时也很轻。

采用 ProtoBuf 作为 IDL，则需要定义 service 类型。生成客户端和服务端代码。用户自行实现服务端代码中的调用接口，并且利用客户端代码来发起请求到服务端。一个完整的例子可以参考 [https://github.com/grpc/grpc-go/blob/master/examples/helloworld](https://github.com/grpc/grpc-go/blob/master/examples/helloworld/)。

以上面 proto 文件为例，需要执行时添加 gRPC 的 plugin：

```sh
$ protoc --go_out=plugins=grpc:. hello.proto
```

gRPC 更多原理可以参考[官方文档：http://www.grpc.io/docs](http://www.grpc.io/docs/)。

#### 生成服务端代码

服务端相关代码如下，主要定义了 HelloServiceServer 接口，用户可以自行编写实现代码。

```go
type HelloServiceServer interface {
        SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
}

func RegisterHelloServiceServer(s *grpc.Server, srv HelloServiceServer) {
        s.RegisterService(&_HelloService_serviceDesc, srv)
}
```

用户需要自行实现服务端接口，代码如下。

比较重要的，创建并启动一个 gRPC 服务的过程：

* 创建监听套接字：`lis, err := net.Listen("tcp", port)`；
* 创建服务端：`grpc.NewServer()`；
* 注册服务：`pb.RegisterHelloServiceServer()`；
* 启动服务端：`s.Serve(lis)`。

```go
type server struct{}

// 这里实现服务端接口中的方法。
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

// 创建并启动一个 gRPC 服务的过程：创建监听套接字、创建服务端、注册服务、启动服务端。
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	s.Serve(lis)
}
```

编译并启动服务端。

#### 生成客户端代码

生成的 go 文件中客户端相关代码如下，主要和实现了 HelloServiceClient 接口。用户可以通过 gRPC 来直接调用这个接口。

```go
type HelloServiceClient interface {
        SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
}

type helloServiceClient struct {
        cc *grpc.ClientConn
}

func NewHelloServiceClient(cc *grpc.ClientConn) HelloServiceClient {
        return &helloServiceClient{cc}
}

func (c *helloServiceClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
        out := new(HelloResponse)
        err := grpc.Invoke(ctx, "/hello.HelloService/SayHello", in, out, c.cc, opts...)
        if err != nil {
                return nil, err
        }
        return out, nil
}
```

用户直接调用接口方法：创建连接、创建客户端、调用接口。

```go
func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
```

编译并启动客户端，查看到服务端返回的消息。

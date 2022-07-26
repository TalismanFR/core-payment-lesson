// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.3
// source: auth-service.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IdentityProviderClient is the client API for IdentityProvider service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IdentityProviderClient interface {
	SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*JWTokens, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*JWTokens, error)
	Verify(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error)
	Refresh(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*JWTokens, error)
	Revoke(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error)
}

type identityProviderClient struct {
	cc grpc.ClientConnInterface
}

func NewIdentityProviderClient(cc grpc.ClientConnInterface) IdentityProviderClient {
	return &identityProviderClient{cc}
}

func (c *identityProviderClient) SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*JWTokens, error) {
	out := new(JWTokens)
	err := c.cc.Invoke(ctx, "/server.IdentityProvider/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityProviderClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*JWTokens, error) {
	out := new(JWTokens)
	err := c.cc.Invoke(ctx, "/server.IdentityProvider/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityProviderClient) Verify(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error) {
	out := new(wrapperspb.BoolValue)
	err := c.cc.Invoke(ctx, "/server.IdentityProvider/Verify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityProviderClient) Refresh(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*JWTokens, error) {
	out := new(JWTokens)
	err := c.cc.Invoke(ctx, "/server.IdentityProvider/Refresh", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityProviderClient) Revoke(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error) {
	out := new(wrapperspb.BoolValue)
	err := c.cc.Invoke(ctx, "/server.IdentityProvider/Revoke", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IdentityProviderServer is the server API for IdentityProvider service.
// All implementations must embed UnimplementedIdentityProviderServer
// for forward compatibility
type IdentityProviderServer interface {
	SignUp(context.Context, *SignUpRequest) (*JWTokens, error)
	Login(context.Context, *LoginRequest) (*JWTokens, error)
	Verify(context.Context, *wrapperspb.StringValue) (*wrapperspb.BoolValue, error)
	Refresh(context.Context, *wrapperspb.StringValue) (*JWTokens, error)
	Revoke(context.Context, *wrapperspb.StringValue) (*wrapperspb.BoolValue, error)
	mustEmbedUnimplementedIdentityProviderServer()
}

// UnimplementedIdentityProviderServer must be embedded to have forward compatible implementations.
type UnimplementedIdentityProviderServer struct {
}

func (UnimplementedIdentityProviderServer) SignUp(context.Context, *SignUpRequest) (*JWTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedIdentityProviderServer) Login(context.Context, *LoginRequest) (*JWTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedIdentityProviderServer) Verify(context.Context, *wrapperspb.StringValue) (*wrapperspb.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verify not implemented")
}
func (UnimplementedIdentityProviderServer) Refresh(context.Context, *wrapperspb.StringValue) (*JWTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (UnimplementedIdentityProviderServer) Revoke(context.Context, *wrapperspb.StringValue) (*wrapperspb.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revoke not implemented")
}
func (UnimplementedIdentityProviderServer) mustEmbedUnimplementedIdentityProviderServer() {}

// UnsafeIdentityProviderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IdentityProviderServer will
// result in compilation errors.
type UnsafeIdentityProviderServer interface {
	mustEmbedUnimplementedIdentityProviderServer()
}

func RegisterIdentityProviderServer(s grpc.ServiceRegistrar, srv IdentityProviderServer) {
	s.RegisterService(&IdentityProvider_ServiceDesc, srv)
}

func _IdentityProvider_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityProviderServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.IdentityProvider/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityProviderServer).SignUp(ctx, req.(*SignUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityProvider_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityProviderServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.IdentityProvider/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityProviderServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityProvider_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityProviderServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.IdentityProvider/Verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityProviderServer).Verify(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityProvider_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityProviderServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.IdentityProvider/Refresh",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityProviderServer).Refresh(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityProvider_Revoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityProviderServer).Revoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.IdentityProvider/Revoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityProviderServer).Revoke(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

// IdentityProvider_ServiceDesc is the grpc.ServiceDesc for IdentityProvider service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IdentityProvider_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "server.IdentityProvider",
	HandlerType: (*IdentityProviderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignUp",
			Handler:    _IdentityProvider_SignUp_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _IdentityProvider_Login_Handler,
		},
		{
			MethodName: "Verify",
			Handler:    _IdentityProvider_Verify_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _IdentityProvider_Refresh_Handler,
		},
		{
			MethodName: "Revoke",
			Handler:    _IdentityProvider_Revoke_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth-service.proto",
}
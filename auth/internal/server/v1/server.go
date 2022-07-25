// TODO: request validation
// TODO: return grpc error codes

package v1

import (
	"auth/internal/api/v1"
	"auth/internal/auth"
	"auth/internal/config"
	"auth/pkg/logger"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net"
)

var _ v1.IdentityProviderServer = (*Server)(nil)

type Server struct {
	v1.UnimplementedIdentityProviderServer
	a auth.AuthenticationService
	l logger.Logger
}

func New(a auth.AuthenticationService, l logger.Logger) *Server {
	return &Server{a: a, l: l}
}

func (s *Server) Run(cfg *config.GrpcConfig) error {
	gs := grpc.NewServer()
	v1.RegisterIdentityProviderServer(gs, s)
	ls, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return err
	}
	return gs.Serve(ls)
}

func (s *Server) SignUp(ctx context.Context, request *v1.SignUpRequest) (*v1.JWTokens, error) {
	sr := auth.SignUpRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
		Role:     request.GetRole(),
	}

	tokens, err := s.a.Signup(ctx, sr)
	if err != nil {
		return nil, err
	}

	return s.newJWTokens(tokens.Access, tokens.Refresh), nil
}

func (s *Server) Login(ctx context.Context, request *v1.LoginRequest) (*v1.JWTokens, error) {
	lr := auth.LoginRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	tokens, err := s.a.Login(ctx, lr)
	if err != nil {
		return nil, err
	}

	return s.newJWTokens(tokens.Access, tokens.Refresh), nil
}

func (s *Server) Refresh(ctx context.Context, value *wrapperspb.StringValue) (*v1.JWTokens, error) {
	refreshToken := value.GetValue()

	tokens, err := s.a.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return s.newJWTokens(tokens.Access, tokens.Refresh), nil
}

func (s *Server) Verify(ctx context.Context, value *wrapperspb.StringValue) (*wrapperspb.BoolValue, error) {
	accessToken := value.GetValue()

	err := s.a.Verify(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return wrapperspb.Bool(true), nil
}

func (s *Server) Revoke(ctx context.Context, value *wrapperspb.StringValue) (*wrapperspb.BoolValue, error) {
	accessToken := value.GetValue()

	err := s.a.Revoke(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return wrapperspb.Bool(true), nil
}

func (s *Server) newJWTokens(access string, refresh string) *v1.JWTokens {
	return &v1.JWTokens{
		Access:  access,
		Refresh: refresh,
	}
}

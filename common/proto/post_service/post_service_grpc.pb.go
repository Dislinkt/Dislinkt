// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: post_service.proto

package post_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PostServiceClient is the client API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostServiceClient interface {
	GetRecent(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetMultipleResponse, error)
	GetAllByUserId(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetMultipleResponse, error)
	GetAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetMultipleResponse, error)
	CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*Empty, error)
	CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error)
	LikePost(ctx context.Context, in *ReactionRequest, opts ...grpc.CallOption) (*Empty, error)
	DislikePost(ctx context.Context, in *ReactionRequest, opts ...grpc.CallOption) (*Empty, error)
	CreateJobOffer(ctx context.Context, in *CreateJobOfferRequest, opts ...grpc.CallOption) (*Empty, error)
	GetAllJobOffers(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllJobOffers, error)
	GetAllLikesForPost(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReactionsResponse, error)
	GetAllDislikesForPost(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReactionsResponse, error)
	GetAllCommentsForPost(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetAllCommentsResponse, error)
}

type postServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostServiceClient(cc grpc.ClientConnInterface) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) GetRecent(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetMultipleResponse, error) {
	out := new(GetMultipleResponse)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/getRecent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetAllByUserId(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetMultipleResponse, error) {
	out := new(GetMultipleResponse)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/getAllByUserId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetMultipleResponse, error) {
	out := new(GetMultipleResponse)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/getAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/createPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error) {
	out := new(CreateCommentResponse)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/createComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) LikePost(ctx context.Context, in *ReactionRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/likePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) DislikePost(ctx context.Context, in *ReactionRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/dislikePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) CreateJobOffer(ctx context.Context, in *CreateJobOfferRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/createJobOffer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetAllJobOffers(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllJobOffers, error) {
	out := new(GetAllJobOffers)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/getAllJobOffers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetAllLikesForPost(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReactionsResponse, error) {
	out := new(GetReactionsResponse)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/getAllLikesForPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetAllDislikesForPost(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReactionsResponse, error) {
	out := new(GetReactionsResponse)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/getAllDislikesForPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetAllCommentsForPost(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetAllCommentsResponse, error) {
	out := new(GetAllCommentsResponse)
	err := c.cc.Invoke(ctx, "/post_service_proto.PostService/getAllCommentsForPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServiceServer is the server API for PostService service.
// All implementations must embed UnimplementedPostServiceServer
// for forward compatibility
type PostServiceServer interface {
	GetRecent(context.Context, *GetRequest) (*GetMultipleResponse, error)
	GetAllByUserId(context.Context, *GetRequest) (*GetMultipleResponse, error)
	GetAll(context.Context, *Empty) (*GetMultipleResponse, error)
	CreatePost(context.Context, *CreatePostRequest) (*Empty, error)
	CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error)
	LikePost(context.Context, *ReactionRequest) (*Empty, error)
	DislikePost(context.Context, *ReactionRequest) (*Empty, error)
	CreateJobOffer(context.Context, *CreateJobOfferRequest) (*Empty, error)
	GetAllJobOffers(context.Context, *Empty) (*GetAllJobOffers, error)
	GetAllLikesForPost(context.Context, *GetRequest) (*GetReactionsResponse, error)
	GetAllDislikesForPost(context.Context, *GetRequest) (*GetReactionsResponse, error)
	GetAllCommentsForPost(context.Context, *GetRequest) (*GetAllCommentsResponse, error)
	mustEmbedUnimplementedPostServiceServer()
}

// UnimplementedPostServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPostServiceServer struct {
}

func (UnimplementedPostServiceServer) GetRecent(context.Context, *GetRequest) (*GetMultipleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecent not implemented")
}
func (UnimplementedPostServiceServer) GetAllByUserId(context.Context, *GetRequest) (*GetMultipleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllByUserId not implemented")
}
func (UnimplementedPostServiceServer) GetAll(context.Context, *Empty) (*GetMultipleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedPostServiceServer) CreatePost(context.Context, *CreatePostRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}
func (UnimplementedPostServiceServer) CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateComment not implemented")
}
func (UnimplementedPostServiceServer) LikePost(context.Context, *ReactionRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikePost not implemented")
}
func (UnimplementedPostServiceServer) DislikePost(context.Context, *ReactionRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DislikePost not implemented")
}
func (UnimplementedPostServiceServer) CreateJobOffer(context.Context, *CreateJobOfferRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateJobOffer not implemented")
}
func (UnimplementedPostServiceServer) GetAllJobOffers(context.Context, *Empty) (*GetAllJobOffers, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllJobOffers not implemented")
}
func (UnimplementedPostServiceServer) GetAllLikesForPost(context.Context, *GetRequest) (*GetReactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllLikesForPost not implemented")
}
func (UnimplementedPostServiceServer) GetAllDislikesForPost(context.Context, *GetRequest) (*GetReactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllDislikesForPost not implemented")
}
func (UnimplementedPostServiceServer) GetAllCommentsForPost(context.Context, *GetRequest) (*GetAllCommentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllCommentsForPost not implemented")
}
func (UnimplementedPostServiceServer) mustEmbedUnimplementedPostServiceServer() {}

// UnsafePostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServiceServer will
// result in compilation errors.
type UnsafePostServiceServer interface {
	mustEmbedUnimplementedPostServiceServer()
}

func RegisterPostServiceServer(s grpc.ServiceRegistrar, srv PostServiceServer) {
	s.RegisterService(&PostService_ServiceDesc, srv)
}

func _PostService_GetRecent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetRecent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/getRecent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetRecent(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetAllByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetAllByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/getAllByUserId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetAllByUserId(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/getAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetAll(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_CreatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).CreatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/createPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).CreatePost(ctx, req.(*CreatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_CreateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).CreateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/createComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).CreateComment(ctx, req.(*CreateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_LikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).LikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/likePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).LikePost(ctx, req.(*ReactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_DislikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).DislikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/dislikePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).DislikePost(ctx, req.(*ReactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_CreateJobOffer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateJobOfferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).CreateJobOffer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/createJobOffer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).CreateJobOffer(ctx, req.(*CreateJobOfferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetAllJobOffers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetAllJobOffers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/getAllJobOffers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetAllJobOffers(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetAllLikesForPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetAllLikesForPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/getAllLikesForPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetAllLikesForPost(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetAllDislikesForPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetAllDislikesForPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/getAllDislikesForPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetAllDislikesForPost(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetAllCommentsForPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetAllCommentsForPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service_proto.PostService/getAllCommentsForPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetAllCommentsForPost(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostService_ServiceDesc is the grpc.ServiceDesc for PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "post_service_proto.PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getRecent",
			Handler:    _PostService_GetRecent_Handler,
		},
		{
			MethodName: "getAllByUserId",
			Handler:    _PostService_GetAllByUserId_Handler,
		},
		{
			MethodName: "getAll",
			Handler:    _PostService_GetAll_Handler,
		},
		{
			MethodName: "createPost",
			Handler:    _PostService_CreatePost_Handler,
		},
		{
			MethodName: "createComment",
			Handler:    _PostService_CreateComment_Handler,
		},
		{
			MethodName: "likePost",
			Handler:    _PostService_LikePost_Handler,
		},
		{
			MethodName: "dislikePost",
			Handler:    _PostService_DislikePost_Handler,
		},
		{
			MethodName: "createJobOffer",
			Handler:    _PostService_CreateJobOffer_Handler,
		},
		{
			MethodName: "getAllJobOffers",
			Handler:    _PostService_GetAllJobOffers_Handler,
		},
		{
			MethodName: "getAllLikesForPost",
			Handler:    _PostService_GetAllLikesForPost_Handler,
		},
		{
			MethodName: "getAllDislikesForPost",
			Handler:    _PostService_GetAllDislikesForPost_Handler,
		},
		{
			MethodName: "getAllCommentsForPost",
			Handler:    _PostService_GetAllCommentsForPost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "post_service.proto",
}

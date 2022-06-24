package services

import (
	"fmt"

	"github.com/Hatsker01/Works/api-getaway/config"
	pb "github.com/Hatsker01/Works/api-getaway/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	PostService() pb.PostServiceClient
}

type serviceManager struct {
	postService pb.PostServiceClient
}

func (s *serviceManager) PostService() pb.PostServiceClient {
	return s.postService
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.PostServiceHost, conf.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		postService: pb.NewPostServiceClient(connPost),
	}

	return serviceManager, nil
}

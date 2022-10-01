package service_mocks

import (
	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/mock"
)

// TODO: mock fiberCtx and Doc.Get

type MockFSClient struct {
	mock.Mock
}

func (m *MockFSClient) MockGetEventById(id string) (*firestore.DocumentSnapshot, error) {
	args := m.Called(id)
	return args.Get(0).(*firestore.DocumentSnapshot), args.Error(0)
}

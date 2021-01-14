package main

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBrokerPublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		topic = "test_topic"
		value = []byte("test_value")
	)

	mockStore := NewMockstorer(ctrl)
	mockStore.EXPECT().Insert(topic, value)

	b := newBroker(mockStore)

	assert.NoError(t, b.Publish(topic, value))
}

func TestBrokerSubscribe(t *testing.T) {
}

func TestConsumerNext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		topic = "test_topic"
	)

	mockStore := NewMockstorer(ctrl)
	mockStore.EXPECT().GetNext(topic)

	b := newBroker(mockStore)
	c := b.Subscribe(topic)

	c.Next()
}

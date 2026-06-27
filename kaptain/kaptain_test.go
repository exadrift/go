package kaptain

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// blocking channel max of one written item at once
var leaderChannel chan string = make(chan string, 1)

type DummyElector struct {
	acquireFailureTimeout    time.Duration
	breakLeaseFailureTimeout time.Duration
}

// NewDummyElector returns a dummy elector
// acquireFailureTimeout represents the time which the elector will wait before encountering a simulated failure during the
// leader acquisition process
func NewDummyElector(acquireFailureTimeout time.Duration, breakLeaseFailureTimeout time.Duration) *DummyElector {
	if acquireFailureTimeout == 0 {
		acquireFailureTimeout = time.Hour // considered an unattainable time for the sake of unit tests
	}
	return &DummyElector{
		acquireFailureTimeout:    acquireFailureTimeout,
		breakLeaseFailureTimeout: breakLeaseFailureTimeout,
	}
}

func (de *DummyElector) Wait(ctx context.Context, applicationKey string, callback func(context.Context)) {
	instanceKey := uuid.NewString()

	select {
	case leaderChannel <- instanceKey:
	case <-time.After(de.acquireFailureTimeout):
		return
	}

	go callback(ctx)

	if de.breakLeaseFailureTimeout > 0 {
		select {
		case <-time.After(de.breakLeaseFailureTimeout):
			// wait until the timer expires
		case <-ctx.Done():
			// wait until the context is cancelled
		}
	} else {
		// wait for context to close
		<-ctx.Done()
	}

	// read back the item
	readKey := <-leaderChannel
	if readKey != instanceKey {
		panic(fmt.Sprintf("expected leader key \"%s\" differs from actual key \"%s\"", instanceKey, readKey))
	}
}

func TestKaptainDoNothing(t *testing.T) {
	leaderElectedCount := 0
	doNothing := func(ctx context.Context) {
		leaderElectedCount++
	}

	appKey := "test"
	cap1, err := NewKaptain(appKey, NewDummyElector(0, 0))
	assert.NoError(t, err)
	cap2, err := NewKaptain(appKey, NewDummyElector(0, 0))
	assert.NoError(t, err)
	cap3, err := NewKaptain(appKey, NewDummyElector(0, 0))
	assert.NoError(t, err)

	waiters := sync.WaitGroup{}
	for _, e := range []*Kaptain{cap1, cap2, cap3} {
		waiters.Add(1)
		go func() {
			defer waiters.Done()
			e.Wait(context.Background(), doNothing)
		}()
	}
	waiters.Wait()
	assert.Equal(t, 3, leaderElectedCount)
}

func TestKaptainDoContextClosureAfterElection(t *testing.T) {
	leaderElectedCount := 0
	waiters := sync.WaitGroup{}
	for range 3 {
		closeContext, cancel := context.WithCancel(context.Background())
		waitForContextClosure := func(ctx context.Context) {
			leaderElectedCount++
			// call the cancel function in the leader loop so that timers aren't used to fire it after
			// leader election is established
			cancel()
			for {
				select {
				case <-time.After(time.Minute):
					panic("we shouldn't hit the timeout")
				case <-closeContext.Done():
					// exit the leader func when the context is closed
					return
				}
			}
		}

		appKey := "test"
		cap, err := NewKaptain(appKey, NewDummyElector(0, 0))
		assert.NoError(t, err)

		waiters.Add(1)
		go func() {
			defer waiters.Done()
			cap.Wait(closeContext, waitForContextClosure)
		}()
	}

	waiters.Wait()
	assert.Equal(t, 3, leaderElectedCount)
}

func TestKaptainDoContextClosureBeforeElection(t *testing.T) {
	waiters := sync.WaitGroup{}
	for range 3 {
		closeContext, cancel := context.WithCancel(context.Background())
		waitForContextClosure := func(ctx context.Context) {
			// call the cancel function in the leader loop so that timers aren't used to fire it after
			// leader election is established
			for {
				select {
				case <-time.After(time.Minute):
					panic("we shouldn't hit the timeout")
				case <-closeContext.Done():
					// exit the leader func when the context is closed
					return
				}
			}
		}

		appKey := "test"
		cap, err := NewKaptain(appKey, NewDummyElector(0, 0))
		assert.NoError(t, err)

		cancel()
		waiters.Add(1)
		go func() {
			defer waiters.Done()
			cap.Wait(closeContext, waitForContextClosure)
		}()
	}

	waiters.Wait()
}

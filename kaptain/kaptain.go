package kaptain

import (
	"context"
	"sync"
)

type Kaptain struct {
	elector        ElectorInterface
	applicationKey string
}

// NewKaptain returns a new Kaptain instance
func NewKaptain(applicationKey string, elector ...ElectorInterface) (*Kaptain, error) {
	if len(elector) > 1 {
		panic("multiple electors are not supported")
	}

	var electorInf ElectorInterface
	if len(elector) == 1 {
		electorInf = elector[0]
	} else {
		// make a default kubernetes elector if none are specified
		var err error
		electorInf, err = NewKubeElector()
		if err != nil {
			return nil, err
		}
	}

	return &Kaptain{
		elector:        electorInf,
		applicationKey: applicationKey,
	}, nil
}

// Wait will block on the leader election process until either the supplied context is cancelled, or
// the leader has been elected and leadership status is relinquished
func (k *Kaptain) Wait(ctx context.Context, leaderFunc func(context.Context)) {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg := sync.WaitGroup{}
	k.elector.Wait(cancelCtx, k.applicationKey, func(ctx context.Context) {
		// cancellation needs to take place in here because if this function exits
		// we want to force the leader hold process to end
		defer cancel()
		defer wg.Done()

		wg.Add(1)
		leaderFunc(cancelCtx)
	})
	// this will force the inner leader function to cancel when this exits, if for any reason we have an exit
	// of the elector wait function, where the context itself was not cancelled before exit
	cancel()
	wg.Wait()
}

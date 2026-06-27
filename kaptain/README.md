# kaptain
kaptain is a leader election helper which wraps and simplifies the leaderelection module provided by the kubernetes module.  a custom leader elector can also be used if desired, by implementing a new elector using the `ElectorInterface`.

## usage
```
import (
  "github.com/exadrift/go/kaptain"
)

func electionFunc(ctx context.Context) {
  // main event loop
  for {
    select {
      case <-ctx.Done():
        // interrupt if the context is cancelled
      ... 
      // either a default clause, if this is a polling loop and context check is non-blocking
      // or wait for interrupting events to fire by checking other channels in here
    }
  }
}

...

leaderElector, err := kaptain.NewKaptain("my-application-name")
if err != nil {
  log.Fatal(err)
}

// this produces a cancellable context, for which you can call the cancel function anytime
// to force the leader election cycle to end and unblock the Wait method below
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// this will block until either the context is cancelled, election fails due to internal error
// or election happened, and then was given up for whatever reason
leaderElector.Wait(ctx, electionFunc)
```

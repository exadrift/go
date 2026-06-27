package kaptain

import (
	"context"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type ElectorInterface interface {
	// Wait is a blocking call, which exits when either the context has been cancelled, or the leader has been
	// elected and election has been forfeitted after being held
	// callback must be executed in a goroutine by the implementer of this interface
	Wait(ctx context.Context, applicationKey string, callback func(ctx context.Context))
}

// KubeElector is a kubernetes based leader elector
type KubeElector struct {
	clientset *kubernetes.Clientset
}

type KubeElectorOptionType int

const (
	KubeElectorOptionClientset KubeElectorOptionType = iota
)

type KubeElectorOption struct {
	Type KubeElectorOptionType
	data any
}

func WithCliensetOption(clientset *kubernetes.Clientset) *KubeElectorOption {
	return &KubeElectorOption{
		Type: KubeElectorOptionClientset,
		data: clientset,
	}
}

// NewKubeElector constructs a new kubernetes based leader elector
func NewKubeElector(options ...*KubeElectorOption) (*KubeElector, error) {
	var clientset *kubernetes.Clientset
	for _, option := range options {
		switch option.Type {
		case KubeElectorOptionClientset:
			clientset = option.data.(*kubernetes.Clientset)
		}
	}

	if clientset == nil {
		// if not set in an option, default it
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
	}

	return &KubeElector{
		clientset: clientset,
	}, nil
}

func (e *KubeElector) Wait(ctx context.Context, applicationKey string, callback func(ctx context.Context)) {
	candidateIdentity := uuid.NewString()

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      applicationKey,
			Namespace: "default",
		},
		Client: e.clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: candidateIdentity,
		},
	}

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: callback,
		},
	})
}

package scheme

import (
	"kubeception.cloud/kubeception/pkg/apis/kubeception/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	Scheme *runtime.Scheme

	CodecFactory serializer.CodecFactory

	Decoder runtime.Decoder
)

func init() {
	Scheme = runtime.NewScheme()

    utilruntime.Must(v1alpha1.AddToScheme(Scheme))

	CodecFactory = serializer.NewCodecFactory(Scheme)
	Decoder = CodecFactory.UniversalDecoder()
}


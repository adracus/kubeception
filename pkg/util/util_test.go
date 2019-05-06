package util

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMachine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Util")
}

var _ = Describe("Utils Suite", func() {
	Describe("#ContextFromStopChannel", func() {
		It("should create a context that is open as long as the stop channel is open", func() {
			stopCh := make(chan struct{})

			ctx := ContextFromStopChannel(stopCh)

			Expect(ctx.Err()).To(BeNil())

			close(stopCh)

			Eventually(ctx.Err).Should(Equal(context.Canceled))
		})
	})
})

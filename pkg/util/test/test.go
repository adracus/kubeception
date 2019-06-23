package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RegisterFailHandlerAndRunSpecs(t *testing.T, description string) {
	RegisterFailHandler(Fail)
	RunSpecs(t, description)
}

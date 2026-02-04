package tests

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test_Stater(t *testing.T) {
	RegisterFailHandler(Fail)
	defer GinkgoRecover()
}

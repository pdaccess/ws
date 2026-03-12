package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/pdaccess/ws/internal/platform/servers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var server *httptest.Server
var mockSvc = NewMockService()

func Test_API(t *testing.T) {
	RegisterFailHandler(Fail)
	defer GinkgoRecover()
	RunSpecs(t, "API Suite")
}

var _ = BeforeSuite(func() {
	server = httptest.NewServer(servers.NewHttpServer(mockSvc))
	DeferCleanup(func() {
		server.Close()
	})
})

var _ = Describe("HTTP Server", func() {
	Context("Health endpoint", func() {
		It("should return 200", func() {
			resp, err := server.Client().Get(server.URL + "/health")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})
})

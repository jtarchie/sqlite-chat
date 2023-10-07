package server_test

import (
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/jtarchie/sqlite-chat/server"
	"github.com/jtarchie/sqlite-chat/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("Dashboard", func() {
	When("the users logs in", func() {
		It("renders the initial channel", func() {
			port, err := freeport.GetFreePort()
			Expect(err).NotTo(HaveOccurred())

			db, err := services.NewDB("file:test.db?cache=shared&mode=memory")
			Expect(err).NotTo(HaveOccurred())

			server, err := server.New(db, "", "", "")
			Expect(err).NotTo(HaveOccurred())

			go func() {
				defer GinkgoRecover()

				err := server.Start(fmt.Sprintf("0.0.0.0:%d", port))
				Expect(err).NotTo(HaveOccurred())
			}()

			defer server.Close()

			response := req.MustGet(fmt.Sprintf("http://localhost:%d/dashboard", port))
			body := response.String()
			Expect(body).To(ContainSubstring("general"))
		})
	})
})

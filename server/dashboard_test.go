package server_test

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jtarchie/sqlite-chat/server"
	"github.com/jtarchie/sqlite-chat/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("Dashboard", func() {
	var (
		port int
		db   *sqlx.DB
	)

	BeforeEach(func() {
		var err error

		port, err = freeport.GetFreePort()
		Expect(err).NotTo(HaveOccurred())

		db, err = services.NewDB("file:test.db?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())

		server, err := server.New(db, "", "", "")
		Expect(err).NotTo(HaveOccurred())

		go func() {
			defer server.Close()
			defer GinkgoRecover()

			err := server.Start(fmt.Sprintf("0.0.0.0:%d", port))
			Expect(err).NotTo(HaveOccurred())
		}()
	})

	AfterEach(func() {
		err := db.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	When("the users logs in", func() {
		It("renders the initial channel", func() {
			response := req.MustGet(fmt.Sprintf("http://localhost:%d/dashboard", port))

			Expect(response.IsSuccessState()).To(BeTrue())
			Expect(response.Response.Request.URL.String()).To(ContainSubstring("/dashboard/channels/1"))

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.String()))
			Expect(err).NotTo(HaveOccurred())

			Expect(doc.Find(`a:contains("general")`).Nodes).To(HaveLen(1))
			Expect(doc.Find(`article:contains("Welcome to the chat")`).Nodes).To(HaveLen(1))
		})
	})
})

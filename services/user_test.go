package services_test

import (
	"github.com/jtarchie/sqlite-chat/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	It("finds the channels for a user", func() {
		db, err := services.NewDB("file:test.db?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		service := services.NewUser(db, "bot@example.com")

		channels, err := service.Channels()
		Expect(err).NotTo(HaveOccurred())
		Expect(channels).To(HaveLen(1))
	})
})

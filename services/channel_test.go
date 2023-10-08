package services_test

import (
	"github.com/jtarchie/sqlite-chat/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Channel", func() {
	It("returns the messages for the channel", func() {
		db, err := services.NewDB("file:test.db?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		service := services.NewChannel(db, 1)

		channels, err := service.Messages()
		Expect(err).NotTo(HaveOccurred())
		Expect(channels).To(HaveLen(1))
	})
})

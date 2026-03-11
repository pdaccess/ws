package tests

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database Integration", func() {
	Context("pgvector extension", func() {
		It("should have vector extension enabled", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var extExists bool
			err := dbPool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'vector')").Scan(&extExists)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(extExists).Should(BeTrue())
		})

		It("should be able to create vector column", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := dbPool.Exec(ctx, "CREATE TABLE IF NOT EXISTS test_vectors (id serial primary key, embedding vector(3))")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = dbPool.Exec(ctx, "INSERT INTO test_vectors (embedding) VALUES ('[1,2,3]')")
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

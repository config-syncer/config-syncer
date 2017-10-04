package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// apiv1 "k8s.io/client-go/pkg/api/v1"
)

var _ = Describe("Book", func() {
	BeforeEach(func() {
		fmt.Println("Hello BeforeEach")
		fmt.Println("Hello BeforeEach--------------------------------")
	})
	JustBeforeEach(func() {
		fmt.Println("Hello Just Before each")
	})

	Describe("Hello 1", func() {
		It("inner help", func() {
			fmt.Println("*_*_*_*_*_", root)
			pods, err := root.KubeClient.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())

			fmt.Println("hello pods==============", pods)
			expected := "Hello"
			actual := "Hello"
			Expect(actual).Should(Equal(expected))
		})
	})
})

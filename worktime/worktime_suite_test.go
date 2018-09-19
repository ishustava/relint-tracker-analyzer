package worktime_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWorktime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Worktime Suite")
}

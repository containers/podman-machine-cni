package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPodmanMachine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PodmanMachine Suite")
}

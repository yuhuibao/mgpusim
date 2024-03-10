package emu

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination "mock_vm_test.go" -package $GOPACKAGE -write_package_comment=false github.com/sarchlab/akita/v3/mem/vm PageTable
//go:generate mockgen -source alu.go -destination "mock_instEmuState_test.go" -package $GOPACKAGE -mock_names=InstEmuState=MockInstEmuState
func TestEmulator(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "GCN3 Emulator")
}

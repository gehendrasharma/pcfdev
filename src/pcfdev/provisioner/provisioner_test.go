package provisioner_test

import (
	"bytes"
	"errors"
	"pcfdev/provisioner"
	"pcfdev/provisioner/mocks"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provisioner", func() {
	Describe("#Provision", func() {
		var (
			p             *provisioner.Provisioner
			mockCtrl      *gomock.Controller
			mockCert      *mocks.MockCert
			mockCmdRunner *mocks.MockCmdRunner
			mockFS        *mocks.MockFS
			mockUI        *mocks.MockUI
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockCert = mocks.NewMockCert(mockCtrl)
			mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)
			mockFS = mocks.NewMockFS(mockCtrl)
			mockUI = mocks.NewMockUI(mockCtrl)

			p = &provisioner.Provisioner{
				Cert:      mockCert,
				CmdRunner: mockCmdRunner,
				FS:        mockFS,
				UI:        mockUI,
			}
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("should provision a VM", func() {
			gomock.InOrder(
				mockCert.EXPECT().GenerateCert("some-domain").Return([]byte("some-cert"), []byte("some-key"), nil),
				mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
				mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
				mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
				mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain"),
				mockUI.EXPECT().PrintHelpText("some-domain"),
			)

			Expect(p.Provision("some-provision-script-path", "some-domain")).To(Succeed())
		})

		Context("when there is an error generating certificate", func() {
			It("should return the error", func() {
				mockCert.EXPECT().GenerateCert("some-domain").Return(nil, nil, errors.New("some-error"))

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error creating the gorouter config directory", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCert("some-domain").Return([]byte("some-cert"), []byte("some-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config").Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error writing the certificate", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCert("some-domain").Return([]byte("some-cert"), []byte("some-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))).Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error writing the private key", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCert("some-domain").Return([]byte("some-cert"), []byte("some-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))).Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error running the provision script", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCert("some-domain").Return([]byte("some-cert"), []byte("some-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
					mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain").Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error printing help text", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCert.EXPECT().GenerateCert("some-domain").Return([]byte("some-cert"), []byte("some-key"), nil),
					mockFS.EXPECT().Mkdir("/var/vcap/jobs/gorouter/config"),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader([]byte("some-cert"))),
					mockFS.EXPECT().Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader([]byte("some-key"))),
					mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain"),
					mockUI.EXPECT().PrintHelpText("some-domain").Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})
	})
})
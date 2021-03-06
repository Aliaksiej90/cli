package servicebroker_test

import (
	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	. "github.com/cloudfoundry/cli/cf/commands/servicebroker"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/cli/testhelpers/matchers"
)

func callListServiceBrokers(args []string, serviceBrokerRepo *testapi.FakeServiceBrokerRepo) (ui *testterm.FakeUI) {
	ui = &testterm.FakeUI{}
	config := testconfig.NewRepositoryWithDefaults()
	cmd := NewListServiceBrokers(ui, config, serviceBrokerRepo)
	testcmd.RunCommand(cmd, args, &testreq.FakeReqFactory{})

	return
}

var _ = Describe("service-brokers command", func() {
	var (
		ui                  *testterm.FakeUI
		config              core_config.Repository
		cmd                 ListServiceBrokers
		repo                *testapi.FakeServiceBrokerRepo
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		config = testconfig.NewRepositoryWithDefaults()
		repo = &testapi.FakeServiceBrokerRepo{}
		cmd = NewListServiceBrokers(ui, config, repo)
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
	})

	Describe("login requirements", func() {
		It("fails if the user is not logged in", func() {
			requirementsFactory.LoginSuccess = false
			Expect(testcmd.RunCommand(cmd, []string{}, requirementsFactory)).To(BeFalse())
		})
		It("should fail with usage when provided any arguments", func() {
			requirementsFactory.LoginSuccess = true
			Expect(testcmd.RunCommand(cmd, []string{"blahblah"}, requirementsFactory)).To(BeFalse())
			Expect(ui.FailedWithUsage).To(BeTrue())
		})
	})

	It("lists service brokers", func() {
		repo.ServiceBrokers = []models.ServiceBroker{models.ServiceBroker{
			Name: "service-broker-to-list-a",
			Guid: "service-broker-to-list-guid-a",
			Url:  "http://service-a-url.com",
		}, models.ServiceBroker{
			Name: "service-broker-to-list-b",
			Guid: "service-broker-to-list-guid-b",
			Url:  "http://service-b-url.com",
		}, models.ServiceBroker{
			Name: "service-broker-to-list-c",
			Guid: "service-broker-to-list-guid-c",
			Url:  "http://service-c-url.com",
		}}

		testcmd.RunCommand(cmd, []string{}, requirementsFactory)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Getting service brokers as", "my-user"},
			[]string{"name", "url"},
			[]string{"service-broker-to-list-a", "http://service-a-url.com"},
			[]string{"service-broker-to-list-b", "http://service-b-url.com"},
			[]string{"service-broker-to-list-c", "http://service-c-url.com"},
		))
	})

	It("lists service brokers by alphabetical order", func() {
		repo.ServiceBrokers = []models.ServiceBroker{models.ServiceBroker{
			Name: "z-service-broker-to-list",
			Guid: "z-service-broker-to-list-guid-a",
			Url:  "http://service-a-url.com",
		}, models.ServiceBroker{
			Name: "a-service-broker-to-list",
			Guid: "a-service-broker-to-list-guid-c",
			Url:  "http://service-c-url.com",
		}, models.ServiceBroker{
			Name: "fun-service-broker-to-list",
			Guid: "fun-service-broker-to-list-guid-b",
			Url:  "http://service-b-url.com",
		}, models.ServiceBroker{
			Name: "123-service-broker-to-list",
			Guid: "123-service-broker-to-list-guid-c",
			Url:  "http://service-d-url.com",
		}}

		testcmd.RunCommand(cmd, []string{}, requirementsFactory)

		Expect(ui.Outputs).To(BeInDisplayOrder(
			[]string{"Getting service brokers as", "my-user"},
			[]string{"name", "url"},
			[]string{"123-service-broker-to-list", "http://service-d-url.com"},
			[]string{"a-service-broker-to-list", "http://service-c-url.com"},
			[]string{"fun-service-broker-to-list", "http://service-b-url.com"},
			[]string{"z-service-broker-to-list", "http://service-a-url.com"},
		))
	})

	It("says when no service brokers were found", func() {
		testcmd.RunCommand(cmd, []string{}, requirementsFactory)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Getting service brokers as", "my-user"},
			[]string{"No service brokers found"},
		))
	})

	It("reports errors when listing service brokers", func() {
		repo.ListErr = true
		testcmd.RunCommand(cmd, []string{}, requirementsFactory)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Getting service brokers as ", "my-user"},
			[]string{"FAILED"},
		))
	})
})

package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestArchApplicationLayer(t *testing.T) {
	archtest.Package(t, applicationLayer).ShouldNotDependOn(
		presentationLayer,
		serviceLayer,
		persistentLayer,
	)
}

func TestApplicationLayerHaveTests(t *testing.T) {
	archtest.Package(t, applicationLayer).IncludeTests()
}

func TestDomainLayerHaveNoDependencies(t *testing.T) {
	archtest.Package(t, domainLayer).ShouldNotDependOn(
		applicationLayer,
		serviceLayer,
		persistentLayer,
		presentationLayer,
	)
}

func TestArchPresentationLayer(t *testing.T) {
	archtest.Package(t, presentationLayer).ShouldNotDependDirectlyOn(
		serviceLayer,
		applicationLayer,
		domainLayer,
	)
}

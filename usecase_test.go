package cleancoder

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/gunit"
	"testing"
)

type PresentCodeCastUsecaseFixture struct {
	*gunit.Fixture
	user     User
	codecast CodeCast
	usecase  PresentCodecastUsecase
}

func TestPresentCodeCastUsecaseFixture(t *testing.T) {
	gunit.Run(new(PresentCodeCastUsecaseFixture), t)
}

func (p *PresentCodeCastUsecaseFixture) Setup() {
	p.usecase = PresentCodecastUsecase{
		gateway: &mockGateway{},
	}
	p.user = User{ID: "u1ID", Username: "User"}
	p.codecast = CodeCast{Title: "A", PublishedDate: "3/1/2022"}
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_shouldNotViewCodecast() {
	p.So(usecase.IsLicenseToviewCodecast(&p.user, &p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithViewableLicense_canViewCodecast() {
	license := License{p.user, p.codecast}
	p.usecase.gateway.SaveLicense(&license)
	p.So(p.usecase.IsLicenseToviewCodecast(&p.user, &p.codecast), ShouldBeTrue)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_canViewOthersCodecast() {
	user2 := User{Username: "User2"}
	license := License{p.user, p.codecast}
	p.usecase.gateway.SaveLicense(&license)
	p.So(usecase.IsLicenseToviewCodecast(&user2, &p.codecast), ShouldBeFalse)
}

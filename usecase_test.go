package cleancoder

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/gunit"
	"testing"
)

type PresentCodeCastUsecaseFixture struct {
	*gunit.Fixture
	user     *User
	codecast *CodeCast
	usecase  PresentCodecastUsecase
}

func TestPresentCodeCastUsecaseFixture(t *testing.T) {
	gunit.Run(new(PresentCodeCastUsecaseFixture), t)
}

func (p *PresentCodeCastUsecaseFixture) Setup() {
	p.usecase = PresentCodecastUsecase{
		gateway: &mockGateway{},
	}
	u := User{Username: "User"}
	p.user = p.usecase.gateway.SaveUser(&u)
	p.codecast = p.usecase.gateway.Save(&CodeCast{Title: "A", PublicationDate: "3/1/2022"})
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_shouldNotViewCodecast() {
	p.So(usecase.IsLicenseToviewCodecast(p.user, p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithViewableLicense_canViewCodecast() {
	p.usecase.gateway.SaveLicense(&License{*p.user, *p.codecast})
	p.So(p.usecase.IsLicenseToviewCodecast(p.user, p.codecast), ShouldBeTrue)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_canViewOthersCodecast() {
	user2 := User{Username: "User2"}
	license := License{*p.user, *p.codecast}
	p.usecase.gateway.SaveLicense(&license)
	p.So(usecase.IsLicenseToviewCodecast(&user2, p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentingNoCodecasts() {
	_ = p.usecase.gateway.Delete(p.codecast)
	codecasts := p.usecase.PreentCodecasts(p.user)
	p.So(codecasts, ShouldBeEmpty)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentOneCodecast() {
	p.codecast.Title = "Some Title"
	p.codecast.PublicationDate = "Tomorrow"
	p.usecase.gateway.Save(p.codecast)
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.Title, ShouldEqual, "Some Title")
	p.So(pc.PublicationDate, ShouldEqual, "Tomorrow")
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsNotViewableIfNoLicense() {
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsViewableIfHaveLicense() {
	p.usecase.gateway.SaveLicense(&License{*p.user, *p.codecast})
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeTrue)
}

package cleancoder

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/gunit"
)

type PresentCodeCastUsecaseFixture struct {
	*gunit.Fixture
	user     *User
	codecast *Codecast
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
	t, _ := time.Parse("02/01/2006", "3/1/2022")
	p.codecast = p.usecase.gateway.Save(&Codecast{Title: "A", PublicationDate: t})
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_shouldNotViewCodecast() {
	p.So(usecase.IsLicenseToViewCodecast(p.user, p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithViewableLicense_canViewCodecast() {
	p.usecase.gateway.SaveLicense(&License{User: *p.user, Codecast: *p.codecast, Type: Viewing})
	p.So(p.usecase.IsLicenseToViewCodecast(p.user, p.codecast), ShouldBeTrue)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_canViewOthersCodecast() {
	user2 := User{Username: "User2"}
	license := License{User: *p.user, Codecast: *p.codecast}
	p.usecase.gateway.SaveLicense(&license)
	p.So(usecase.IsLicenseToViewCodecast(&user2, p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentingNoCodecasts() {
	_ = p.usecase.gateway.Delete(p.codecast)
	codecasts := p.usecase.PreentCodecasts(p.user)
	p.So(codecasts, ShouldBeEmpty)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentOneCodecast() {
	p.codecast.Title = "Some Title"
	p.codecast.PublicationDate = convertDate("02/01/2022")
	p.usecase.gateway.Save(p.codecast)
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.Title, ShouldEqual, "Some Title")
	p.So(pc.PublicationDate, ShouldEqual, "02/01/2022")
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsNotViewableIfNoLicense() {
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsViewableIfHaveLicense() {
	p.usecase.gateway.SaveLicense(&License{User: *p.user, Codecast: *p.codecast, Type: Viewing})
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeTrue)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsDownloadableIfHaveDownloadLicense() {
	p.usecase.gateway.SaveLicense(&License{User: *p.user, Codecast: *p.codecast, Type: Downloading})
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeFalse)
	p.So(pc.IsDownloadable, ShouldBeTrue)
}

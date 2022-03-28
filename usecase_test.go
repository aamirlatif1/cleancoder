package cleancoder

import (
	"github.com/aamirlatif1/cleancoder/data"
	"github.com/aamirlatif1/cleancoder/usecase"
	"testing"
	"time"

	"github.com/aamirlatif1/cleancoder/entity"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/gunit"
)

type PresentCodeCastUsecaseFixture struct {
	*gunit.Fixture
	user     *entity.User
	codecast *entity.Codecast
	uc       usecase.PresentCodecast
}

func TestPresentCodeCastUsecaseFixture(t *testing.T) {
	gunit.Run(new(PresentCodeCastUsecaseFixture), t)
}

func (p *PresentCodeCastUsecaseFixture) Setup() {
	p.uc = usecase.PresentCodecast{
		UserGateway:     &inMemoryUserGateway{},
		LicenseGateway:  &inMemoryLicenseGateway{},
		CodecastGateway: &inMemoryCodecastGateway{},
	}
	u := entity.User{Username: "User"}
	p.user = p.uc.UserGateway.SaveUser(&u)
	t, _ := time.Parse("02/01/2006", "3/1/2022")
	p.codecast = p.uc.CodecastGateway.Save(&entity.Codecast{Title: "A", PublicationDate: t})
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_shouldNotViewCodecast() {
	p.So(uc.IsLicenseToViewCodecast(p.user, p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithViewableLicense_canViewCodecast() {
	p.uc.LicenseGateway.SaveLicense(&entity.License{User: *p.user, Codecast: *p.codecast, Type: data.Viewing})
	p.So(p.uc.IsLicenseToViewCodecast(p.user, p.codecast), ShouldBeTrue)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_canViewOthersCodecast() {
	user2 := entity.User{Username: "User2"}
	license := entity.License{User: *p.user, Codecast: *p.codecast, Type: data.Viewing}
	p.uc.LicenseGateway.SaveLicense(&license)
	p.So(uc.IsLicenseToViewCodecast(&user2, p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentingNoCodecasts() {
	_ = p.uc.CodecastGateway.Delete(p.codecast)
	codecasts := p.uc.PreentCodecasts(p.user)
	p.So(codecasts, ShouldBeEmpty)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentOneCodecast() {
	p.codecast.Title = "Some Title"
	p.codecast.PublicationDate = convertDate("02/01/2022")
	p.uc.CodecastGateway.Save(p.codecast)
	presentableCodecasts := p.uc.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.Title, ShouldEqual, "Some Title")
	p.So(pc.PublicationDate, ShouldEqual, "02/01/2022")
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsNotViewableIfNoLicense() {
	presentableCodecasts := p.uc.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsViewableIfHaveLicense() {
	p.uc.LicenseGateway.SaveLicense(&entity.License{User: *p.user, Codecast: *p.codecast, Type: data.Viewing})
	presentableCodecasts := p.uc.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeTrue)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsDownloadableIfHaveDownloadLicense() {
	p.uc.LicenseGateway.SaveLicense(&entity.License{User: *p.user, Codecast: *p.codecast, Type: data.Downloading})
	presentableCodecasts := p.uc.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeFalse)
	p.So(pc.IsDownloadable, ShouldBeTrue)
}

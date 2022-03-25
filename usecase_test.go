package cleancoder

import (
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
	usecase  PresentCodecastUsecase
}

func TestPresentCodeCastUsecaseFixture(t *testing.T) {
	gunit.Run(new(PresentCodeCastUsecaseFixture), t)
}

func (p *PresentCodeCastUsecaseFixture) Setup() {
	p.usecase = PresentCodecastUsecase{
		userGateway:     &inMemoryUserGateway{},
		licenseGateway:  &inMemoryLicenseGateway{},
		codecastGateway: &inMemoryCodecastGateway{},
	}
	u := entity.User{Username: "User"}
	p.user = p.usecase.userGateway.SaveUser(&u)
	t, _ := time.Parse("02/01/2006", "3/1/2022")
	p.codecast = p.usecase.codecastGateway.Save(&entity.Codecast{Title: "A", PublicationDate: t})
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_shouldNotViewCodecast() {
	p.So(usecase.IsLicenseToViewCodecast(p.user, p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithViewableLicense_canViewCodecast() {
	p.usecase.licenseGateway.SaveLicense(&entity.License{User: *p.user, Codecast: *p.codecast, Type: Viewing})
	p.So(p.usecase.IsLicenseToViewCodecast(p.user, p.codecast), ShouldBeTrue)
}

func (p *PresentCodeCastUsecaseFixture) TestUserWithoutViewableLicense_canViewOthersCodecast() {
	user2 := entity.User{Username: "User2"}
	license := entity.License{User: *p.user, Codecast: *p.codecast, Type: Viewing}
	p.usecase.licenseGateway.SaveLicense(&license)
	p.So(usecase.IsLicenseToViewCodecast(&user2, p.codecast), ShouldBeFalse)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentingNoCodecasts() {
	_ = p.usecase.codecastGateway.Delete(p.codecast)
	codecasts := p.usecase.PreentCodecasts(p.user)
	p.So(codecasts, ShouldBeEmpty)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentOneCodecast() {
	p.codecast.Title = "Some Title"
	p.codecast.PublicationDate = convertDate("02/01/2022")
	p.usecase.codecastGateway.Save(p.codecast)
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
	p.usecase.licenseGateway.SaveLicense(&entity.License{User: *p.user, Codecast: *p.codecast, Type: Viewing})
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeTrue)
}

func (p *PresentCodeCastUsecaseFixture) TestPresentedCodecastIsDownloadableIfHaveDownloadLicense() {
	p.usecase.licenseGateway.SaveLicense(&entity.License{User: *p.user, Codecast: *p.codecast, Type: Downloading})
	presentableCodecasts := p.usecase.PreentCodecasts(p.user)
	p.So(len(presentableCodecasts), ShouldEqual, 1)
	pc := presentableCodecasts[0]
	p.So(pc.IsViewable, ShouldBeFalse)
	p.So(pc.IsDownloadable, ShouldBeTrue)
}

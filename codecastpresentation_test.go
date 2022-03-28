package cleancoder

import (
	"errors"
	"fmt"
	"github.com/aamirlatif1/cleancoder/data"
	"github.com/aamirlatif1/cleancoder/usecase"
	"sort"
	"testing"
	"time"

	"github.com/aamirlatif1/cleancoder/entity"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

const dateLayout = "01/02/2006"

var (
	app = application{
		userGateway:     &inMemoryUserGateway{},
		codecastGateway: &inMemoryCodecastGateway{},
		licenseGateway:  &inMemoryLicenseGateway{},
	}
	username = "U"
	uc       = usecase.PresentCodecast{
		UserGateway:     app.userGateway,
		CodecastGateway: app.codecastGateway,
		LicenseGateway:  app.licenseGateway,
	}
)

var gatekeeper Gatekeeper

func TestPresentNoCodeCasts(t *testing.T) {

	Convey("Given no codecasts", t, func() {
		codecasts := app.codecastGateway.FindAllCodecastsSortedChronologically()
		for _, c := range codecasts {
			_ = app.codecastGateway.Delete(&c)
		}
		So(app.codecastGateway.FindAllCodecastsSortedChronologically(), ShouldBeEmpty)
	})

	Convey("Given user U", t, func() {
		user := entity.User{Username: username}
		app.userGateway.SaveUser(&user)
	})

	Convey("With user U logged in", t, func() {
		user, _ := app.userGateway.FindUser(username)
		So(user, ShouldNotBeNil)
		gatekeeper.SetLoggedInUser(user)
	})

	Convey("Then the following codecasts will be presented for U", t, func() {

		So(gatekeeper.loggedInUser.Username, ShouldEqual, username)
		Convey("there will be no codecasts presented", func() {
			user, _ := app.userGateway.FindUser(username)
			codecasts := uc.PreentCodecasts(user)
			So(codecasts, ShouldBeEmpty)
		})
	})
}

func TestPrentViewableCodecasts(t *testing.T) {

	Convey("Given Codecasts", t, func() {
		records := []struct {
			title     string
			published string
		}{
			{"A", "03/01/2022"},
			{"B", "03/02/2022"},
			{"C", "02/18/2022"},
		}
		for _, record := range records {
			saveCodecast(record.title, convertDate(record.published))
		}
	})

	Convey("Given user U", t, func() {
		user := entity.User{Username: username}
		app.userGateway.SaveUser(&user)
	})

	Convey("with user U logged in", t, func() {
		user, _ := app.userGateway.FindUser(username)
		So(user, ShouldNotBeNil)
		gatekeeper.SetLoggedInUser(user)
	})

	Convey("and with license for U able to view A", t, func() {
		user, _ := app.userGateway.FindUser(username)
		codecast, _ := app.codecastGateway.FindCodecastByTitle("A")
		liccense := entity.License{
			User:     *user,
			Codecast: *codecast,
			Type:     data.Viewing,
		}
		app.licenseGateway.SaveLicense(&liccense)
		So(uc.IsLicenseToViewCodecast(user, codecast), ShouldBeTrue)
	})

	Convey("and with viewable and downloadable license for U able to view and download B", t, func() {
		user, _ := app.userGateway.FindUser(username)
		codecast, _ := app.codecastGateway.FindCodecastByTitle("B")
		app.licenseGateway.SaveLicense(&entity.License{
			User:     *user,
			Codecast: *codecast,
			Type:     data.Viewing,
		})
		app.licenseGateway.SaveLicense(&entity.License{
			User:     *user,
			Codecast: *codecast,
			Type:     data.Downloading,
		})
		So(uc.IsLicenseToViewCodecast(user, codecast), ShouldBeTrue)
	})

	Convey("then the following codecasts will be presented for U", t, func() {
		expectedPc := []data.PresentableCodecast{
			{"C", "C", "C", "02/18/2022", false, false},
			{"A", "A", "A", "03/01/2022", true, false},
			{"B", "B", "B", "03/02/2022", true, true},
		}
		actualPc := uc.PreentCodecasts(&gatekeeper.loggedInUser)
		So(actualPc, ShouldResemble, expectedPc)
	})

}

func convertDate(date string) time.Time {
	t, err := time.Parse(dateLayout, date)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

func saveCodecast(title string, publishedDate time.Time) {
	codecast := entity.Codecast{
		Title:           title,
		PublicationDate: publishedDate,
	}
	app.codecastGateway.Save(&codecast)
}

type inMemoryLicenseGateway struct {
	licenses []entity.License
}

func (m *inMemoryLicenseGateway) SaveLicense(liccense *entity.License) {
	m.licenses = append(m.licenses, *liccense)
}

func (m *inMemoryLicenseGateway) FindLicensesForUserAndCodecast(user *entity.User, codecast *entity.Codecast) []entity.License {
	var licenses = make([]entity.License, 0)
	for _, license := range m.licenses {
		if license.User.IsSame(user) && license.Codecast.IsSame(codecast) {
			licenses = append(licenses, license)
		}
	}
	return licenses
}

type inMemoryCodecastGateway struct {
	codecasts []entity.Codecast
}

func (m *inMemoryCodecastGateway) FindAllCodecastsSortedChronologically() []entity.Codecast {
	cc := make([]entity.Codecast, len(m.codecasts))
	copy(cc, m.codecasts)
	sort.SliceStable(cc, func(i, j int) bool {
		return cc[i].PublicationDate.Unix() < cc[j].PublicationDate.Unix()
	})
	return cc
}

func (m *inMemoryCodecastGateway) Delete(codecast *entity.Codecast) error {
	m.codecasts = remove(m.codecasts, *codecast)
	return nil
}

func (m *inMemoryCodecastGateway) Save(codecast *entity.Codecast) *entity.Codecast {
	if codecast.ID == "" {
		establishID(codecast)
		m.codecasts = append(m.codecasts, *codecast)
	} else {
		m.updateCodeCast(codecast)
	}
	return codecast
}

func (m *inMemoryCodecastGateway) FindCodecastByTitle(title string) (*entity.Codecast, error) {
	for _, codecast := range m.codecasts {
		if codecast.Title == title {
			return &codecast, nil
		}
	}
	return nil, errors.New("resource not found")
}

func (m *inMemoryCodecastGateway) updateCodeCast(codecast *entity.Codecast) {
	for i, cc := range m.codecasts {
		if cc.ID == codecast.ID {
			p := &m.codecasts[i]
			p.Title = codecast.Title
			p.PublicationDate = codecast.PublicationDate
			break
		}
	}
}

type inMemoryUserGateway struct {
	users []entity.User
}

func (m *inMemoryUserGateway) FindUser(username string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("resource not found")
}

func (m *inMemoryUserGateway) SaveUser(user *entity.User) *entity.User {
	establishID(user)
	m.users = append(m.users, *user)
	return user
}

func remove(codecasts []entity.Codecast, key entity.Codecast) []entity.Codecast {
	for i, v := range codecasts {
		if v == key {
			return append(codecasts[:i], codecasts[i+1:]...)
		}
	}
	return codecasts
}

func establishID(entity entity.Entity) {
	if entity.GetID() == "" {
		entity.SetID(uuid.NewString())
	}
}

package cleancoder

import (
	"errors"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

var (
	app = application{
		gateway: &mockGateway{},
	}
	username      = "U"
	codecastTitle = "A"
	usecase       = PresentCodecastUsecase{gateway: app.gateway}
)

var gatekeeper Gatekeeper

func TestPresentNoCodeCasts(t *testing.T) {

	Convey("Given no codecasts", t, func() {
		codecasts := app.gateway.FindAllCodecasts()
		for _, c := range codecasts {
			_ = app.gateway.Delete(&c)
		}
		So(app.gateway.FindAllCodecasts(), ShouldBeEmpty)
	})

	Convey("Given user U", t, func() {
		user := User{Username: username}
		app.gateway.SaveUser(&user)
	})

	Convey("With user U logged in", t, func() {
		user, _ := app.gateway.FindUser(username)
		So(user, ShouldNotBeNil)
		gatekeeper.SetLoggedInUser(user)
	})

	Convey("Then the following codecasts will be presented for U", t, func() {

		So(gatekeeper.loggedInUser.Username, ShouldEqual, username)
		Convey("there will be no codecasts presented", func() {
			user, _ := app.gateway.FindUser(username)
			codecasts := usecase.PreentCodecasts(user)
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
			{"A", "3/1/2022"},
			{"B", "3/2/2022"},
			{"C", "2/18/2022"},
		}
		for _, record := range records {
			saveCodecast(record.title, record.published)
		}
	})

	Convey("Given user U", t, func() {
		user := User{Username: username}
		app.gateway.SaveUser(&user)
	})

	Convey("with user U logged in", t, func() {
		user, _ := app.gateway.FindUser(username)
		So(user, ShouldNotBeNil)
		gatekeeper.SetLoggedInUser(user)
	})

	Convey("and with license for U able to view A", t, func() {
		user, _ := app.gateway.FindUser(username)
		codecast, _ := app.gateway.FindCodecastByTitle(codecastTitle)
		liccense := License{
			User:     *user,
			codecast: *codecast,
		}
		app.gateway.SaveLicense(&liccense)
		So(usecase.IsLicenseToviewCodecast(user, codecast), ShouldBeTrue)
	})

	Convey("then the following codecasts will be presented for U", t, func() {
		expectedPc := []PresentableCodecast{
			{"C", "C", "C", "2/18/2022", false, false},
			{"A", "A", "A", "3/1/2022", true, false},
			{"B", "B", "B", "3/2/2022", false, false},
		}
		actualPc := usecase.PreentCodecasts(&gatekeeper.loggedInUser)
		So(actualPc, ShouldResemble, expectedPc)
	})

}

func saveCodecast(title, publishedDate string) {
	codecast := CodeCast{
		Title:           title,
		PublicationDate: publishedDate,
	}
	app.gateway.Save(&codecast)
}

type mockGateway struct {
	codecasts []CodeCast
	users     []User
	licenses  []License
}

func (m *mockGateway) FindLicensesForUserAndCodecast(user *User, codecast *CodeCast) []License {
	var licenses []License
	for _, license := range m.licenses {
		if license.User.IsSame(user) && license.codecast.IsSame(codecast) {
			licenses = append(licenses, license)
		}
	}
	return licenses
}

func (m *mockGateway) SaveLicense(liccense *License) {
	m.licenses = append(m.licenses, *liccense)
}

func (m *mockGateway) FindCodecastByTitle(title string) (*CodeCast, error) {
	for _, codecast := range m.codecasts {
		if codecast.Title == title {
			return &codecast, nil
		}
	}
	return nil, errors.New("resource not found")
}

func (m *mockGateway) FindUser(username string) (*User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("resource not found")
}

func (m *mockGateway) SaveUser(user *User) *User {
	establishID(user)
	m.users = append(m.users, *user)
	return user
}

func (m *mockGateway) Save(codecast *CodeCast) *CodeCast {
	if codecast.ID == "" {
		establishID(codecast)
		m.codecasts = append(m.codecasts, *codecast)
	} else {
		m.updateCodeCast(codecast)
	}
	return codecast
}

func (m *mockGateway) updateCodeCast(codecast *CodeCast) {
	for i, cc := range m.codecasts {
		if cc.ID == codecast.ID {
			p := &m.codecasts[i]
			p.Title = codecast.Title
			p.PublicationDate = codecast.PublicationDate
			break
		}
	}
}

func (m *mockGateway) FindAllCodecasts() []CodeCast {
	cc := make([]CodeCast, len(m.codecasts))
	copy(cc, m.codecasts)
	sort.SliceStable(cc, func(i, j int) bool {
		return cc[i].PublicationDate < cc[j].PublicationDate
	})
	return cc
}

func (m *mockGateway) Delete(codecast *CodeCast) error {
	m.codecasts = remove(m.codecasts, *codecast)
	return nil
}

func remove(codecasts []CodeCast, key CodeCast) []CodeCast {
	for i, v := range codecasts {
		if v == key {
			return append(codecasts[:i], codecasts[i+1:]...)
		}
	}
	return codecasts
}

func establishID(entity Entity) {
	if entity.GetID() == "" {
		entity.SetID(uuid.NewString())
	}
}

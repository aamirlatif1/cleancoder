package cleancoder

import (
	"errors"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
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
			codecasts := usecase.PreentCodecasts()
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
		viewableCodecasts := []struct {
			title        string
			picture      string
			description  string
			viewable     bool
			downloadable bool
		}{
			{"C", "C", "C", false, false},
			{"A", "A", "A", true, false},
			{"B", "B", "B", false, false},
		}
		_ = viewableCodecasts
	})

}

func saveCodecast(title, publishedDate string) {
	codecast := CodeCast{
		Title:         title,
		PublishedDate: publishedDate,
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

func (m *mockGateway) SaveUser(user *User) {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	m.users = append(m.users, *user)

}

func (m *mockGateway) Save(codecast *CodeCast) {
	m.codecasts = append(m.codecasts, *codecast)
}

func (m *mockGateway) FindAllCodecasts() []CodeCast {
	cc := make([]CodeCast, len(m.codecasts))
	copy(cc, m.codecasts)
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

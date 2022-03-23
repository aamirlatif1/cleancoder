package cleancoder

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	app = application{
		gateway: &mockGateway{},
	}
	username = "U"
)

var gatekeeper Gatekeeper
var usecase PresentCodecastUsecase

func TestPresentNoCodeCasts(t *testing.T) {
	Convey("Given Codecasts", t, func() {
		records := []struct{
			title string
			published string
		} {
			{"A", "1/1/2022"},
			{"B", "2/1/2022"},
			{"C", "3/1/2022"},
		}
		for _, record := range records {
			saveCodecast(record.title, record.published)
		}
	})

	Convey("Given no codecasts", t, func() {
		codecasts := app.gateway.FindAllCodecasts()
		for _, c := range codecasts {
			app.gateway.Delete(&c)
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

type mockGateway struct {
	codecasts []CodeCast
	users []User
}

func (m *mockGateway) FindUser(username string) (*User, error) {
	for _, user := range m.users{
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("resource not found")
}

func (m *mockGateway) SaveUser(user *User) {
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

func saveCodecast(title, publishedDate string ) {
	codecast := CodeCast{
		Title: title,
		PublishedDate: publishedDate,
	}
	app.gateway.Save(&codecast)
}

func remove( codecasts []CodeCast,  key CodeCast) []CodeCast {
	for i, v := range codecasts {
		if v == key {
			return append(codecasts[:i], codecasts[i+1:]...)
		}
	}
	return codecasts
}




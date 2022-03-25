package cleancoder

import (
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	app = application{
		gateway: &mockGateway{},
	}
	username = "U"
	usecase  = PresentCodecastUsecase{gateway: app.gateway}
)

var gatekeeper Gatekeeper

func TestPresentNoCodeCasts(t *testing.T) {

	Convey("Given no codecasts", t, func() {
		codecasts := app.gateway.FindAllCodecastsSortedChronologically()
		for _, c := range codecasts {
			_ = app.gateway.Delete(&c)
		}
		So(app.gateway.FindAllCodecastsSortedChronologically(), ShouldBeEmpty)
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
			{"A", "03/01/2022"},
			{"B", "03/02/2022"},
			{"C", "02/18/2022"},
		}
		for _, record := range records {
			saveCodecast(record.title, convertDate(record.published))
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
		codecast, _ := app.gateway.FindCodecastByTitle("A")
		liccense := License{
			User:     *user,
			Codecast: *codecast,
			Type:     Viewing,
		}
		app.gateway.SaveLicense(&liccense)
		So(usecase.IsLicenseToViewCodecast(user, codecast), ShouldBeTrue)
	})

	Convey("and with viewable and downloadable license for U able to view and download B", t, func() {
		user, _ := app.gateway.FindUser(username)
		codecast, _ := app.gateway.FindCodecastByTitle("B")
		app.gateway.SaveLicense(&License{
			User:     *user,
			Codecast: *codecast,
			Type:     Viewing,
		})
		app.gateway.SaveLicense(&License{
			User:     *user,
			Codecast: *codecast,
			Type:     Downloading,
		})
		So(usecase.IsLicenseToViewCodecast(user, codecast), ShouldBeTrue)
	})

	Convey("then the following codecasts will be presented for U", t, func() {
		expectedPc := []PresentableCodecast{
			{"C", "C", "C", "02/18/2022", false, false},
			{"A", "A", "A", "03/01/2022", true, false},
			{"B", "B", "B", "03/02/2022", true, true},
		}
		actualPc := usecase.PreentCodecasts(&gatekeeper.loggedInUser)
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
	codecast := Codecast{
		Title:           title,
		PublicationDate: publishedDate,
	}
	app.gateway.Save(&codecast)
}

type mockGateway struct {
	codecasts []Codecast
	users     []User
	licenses  []License
}

func (m *mockGateway) FindLicensesForUserAndCodecast(user *User, codecast *Codecast) []License {
	var licenses []License
	for _, license := range m.licenses {
		if license.User.IsSame(user) && license.Codecast.IsSame(codecast) {
			licenses = append(licenses, license)
		}
	}
	return licenses
}

func (m *mockGateway) SaveLicense(liccense *License) {
	m.licenses = append(m.licenses, *liccense)
}

func (m *mockGateway) FindCodecastByTitle(title string) (*Codecast, error) {
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

func (m *mockGateway) Save(codecast *Codecast) *Codecast {
	if codecast.ID == "" {
		establishID(codecast)
		m.codecasts = append(m.codecasts, *codecast)
	} else {
		m.updateCodeCast(codecast)
	}
	return codecast
}

func (m *mockGateway) updateCodeCast(codecast *Codecast) {
	for i, cc := range m.codecasts {
		if cc.ID == codecast.ID {
			p := &m.codecasts[i]
			p.Title = codecast.Title
			p.PublicationDate = codecast.PublicationDate
			break
		}
	}
}

func (m *mockGateway) FindAllCodecastsSortedChronologically() []Codecast {
	cc := make([]Codecast, len(m.codecasts))
	copy(cc, m.codecasts)
	sort.SliceStable(cc, func(i, j int) bool {
		return cc[i].PublicationDate.Unix() < cc[j].PublicationDate.Unix()
	})
	return cc
}

func (m *mockGateway) Delete(codecast *Codecast) error {
	m.codecasts = remove(m.codecasts, *codecast)
	return nil
}

func remove(codecasts []Codecast, key Codecast) []Codecast {
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

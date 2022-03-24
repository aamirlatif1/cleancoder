package cleancoder

import (
	. "github.com/smartystreets/assertions"
	"github.com/smartystreets/gunit"
	"testing"
)

func TestUserFixture(t *testing.T) {
	gunit.Run(new(UserFixture), t)
}

type UserFixture struct {
	*gunit.Fixture
}

func (u *UserFixture) Setup() {
}

func (u *UserFixture) TestTwoDifferentUserAreNotSame() {
	u1 := User{ID: "u1ID", Username: "u1"}
	u2 := User{ID: "u2ID", Username: "u2"}
	u.So(u1.IsSame(&u2), ShouldBeFalse)
}

func (u *UserFixture) TestOneUserIsSameAsItself() {
	u1 := User{Username: "u1"}
	u.So(u1.IsSame(&u1), ShouldBeTrue)
}

func (u *UserFixture) TestUserWithSameIdAreTheSame() {
	u1 := User{Username: "u1"}
	u2 := User{Username: "u2"}
	u1.ID = "u1ID"
	u2.ID = "u1ID"
	u.So(u1.IsSame(&u2), ShouldBeTrue)
}

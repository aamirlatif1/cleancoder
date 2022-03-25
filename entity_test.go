package cleancoder

import (
	"testing"

	. "github.com/smartystreets/assertions"
	"github.com/smartystreets/gunit"
)

func TestUserFixture(t *testing.T) {
	gunit.Run(new(EntityFixture), t)
}

type EntityFixture struct {
	*gunit.Fixture
}

func (u *EntityFixture) Setup() {
}

func (u *EntityFixture) TestTwoDifferentEntitiesAreNotSame() {
	e1 := newEntity()
	e2 := newEntity()
	e1.SetID("e1ID")
	e2.SetID("e2ID")
	u.So(e1.IsSame(e2), ShouldBeFalse)
}

func (u *EntityFixture) TestOneEntityIsSameAsItself() {
	e1 := newEntity()
	e1.SetID("e1ID")
	u.So(e1.IsSame(e1), ShouldBeTrue)
}

func (u *EntityFixture) TestEntityWithSameIdAreTheSame() {
	e1 := newEntity()
	e2 := newEntity()
	e1.SetID("e1ID")
	e2.SetID("e1ID")
	u.So(e1.IsSame(e2), ShouldBeTrue)
}

func (u *EntityFixture) TestEntityWithNullIdsAreNeverSame() {
	e1 := newEntity()
	e2 := newEntity()
	u.So(e1.IsSame(e2), ShouldBeFalse)
}

func newEntity() Entity {
	return &Codecast{}
}

package dogowners

import (
	"testing"
	"time"

	"github.com/almendar/golang-gorm-chi-postgres/shared"
	"github.com/google/go-cmp/cmp"
)

func TestStorage(t *testing.T) {
	gorm, err := shared.DefaultGormHandle()
	if err != nil {
		t.Fatal(err)
	}
	RunMigration(gorm)

	owner := &OwnerDBModel{
		Name:     "John",
		Email:    nil,
		Age:      20,
		Birthday: time.Now().Truncate(24 * time.Hour).Add(-20 * 365 * 24 * time.Hour),
		Dogs:     make([]DogDBModel, 0),
	}

	dogs := []DogDBModel{
		{
			Name:     "Dog1",
			Birthday: time.Now().Truncate(24 * time.Hour).Add(-5 * 365 * 24 * time.Hour),
		},
		{
			Name:     "Dog2",
			Birthday: time.Now().Truncate(24 * time.Hour).Add(-3 * 365 * 24 * time.Hour),
		},
	}

	t.Run("SaveOwner", func(t *testing.T) {
		storage := NewDatabase(gorm)

		err := storage.SaveOwner(owner)
		if err != nil {
			t.Fatal(err)
		}
		if owner.ID == 0 {
			t.Fatal("ID is not set")
		}

		owner2, err := storage.Owner(owner.ID)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(owner, owner2); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("SaveDog", func(t *testing.T) {
		storage := NewDatabase(gorm)

		ownerCopy := owner
		ownerCopy.Dogs = dogs

		for i := range dogs {
			dogs[i].OwnerID = ownerCopy.ID
			err := storage.SaveDog(&dogs[i])
			if err != nil {
				t.Fatal(err)
			}
			if dogs[i].ID == 0 {
				t.Fatal("ID is not set")
			}
		}

		owner2, err := storage.Owner(ownerCopy.ID)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(ownerCopy, owner2); diff != "" {
			t.Fatal(diff)
		}
	})
}

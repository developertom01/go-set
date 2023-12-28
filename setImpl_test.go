package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	set := NewSet()

	set.Add(2)
	assert.True(t, set.Has(2))
	assert.Equal(t, 1, set.Len())
	assert.False(t, set.IsEmpty())

	set.Remove(2)
	assert.True(t, set.IsEmpty())

	set.Add(1)
	set.Add(1)
	set.Add(1)
	assert.Equal(t, 1, set.Len())
	set.Remove(1)
	assert.True(t, set.IsEmpty())

	set.Add(3)
	set.Add(4)
	set.Add(5)
	assert.Equal(t, 3, set.Len())
	set2 := NewSet()
	set2.Add(3)
	set2.Add(5)
	assert.Equal(t, 2, set2.Len())
	assert.True(t, set.Contains(set2))

	set2.Add(8)
	assert.False(t, set.Has(set2))

}

func TestSetToList(t *testing.T) {
	set := NewSet()

	set.Add(1)
	set.Add(2)
	set.Add(3)

	list := set.ToSlice()
	assert.Equal(t, 3, len(list))

}

func TestSetUnionOperations(t *testing.T) {
	var (
		set1, set2 = NewSet(), NewSet()
	)
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2.Add(4)
	set2.Add(5)
	set2.Add(6)

	unionSet := set1.Union(set2)
	assert.Equal(t, 6, unionSet.Len())
	assert.True(t, unionSet.Contains(set1))
	assert.True(t, unionSet.Contains(set2))
}

func TestSetIntersectionOperations(t *testing.T) {
	var (
		set1, set2 = NewSet(), NewSet()
	)
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2.Add(4)
	set2.Add(5)
	set2.Add(6)

	intersectionSet := set1.Intersection(set2)
	assert.Equal(t, 0, intersectionSet.Len())

	set2.Add(2)
	set2.Add(3)
	intersectionSet = set1.Intersection(set2)
	assert.Equal(t, 2, intersectionSet.Len())
	assert.True(t, set1.Contains(intersectionSet))
	assert.True(t, set2.Contains(intersectionSet))

}

func TestSetComplementOperations(t *testing.T) {
	var (
		set1, set2 = NewSet(), NewSet()
	)
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2.Add(1)
	set2.Add(2)
	set2.Add(3)

	complementSet := set1.Complement(set2)
	assert.Equal(t, 0, complementSet.Len())
	set2.Add(4)
	set2.Add(5)
	complementSet = set1.Complement(set2)
	assert.Equal(t, 2, complementSet.Len())
	assert.False(t, set1.Contains(complementSet))

}

# Go implementation for set with support for concurrent access

## Example
``` go
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
```
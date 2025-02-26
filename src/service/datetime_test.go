package service

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatesNormalizedDateTime(t *testing.T) {
	for _, x := range []struct {
		date Date
		time Time
	}{
		{Ɀ_Date_(1000, 7, 15), Ɀ_Time_(15, 00)},
		{Ɀ_Date_(1000, 7, 14), Ɀ_TimeTomorrow_(15, 00)},
		{Ɀ_Date_(1000, 7, 16), Ɀ_TimeYesterday_(15, 00)},
	} {
		dt := NewDateTime(x.date, x.time)
		assert.Equal(t, "1000-07-15", dt.Date.ToString())
		assert.Equal(t, "15:00", dt.Time.ToString())
	}
}

func TestEqualsDateTime(t *testing.T) {
	dt1 := NewDateTime(Ɀ_Date_(1000, 7, 15), Ɀ_Time_(12, 00))
	dt2 := NewDateTime(Ɀ_Date_(1000, 7, 15), Ɀ_Time_(12, 01))
	dt3 := NewDateTime(Ɀ_Date_(1000, 7, 16), Ɀ_Time_(12, 00))
	assert.True(t, dt1.IsEqual(dt1))
	assert.False(t, dt1.IsEqual(dt2))
	assert.False(t, dt1.IsEqual(dt3))
	assert.False(t, dt2.IsEqual(dt3))
}

func TestIsAfterOrEqualsDateTime(t *testing.T) {
	dt1 := NewDateTime(Ɀ_Date_(1000, 7, 14), Ɀ_Time_(13, 00))
	dt2 := NewDateTime(Ɀ_Date_(1000, 7, 15), Ɀ_Time_(11, 59))
	dt3 := NewDateTime(Ɀ_Date_(1000, 7, 15), Ɀ_Time_(12, 00))
	dt4 := NewDateTime(Ɀ_Date_(1000, 7, 15), Ɀ_Time_(12, 01))
	dt5 := NewDateTime(Ɀ_Date_(1000, 7, 16), Ɀ_Time_(11, 01))

	assert.True(t, dt2.IsAfterOrEqual(dt1))
	assert.True(t, dt3.IsAfterOrEqual(dt2))
	assert.True(t, dt4.IsAfterOrEqual(dt3))
	assert.True(t, dt5.IsAfterOrEqual(dt4))

	assert.True(t, dt5.IsAfterOrEqual(dt1))
	assert.True(t, dt5.IsAfterOrEqual(dt1))

	assert.False(t, dt1.IsAfterOrEqual(dt2))
	assert.False(t, dt1.IsAfterOrEqual(dt3))
	assert.False(t, dt1.IsAfterOrEqual(dt5))
	assert.False(t, dt2.IsAfterOrEqual(dt3))
}

package datetime

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// "klog/testutil" REMINDER: can not use `testutil` because of circular import
	"testing"
)

func TestCreateTimeRange(t *testing.T) {
	time1, _ := NewTime(11, 25)
	time2, _ := NewTime(17, 10)
	tr, err := NewTimeRange(time1, time2)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, tr.Start(), time1)
	assert.Equal(t, tr.End(), time2)
}

func TestCreateOverlappingTimeRange(t *testing.T) {
	time1, _ := NewTime(23, 30)
	time2, _ := NewTime(17, 10)
	tr, err := NewOverlappingTimeRange(time1, true, time2, false)
	require.Nil(t, err)
	require.NotNil(t, tr)
	assert.Equal(t, tr.Start(), time1)
	assert.Equal(t, tr.End(), time2)
	assert.Equal(t, tr.IsStartYesterday(), true)
	assert.Equal(t, tr.IsEndTomorrow(), false)
}

func TestCreationFailsIfStartIsBeforeEnd(t *testing.T) {
	time1, _ := NewTime(14, 00)
	time2, _ := NewTime(13, 59)
	tr, err := NewTimeRange(time1, time2)
	assert.Nil(t, tr)
	assert.Equal(t, errors.New("ILLEGAL_RANGE"), err)
}

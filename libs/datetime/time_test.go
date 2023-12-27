package datetime_test

import (
	"errors"
	"testing"
	"time"

	"github.com/mauricetjmurphy/ms-common/libs/datetime"
	"github.com/stretchr/testify/assert"
)

func TestDateConverter(t *testing.T) {
	var (
		tests = []struct {
			Name    string
			Layout  string
			Date    string
			Want    string
			WantErr error
		}{
			{
				Name:   "Test with Date format dd/mm/yy",
				Layout: datetime.DateDMYSlash,
				Date:   "31/07/21",
				Want:   "2021-07-31 00:00:00 +0000 UTC",
			},
			{
				Name:   "Test with Date format mm/dd/yy",
				Layout: datetime.DateMDYSlash,
				Date:   "06/30/21",
				Want:   "2021-06-30 00:00:00 +0000 UTC",
			},
			{
				Name:   "Test with Datetime format mm/dd/yy hh:mm",
				Date:   "10/29/21 10:10",
				Layout: datetime.DateMDYTimeSlash,
				Want:   "2021-10-29 10:10:00 +0000 UTC",
			},
			{
				Name:    "Test with empty Date input",
				Date:    "",
				Layout:  datetime.DateMDYSlash,
				WantErr: errors.New("requires value [] and layout [01/02/06]"),
			},
			{
				Name:    "Test with empty Layout",
				Date:    "06/30/21",
				Layout:  "",
				WantErr: errors.New("requires value [06/30/21] and layout []"),
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			got, err := datetime.NewTime(tt.Layout, tt.Date)

			if err != nil && err.Error() != tt.WantErr.Error() {
				t.Errorf("layout (%v): date %v, error %v and want %v", tt.Layout, tt.Date, err, tt.WantErr)
			}

			if tt.Want != "" && got.String() != tt.Want {
				t.Errorf("Time = %v, want %v", got.String(), tt.Want)
			}
		})
	}

}

func Test_DateTime_Sort(t *testing.T) {
	//Given
	var (
		t1 = time.Date(2022, 01, 01, 07, 0, 0, 0, time.UTC)
		t2 = time.Date(2022, 01, 01, 07, 0, 0, 0, time.UTC)
		t3 = time.Date(2022, 01, 02, 07, 0, 0, 0, time.UTC)
		t4 = time.Date(2021, 01, 01, 00, 0, 0, 0, time.UTC)
	)

	var ex1 = []time.Time{t1, t2, t3, t4}

	//Then
	datetime.Sort(ex1)

	//Then
	assert.Len(t, ex1, 4)
	assert.Equal(t, t4, ex1[0])
	assert.Equal(t, t1, ex1[1])
	assert.Equal(t, t2, ex1[2])
	assert.Equal(t, t3, ex1[3])
}

package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Mayuko-Matsumoto/rikon-kaiketsukin-tracker/backend/internal/domain"
)

func makeHolidayMap(dates ...string) domain.HolidayChecker {
	set := make(map[string]struct{}, len(dates))
	for _, d := range dates {
		set[d] = struct{}{}
	}
	return func(date time.Time) bool {
		_, ok := set[date.Format("2006-01-02")]
		return ok
	}
}

func TestCalculateDueDate(t *testing.T) {
	t.Parallel()

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	tests := []struct {
		name      string
		year      int
		month     time.Month
		loc       *time.Location
		isHoliday domain.HolidayChecker
		wantDate  string
		wantErr   error
	}{
		{
			name:     "通常月末平日: 2026年11月30日(月)",
			year:     2026,
			month:    time.November,
			loc:      jst,
			wantDate: "2026-11-30",
		},
		{
			name:     "月末が日曜: 2026年5月31日(日)は29日(金)へ前倒し",
			year:     2026,
			month:    time.May,
			loc:      jst,
			wantDate: "2026-05-29",
		},
		{
			name:     "月末が土曜: 2026年10月31日(土)は30日(金)へ前倒し",
			year:     2026,
			month:    time.October,
			loc:      jst,
			wantDate: "2026-10-30",
		},
		{
			name:     "閏年の2月: 2028年2月29日(火)",
			year:     2028,
			month:    time.February,
			loc:      jst,
			wantDate: "2028-02-29",
		},
		{
			name:     "平年の2月(月末土曜): 2026年2月28日(土)は27日(金)へ前倒し",
			year:     2026,
			month:    time.February,
			loc:      jst,
			wantDate: "2026-02-27",
		},
		{
			name:      "月末金曜が祝日: 2026年7月31日(金)を祝日注入すると1日前倒しで30日(木)",
			year:      2026,
			month:     time.July,
			loc:       jst,
			isHoliday: makeHolidayMap("2026-07-31"),
			wantDate:  "2026-07-30",
		},
		{
			name:      "月末日曜かつ直前金曜が祝日: 2026年5月31日(日)と29日(金)祝日で28日(木)へ前倒し",
			year:      2026,
			month:     time.May,
			loc:       jst,
			isHoliday: makeHolidayMap("2026-05-29"),
			wantDate:  "2026-05-28",
		},
		{
			name:     "最終回(令和9年10月): 2027年10月31日(日)は29日(金)へ前倒し",
			year:     2027,
			month:    time.October,
			loc:      jst,
			wantDate: "2027-10-29",
		},
		{
			name:     "12月標準ルール: 2026年12月31日(木)",
			year:     2026,
			month:    time.December,
			loc:      jst,
			wantDate: "2026-12-31",
		},
		{
			name:     "タイムゾーンnil時はJSTデフォルト",
			year:     2026,
			month:    time.March,
			loc:      nil,
			wantDate: "2026-03-31",
		},
		{
			name:    "異常系: 月が0(1月未満)",
			year:    2026,
			month:   time.Month(0),
			loc:     jst,
			wantErr: domain.ErrInvalidMonth,
		},
		{
			name:    "異常系: 月が13(12月超過)",
			year:    2026,
			month:   time.Month(13),
			loc:     jst,
			wantErr: domain.ErrInvalidMonth,
		},
		{
			name:    "異常系: 年が0以下",
			year:    0,
			month:   time.November,
			loc:     jst,
			wantErr: domain.ErrInvalidYear,
		},
		{
			name:  "異常系: 全日が祝日の場合は無限ループせずエラー",
			year:  2026,
			month: time.November,
			loc:   jst,
			isHoliday: func(time.Time) bool {
				return true
			},
			wantErr: domain.ErrNoBusinessDay,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.CalculateDueDate(tt.year, tt.month, tt.loc, tt.isHoliday)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if !got.IsZero() {
					t.Fatalf("expected zero time on error, got %v", got)
				}
				return
			}

			if got.Format("2006-01-02") != tt.wantDate {
				t.Errorf("got = %s, want = %s", got.Format("2006-01-02"), tt.wantDate)
			}
			if tt.loc == nil {
				name, offset := got.Zone()
				if name != "Asia/Tokyo" || offset != 9*60*60 {
					t.Errorf("default zone = (%s, %d), want (Asia/Tokyo, %d)", name, offset, 9*60*60)
				}
			}
		})
	}
}

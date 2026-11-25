package domain

import (
	"errors"
	"time"
)

// HolidayChecker は指定された日付が日本の祝日かどうかを判定する関数型。
type HolidayChecker func(date time.Time) bool

var (
	// ErrInvalidYear は 1 未満の不正な年が指定された場合のエラー。
	ErrInvalidYear = errors.New("invalid year: must be >= 1")
	// ErrInvalidMonth は 1月〜12月 以外の月が指定された場合のエラー。
	ErrInvalidMonth = errors.New("invalid month: must be between January and December")
	// ErrNoBusinessDay は対象月内に営業日が1日も見つからなかった場合のエラー（無限ループ防止）。
	ErrNoBusinessDay = errors.New("no business day found in target month")
)

// CalculateDueDate は指定した年・月の支払期日（月末営業日）を算出する。
// 月末日が土曜日・日曜日・祝日（isHoliday が true）の場合は、直前の平日に前倒しする。
// loc が nil の場合は JST（Asia/Tokyo, UTC+9）を使用する。
func CalculateDueDate(year int, month time.Month, loc *time.Location, isHoliday HolidayChecker) (time.Time, error) {
	// 1. バリデーション：不正な年・月を弾く
	if year < 1 {
		return time.Time{}, ErrInvalidYear
	}
	if month < time.January || month > time.December {
		return time.Time{}, ErrInvalidMonth
	}

	// 2. タイムゾーンのデフォルト設定：loc が nil の場合は JST（UTC+9）をセットする
	if loc == nil {
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}

	// 3. 対象年・月の「暦上の最終日（月末日）」を生成する
	dueDate := time.Date(year, month+1, 0, 0, 0, 0, 0, loc)

	// 4. ループ処理：判定日が「土日」または「祝日」である間、1日ずつ前倒しする
	for {
		weekday := dueDate.Weekday()
		isWeekend := weekday == time.Saturday || weekday == time.Sunday
		holiday := isHoliday != nil && isHoliday(dueDate)

		if !isWeekend && !holiday {
			break
		}

		dueDate = dueDate.AddDate(0, 0, -1)
		if dueDate.Month() != month {
			return time.Time{}, ErrNoBusinessDay
		}
	}

	// 5. 確定した支払期日と nil（エラーなし）を返す
	return dueDate, nil
}

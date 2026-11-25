package domain

import (
	"errors"
	"time"
)

// HolidayChecker は指定された日付が日本の祝日かどうかを判定する関数型。
type HolidayChecker func(date time.Time) bool

var (
	// ErrInvalidMonth は 1月〜12月 以外の月が指定された場合のエラー。
	ErrInvalidMonth = errors.New("invalid month: must be between January and December")
	// ErrNoBusinessDay は対象月内に営業日が1日も見つからなかった場合のエラー（無限ループ防止）。
	ErrNoBusinessDay = errors.New("no business day found in target month")
)

// CalculateDueDate は指定した年・月の支払期日（月末営業日）を算出する。
// 月末日が土曜日・日曜日・祝日（isHoliday が true）の場合は、直前の平日に前倒しする。
// loc が nil の場合は JST（Asia/Tokyo, UTC+9）を使用する。
func CalculateDueDate(year int, month time.Month, loc *time.Location, isHoliday HolidayChecker) (time.Time, error) {
	// 1. バリデーション：month が time.January 未満 または time.December より大きい場合
	//    ゼロ値の time.Time{} と ErrInvalidMonth を返す

	// 2. タイムゾーンのデフォルト設定：loc が nil の場合
	//    time.FixedZone("Asia/Tokyo", 9*60*60) を loc にセットする

	// 3. 対象年・月の「暦上の最終日（月末日）」を生成する
	//    ヒント: time.Date(year, month+1, 0, 0, 0, 0, 0, loc)

	// 4. ループ処理：判定日が「土日」または「祝日」である間、1日ずつ前倒しする
	//    4-1. 判定日の曜日（Weekday()）が 土曜(time.Saturday) または 日曜(time.Sunday) か判定
	//    4-2. isHoliday が nil でなく、かつ isHoliday(判定日) が true か判定
	//    4-3. 土日でも祝日でもなければ（平日なら）ループを抜ける
	//    4-4. 判定日を 1日前（AddDate(0, 0, -1)）に更新する
	//    4-5. もし前倒しした結果、月が変わってしまったら（判定日.Month() != month）
	//         無限ループ・前月跨ぎ防止のため time.Time{} と ErrNoBusinessDay を返す

	// 5. 確定した支払期日と nil（エラーなし）を返す
}

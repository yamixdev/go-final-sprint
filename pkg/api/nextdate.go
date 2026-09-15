package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.Format(dateFormat) > now.Format(dateFormat)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("repeat rule is empty")
	}

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid yearly repeat rule")
		}

		for {
			date = date.AddDate(1, 0, 0)

			if afterNow(date, now) {
				break
			}

		}
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid daily rule")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if interval < 1 || interval > 400 {
			return "", errors.New("invalid day interval")
		}

		for {
			date = date.AddDate(0, 0, interval)

			if afterNow(date, now) {
				break
			}
		}
	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid weekly rule")
		}

		var weekdays [8]bool

		for _, value := range strings.Split(parts[1], ",") {
			day, err := strconv.Atoi(value)
			if err != nil {
				return "", err
			}

			if day < 1 || day > 7 {
				return "", errors.New("invalid weekday")
			}

			weekdays[day] = true
		}

		// Если исходная дата уже прошла, начинаем искать от сегодняшней.
		if !afterNow(date, now) {
			date, _ = time.Parse(dateFormat, now.Format(dateFormat))
		}

		for {
			date = date.AddDate(0, 0, 1)

			weekday := int(date.Weekday())

			// В Go воскресенье = 0, а по заданию воскресенье = 7.
			if weekday == 0 {
				weekday = 7
			}

			if weekdays[weekday] {
				break
			}
		}
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid monthly rule")
		}

		var days [32]bool
		var months [13]bool

		lastDay := false
		penultimateDay := false

		// Разбираем дни месяца.
		for _, value := range strings.Split(parts[1], ",") {
			day, err := strconv.Atoi(value)
			if err != nil {
				return "", err
			}

			switch {
			case day >= 1 && day <= 31:
				days[day] = true
			case day == -1:
				lastDay = true
			case day == -2:
				penultimateDay = true
			default:
				return "", errors.New("invalid month day")
			}
		}

		// Если месяцы указаны — разрешаем только их.
		if len(parts) == 3 {
			for _, value := range strings.Split(parts[2], ",") {
				month, err := strconv.Atoi(value)
				if err != nil {
					return "", err
				}

				if month < 1 || month > 12 {
					return "", errors.New("invalid month")
				}

				months[month] = true
			}
		} else {
			// Месяцы не указаны — подходят все.
			for month := 1; month <= 12; month++ {
				months[month] = true
			}
		}

		// Если исходная дата уже прошла, начинаем поиск от now.
		if !afterNow(date, now) {
			date, _ = time.Parse(dateFormat, now.Format(dateFormat))
		}

		// Ищем следующую подходящую дату.
		for i := 0; i < 366*10; i++ {
			date = date.AddDate(0, 0, 1)

			if !months[int(date.Month())] {
				continue
			}

			day := date.Day()

			// День 0 следующего месяца — последний день текущего.
			last := time.Date(
				date.Year(),
				date.Month()+1,
				0,
				0, 0, 0, 0,
				date.Location(),
			).Day()

			if days[day] ||
				(lastDay && day == last) ||
				(penultimateDay && day == last-1) {
				return date.Format(dateFormat), nil
			}
		}

		return "", errors.New("next date not found")

	default:
		return "", errors.New("unsupported repeat rule")
	}

	return date.Format(dateFormat), nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowString := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowString == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}

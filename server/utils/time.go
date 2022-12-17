package utils 

import "time"

func GetDateOfBirth(year, month, day int) time.Time {
    dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
    return dob
}

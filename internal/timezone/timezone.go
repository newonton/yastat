package timezone

import "time"

var MoscowZone = time.FixedZone("Europe/Moscow", int((time.Hour * 3).Seconds()))

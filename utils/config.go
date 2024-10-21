package utils

import (
	"flag"
)

type BaseConfig struct {
	Count      int
	Maxium     int
	IsUnique   bool
	IsStored   bool
	Recall     bool
	HintLevel  int
	DBFilename string
}

func InitConfig(Config *BaseConfig) {

	flag.IntVar(&Config.Count, "c", 30, "count of random numbers.")
	flag.IntVar(&Config.Maxium, "m", 100, "maxium for random numbers.")
	flag.BoolVar(&Config.IsUnique, "u", true, "all generated numbers will be unique.")
	flag.BoolVar(&Config.Recall, "r", false, "recall test mode.")
	flag.BoolVar(&Config.IsStored, "s", true, "store generated numbers for recall tests.")
	flag.IntVar(&Config.HintLevel, "hint", 2, "show hints when recall tests fail: 0 for no hint, 1 for diff hint, 2 for full hint")
	flag.StringVar(&Config.DBFilename, "db", "data.db", "database filename")

	flag.Parse()
}

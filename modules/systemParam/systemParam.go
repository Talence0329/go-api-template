package systemparam

import (
	"backend/basic/database"
	"fmt"
	"time"
)

var cfg = Config{
	RefreshTime: time.Minute * 5,
}
var systemParamExport = SystemParamExport{
	SystemParam: map[SystemParamKey]SystemParamData{},
}

func Init(initCfg Config) error {
	cfg = initCfg
	if err := ReadAll(); err != nil {
		return err
	}

	return nil
}

func (key SystemParamKey) Get() string {
	systemParamExport.mux.RLock()
	if value, isExist := systemParamExport.SystemParam[key]; !isExist {
		systemParamExport.mux.RUnlock()
		return ""
	} else {
		if time.Since(value.LastTime) > cfg.RefreshTime {
			systemParamExport.mux.RUnlock()
			if err := key.Read(); err != nil {
				fmt.Println(err)
			}
		} else {
			defer systemParamExport.mux.RUnlock()
		}
	}

	return systemParamExport.SystemParam[key].Value
}

func (key SystemParamKey) Read() error {
	systemParamExport.mux.Lock()
	defer systemParamExport.mux.Unlock()
	if db, err := database.BACKSTAGE.DB(); err != nil {
		return err
	} else {
		if data, err := getOne(db, string(key)); err != nil {
			return err
		} else {
			data.LastTime = time.Now()
			systemParamExport.SystemParam[key] = data
		}
	}

	return nil
}

func ReadAll() error {
	systemParamExport.mux.Lock()
	defer systemParamExport.mux.Unlock()
	if db, err := database.BACKSTAGE.DB(); err != nil {
		return err
	} else {
		if valueList, err := getAll(db); err != nil {
			return err
		} else {
			for _, data := range valueList {
				data.LastTime = time.Now()
				systemParamExport.SystemParam[data.Key] = data
			}
		}
	}

	return nil
}

// getParamListLack : 檢查aList缺少bList哪些內容
func getParamListLack(aList []SystemParamKey, bList []SystemParamKey) []SystemParamKey {
	diff := []SystemParamKey{}

	for _, bStr := range bList {
		isExist := false
		for _, aStr := range aList {
			if bStr == aStr {
				isExist = true
			}
		}
		if !isExist {
			diff = append(diff, bStr)
		}
	}

	return diff
}

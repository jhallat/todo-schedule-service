package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/jhallat/todo-schedule-service/database"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var scheduleMap = struct {
	sync.RWMutex
	m map[int]WeeklySchedule
}{m: make(map[int]WeeklySchedule)}

func init() {
	fmt.Println("loading weekly schedule...")
	sMap, err := loadScheduleMap()
	scheduleMap.m = sMap
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d schedules loaded ... \n", len(scheduleMap.m))
}

func loadScheduleMap() (map[int]WeeklySchedule, error) {
	fileName := "weekly-schedule.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}
	file, _ := ioutil.ReadFile(fileName)
	scheduleList := make([]WeeklySchedule, 0)
	err = json.Unmarshal([]byte(file), &scheduleList)
	if err != nil {
		log.Fatal(err)
	}
	sMap := make(map[int]WeeklySchedule)
	for i := 0; i < len(scheduleList); i++ {
		sMap[scheduleList[i].Id] = scheduleList[i]
	}
	return sMap, nil
}

func getSchedule(id int) *WeeklySchedule {
	row := database.DbConnection.QueryRow(`SELECT id, 
       description, 
       days as selectedDays 
       FROM weeklyschedule
       WHERE id = ?`)
	scheduleMap.RLock()
	defer scheduleMap.RUnlock()
	if schedule, ok := scheduleMap.m[id]; ok {
		return &schedule
	}
	return nil
}

func removeSchedule(id int) {
	scheduleMap.Lock()
	defer scheduleMap.Unlock()
	delete(scheduleMap.m, id)
}

func getScheduleList() ([]WeeklySchedule, error) {
	results, err := database.DbConnection.Query(`SELECT id, 
       description, 
       days as selectedDays 
       FROM weeklyschedule`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	schedules := make([]WeeklySchedule, 0)
	for results.Next() {
		var schedule WeeklySchedule
		results.Scan(&schedule.Id,
			         &schedule.Description,
			         &schedule.SelectedDays)
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func getScheduleIds() []int {
	scheduleMap.RLock()
	scheduleIds := []int{}
	for key := range scheduleMap.m {
		scheduleIds = append(scheduleIds, key)
	}
	scheduleMap.RUnlock()
	sort.Ints(scheduleIds)
	return scheduleIds
}

func getNextScheduleId() int {
	scheduleIds := getScheduleIds()
	return scheduleIds[len(scheduleIds)-1]+1
}

func addOrUpdateSchedule(schedule WeeklySchedule) (int, error) {
	addOrUpdateId := -1
	if schedule.Id > 0 {
		oldSchedule := getSchedule(schedule.Id)
		if oldSchedule == nil {
			return 0, fmt.Errorf("schedule id [%d] doesn't exist", schedule.Id)
		}
		addOrUpdateId = schedule.Id
	} else {
		addOrUpdateId = getNextScheduleId()
		schedule.Id = addOrUpdateId
	}
	scheduleMap.Lock()
	scheduleMap.m[addOrUpdateId] = schedule
	scheduleMap.Unlock()
	return addOrUpdateId, nil
}
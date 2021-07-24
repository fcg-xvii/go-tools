package json

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestJSON(t *testing.T) {
	m := Map{
		"jsrc": []int{1, 2, 3, 4},
		"kate": nil,
		"m1":   Map{"one": 1},
		"m2":   map[string]interface{}{"one": 2},
	}
	t.Log(m, string(m.JSON()))
	t.Log(string(m.ValueJSON("jsrc", []byte{})))
	t.Log(string(m.ValueJSON("jsrc1", []byte("{ }"))))
	t.Log(m.Map("m1", Map{}))
	t.Log(m.Map("m2", Map{}))
}

func TestInterface(t *testing.T) {
	m := MapFromInterface(Map{"one": 1})
	log.Println(m)
}

type TimeInterval struct {
	Start  int
	Finish int
}

func (s *TimeInterval) JSONDecode(dec *JSONDecoder) (err error) {
	var sl []int
	if err = dec.Decode(&sl); err == nil {
		if len(sl) != 2 {
			return errors.New("TimeInterval: expected 2 int values")
		}
		s.Start, s.Finish = sl[0], sl[1]
	}
	return
}

type TimeIntervals []TimeInterval

func (s *TimeIntervals) JSONDecode(dec *JSONDecoder) (isNil bool, err error) {
	var sl [][]int
	if err = dec.Decode(&sl); err == nil && sl != nil {
		*s = make(TimeIntervals, 0, len(sl))
		for i, interval := range sl {
			if len(interval) != 2 {
				return true, fmt.Errorf("Index %v - expected 2 elements list", i)
			}
			*s = append(*s, TimeInterval{
				Start:  interval[0],
				Finish: interval[1],
			})
		}
	}
	isNil = sl == nil
	return
}

type NestedObject struct {
	ID        int
	Name      string
	Embedded  *TObject
	Intervals TimeIntervals
}

type TObject struct {
	ID        int
	Name      string
	Embedded  *NestedObject
	Intervals TimeIntervals
}

func (s *TObject) JSONField(fieldName string, storeTemp Map) (ptr interface{}, err error) {
	switch fieldName {
	case "id":
		ptr = &s.ID
	case "name":
		storeTemp["idd"] = &s.Name
		ptr = &s.Name
	case "embedded":
		ptr = &s.Embedded
	case "intervals":
		ptr = &s.Intervals
	}
	return
}

func (s *TObject) JSONFinish(storeTemp Map) error {
	log.Println("FINISH...", storeTemp.String("idd", "!"))
	return nil
}

func TestDecoder(t *testing.T) {
	/*
		//src := []byte("[ [ 1, 2 ], [ 3, 4 ] ]")
		src := []byte("null")
		var sl *TimeIntervals
		err := DecodeBytes(src, &sl)
		log.Println(sl, err)
	*/
	// object
	fObj, err := os.Open("test_object.json")
	if err != nil {
		t.Error(err)
	}

	var obj TObject
	if err := Decode(fObj, &obj); err != nil {
		t.Error(err)
	}
	fObj.Close()
	log.Println("OBJ", obj, err, obj.Embedded, obj.Intervals)
	//log.Println("OBJ", obj, obj.embedded, obj.intervals)
	// slice
	fObj, err = os.Open("test_array.json")
	if err != nil {
		t.Error(err)
	}

	var arr []*TObject
	if err := Decode(fObj, &arr); err != nil {
		t.Error(err)
	}
	fObj.Close()
	log.Println("ARR", arr)
	for _, v := range arr {
		log.Println(v)
	}
}

/////////////////////////////////////////////////////////////////// project

type LeadVerifySettings struct {
	BotEnabled     bool
	ManagerEnabled bool
	CallDecryption bool
	//CallLeadIntervals          []*TimeInterval
	CallsMaxCount int
	FirstDelay    int
	//QualityCriteries           []string
	RecallMinPeriod            int
	ProductDescriptionFilename string
	ScriptFilename             string
}

type Project struct {
	ID                 int64
	Name               string
	LeadVerifySettings *LeadVerifySettings
}

func (s *Project) JSONField(fieldName string, storeTemp Map) (ptr interface{}, err error) {
	switch fieldName {
	case "id":
		ptr = &s.ID
	case "name":
		ptr = &s.Name
	case "lead_verify_settings":
		ptr = &s.LeadVerifySettings
	}
	return
}

func TestProject(t *testing.T) {
	src := `{"id": 1536, "name": "Первый проект 2", "lead_verify_settings": {"bot_enabled": false, "first_delay": 81000, "call_decryption": true, "calls_max_count": 100, "manager_enabled": false, "script_filename": null, "quality_criteries": ["Ок (качественный лид)", "номер не существует", "номер не отвечает/заблокирован/сбрасывает/молчание в трубке ХХХ дней", "номер заблокирован", "не автор заявки", "не наше ГЕО", "не заинтересован в продукте (у клиента другой вопрос)", "нет 18 лет", "несоответствие языка", "неразборчивая речь"], "recall_min_period": 9600, "call_lead_day_intervals": [[3600, 7200]], "product_description_filename": null}}`
	log.Println(string(src))
	var proj Project
	err := DecodeBytes([]byte(src), &proj)
	log.Println(proj, err)
}

func TestMap(t *testing.T) {
	m := Map{
		"one": float64(10.2),
	}
	log.Println(m.Int32("one", 0))
}

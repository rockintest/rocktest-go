package scenario

import (
	"fmt"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

type Seq struct {
	value int
	step  int
}

func getSeqs(scenario *Scenario) (map[string]*Seq, error) {
	// Put in the scenario store "seqs" a map of *Seq
	iseqs := scenario.GetStore("seqs")
	var seqs map[string]*Seq
	var ok bool

	if iseqs == nil {
		seqs = make(map[string]*Seq)
		scenario.PutStore("seqs", seqs)
	} else {
		seqs, ok = iseqs.(map[string]*Seq)
		if !ok {
			return nil, fmt.Errorf("%s", "Cast error")
		}
	}

	return seqs, nil
}

func (module *Module) Id_initseq(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	seqs, err := getSeqs(scenario)

	if err != nil {
		return err
	}

	name, _ := scenario.GetString(paramsEx, "name", "default")
	val, _ := scenario.GetNumber(paramsEx, "value", 0)
	step, _ := scenario.GetNumber(paramsEx, "step", 1)

	log.Debugf("Init sequence %s. Val=%s Step=%d", name, val, step)

	seq := new(Seq)
	seq.step = step
	seq.value = val - step

	seqs[name] = seq

	return nil
}

func (module *Module) Id_seqMeta() Meta {
	return Meta{Ext: "unused", Params: []string{"name"}}
}

func (module *Module) Id_seq(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	seqs, err := getSeqs(scenario)

	if err != nil {
		return err
	}

	name, _ := scenario.GetString(paramsEx, "name", "default")
	if name == "" {
		name = "default"
	}

	log.Tracef("Get next value for sequence %s", name)

	seq, ok := seqs[name]

	var ret int

	if !ok {
		seq := new(Seq)
		seq.step = 1
		seq.value = 0
		seqs[name] = seq
		ret = 0
	} else {
		ret = seq.value + seq.step
		seq.value = ret
	}

	scenario.PutContextAs(paramsEx, "seq", "result", ret)
	scenario.PutContext("??", ret)

	return nil

}

func (module *Module) Id_uuid(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	ret := fmt.Sprintf("%v", uuid.New())

	scenario.PutContextAs(paramsEx, "uuid", "result", ret)
	scenario.PutContext("??", ret)

	return nil
}

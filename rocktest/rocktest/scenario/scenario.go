package scenario

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	"io.rocktest/rocktest/text"

	log "github.com/sirupsen/logrus"
)

type Scenario struct {
	Steps []interface{}
	M     Module

	// Contains the variables
	Context map[string]string

	Subst *text.StringSubstitutor

	Quoter      text.ParamQuoter
	SubstQuoter *text.StringSubstitutor

	Executor      *InlineExecutor
	SubstExecutor *text.StringSubstitutor

	// Caller, when the scenario is called from another one
	Caller *Scenario

	// Root, from where looking from modules
	Root string

	Skip bool
}

func NewScenario() *Scenario {
	ret := new(Scenario)
	ret.Context = make(map[string]string)
	ret.Subst = text.NewStringSubstitutorByMap(ret.Context)
	ret.SubstQuoter = text.NewStringSubstitutorByLookuper(ret.Quoter)
	ret.Executor = NewInlineExecutor(ret)
	ret.SubstExecutor = text.NewStringSubstitutorByLookuper(ret.Executor)
	ret.Skip = false
	return ret
}

func nodeToList(node interface{}) []interface{} {
	m, ok := node.([]interface{})
	if !ok {
		panic(fmt.Sprintf("%v is not of type list", node))
	}
	return m
}

func nodeToMap(node interface{}) map[string]interface{} {
	m, ok := node.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("%v is not of type map", node))
	}
	return m
}

func humanStr(src interface{}) string {
	if src != nil {
		return fmt.Sprint(src)
	} else {
		return ""
	}
}

func (s *Scenario) isBuiltin(name string) bool {
	return name == "module" || name == "step"
}

func (s *Scenario) RunSteps(steps []interface{}) error {
	var i int = 0
	var stop bool = false

	for _, v := range steps {
		i++

		log.Infof("---------- Step %s/%d ----------", s.Context["module"], i)

		s.Context["step"] = fmt.Sprint(i)

		node := nodeToMap(v)
		step := NewStep(node, s)

		if !s.Skip && !strings.HasPrefix(step.Type, "--") {
			log.Infof("[%s/%d] %s - %s", s.Context["module"], i, strings.ToUpper(step.Type), humanStr(step.Params["value"]))
			if log.IsLevelEnabled(log.DebugLevel) {
				for name, val := range step.Params {
					log.Debugf("  %s = %s", name, val)
				}
			}
		}

		switch step.Type {
		case "Exit":
			stop = true
		default:

			if step.Type == "Resume" || (!s.Skip && !strings.HasPrefix(step.Type, "--")) {
				err := step.Exec()

				if err != nil {

					yamlMap := make(map[string]interface{})
					yamlMap["stepNumber"] = i
					yamlMap["scenario"] = s.Context["module"]
					if step.Desc != "" {
						yamlMap["desc"] = step.Desc
					}
					yamlMap["error"] = err.Error()
					yamlMap["step"] = node

					yamlString, _ := yaml.Marshal(yamlMap)

					return errors.New(string(yamlString))
				}

			}
		}

		if stop {
			break
		}
	}

	return nil
}

func (s *Scenario) RunFromRoot(scen string) error {
	log.Infof("Run scenario %s", scen)

	yamlFile, err := ioutil.ReadFile(s.Root + "/" + scen)
	if err != nil {
		return err
	}

	basename := filepath.Base(scen)

	s.Context["module"] = strings.TrimSuffix(basename, filepath.Ext(basename))

	err = yaml.Unmarshal(yamlFile, &s.Steps)

	if err != nil {
		return err
	}

	steps := nodeToList(s.Steps)

	return s.RunSteps(steps)

}

func (s *Scenario) Run(scen string) error {
	log.Infof("Run scenario %s", scen)

	yamlFile, err := ioutil.ReadFile(scen)
	if err != nil {
		return err
	}

	basename := filepath.Base(scen)

	abs, _ := filepath.Abs(scen)

	s.Root = filepath.Dir(abs)
	s.Context["module"] = strings.TrimSuffix(basename, filepath.Ext(basename))

	err = yaml.Unmarshal(yamlFile, &s.Steps)

	if err != nil {
		return err
	}

	steps := nodeToList(s.Steps)

	return s.RunSteps(steps)

}

func (s *Scenario) PutContext(name string, value interface{}) error {

	switch str := value.(type) {
	case string:
		s.Context[name] = str
		log.Debugf("Set %s: %s = %v", s.Context["module"], name, value)
	case int:
		s.Context[name] = fmt.Sprint(str)
		log.Debugf("Set %s: %s = %v", s.Context["module"], name, value)
	case bool:
		s.Context[name] = fmt.Sprint(str)
		log.Debugf("Set %s: %s = %v", s.Context["module"], name, value)
	default:
		log.Debugf("NotSet %s: %s = %v (type must be string or int, not %T)", s.Context["module"], name, value, value)
		return fmt.Errorf("variable value type must be string or int, not %T", value)
	}

	return nil
}

func (s *Scenario) AddVariables(params map[string]interface{}) error {

	for k, v := range params {

		s.PutContext(k, v)

	}

	return nil
}

func (s *Scenario) CopyVariable(name string, source *Scenario) error {

	val, found := source.Context[name]
	if found {
		err := s.PutContext(name, val)
		if err != nil {
			return err
		}
	}

	return nil
}

// Add all the variables from params, excluding builtin variables
func (s *Scenario) CopyVariables(source *Scenario) error {

	for k, v := range source.Context {

		if s.isBuiltin(k) {
			continue
		}

		s.PutContext(k, v)
	}

	return nil
}

// Replace the values of the variables in a list
func (s *Scenario) ExpandList(params []interface{}) []interface{} {

	ret := make([]interface{}, len(params))

	for i, v := range params {
		ret[i] = s.Expand(v)
	}

	return ret
}

// Replace the values
func (s *Scenario) Expand(params interface{}) interface{} {

	switch paramCast := params.(type) {
	case string:
		return s.ExpandString(paramCast)
	case []interface{}:
		return s.ExpandList(paramCast)
	case map[string]interface{}:
		return s.ExpandMap(paramCast)
	default:
		return params
	}

}

// Replace the value of the variable
func (s *Scenario) ExpandString(param string) string {

	// First, quote the parameters.
	// ${$module(p1,p2)} => ${$module(<<[p1]>>,<<[p2]>>)}
	// This way, if p1 or p2 contain commas, it will work
	ret := s.SubstQuoter.Replace(param)

	// Then substitute the variables
	// ${rock} => test
	ret2 := s.Subst.Replace(ret)

	// Finaly, call the modules inline
	// ${$tolower(ROCK)} => rock
	// The function Tolowere is called and the result is returned
	ret3 := s.SubstExecutor.Replace(ret2)

	log.Tracef("%s => %s => %s => %s", param, ret, ret2, ret3)

	return ret3
}

// Replace the values of the variables in a map
func (s *Scenario) ExpandMap(params map[string]interface{}) map[string]interface{} {

	ret := make(map[string]interface{})

	for k, v := range params {
		ret[k] = s.Expand(v)
	}

	return ret

}

// Get a parameter as string. Returns the default value if not found.
// If the value is not found, and there is no default value, returns an error.
// If the value is not a string, return an error
func (s *Scenario) GetString(params map[string]interface{}, key string, def interface{}) (string, error) {

	if params == nil {
		if def == nil {
			return "", errors.New("Params map empty, and no default value provided for key " + key)
		} else {
			return def.(string), nil
		}
	}

	ret, ok := params[key]

	if ok {

		if ret == "<nil>" {

			return "", nil

		} else {
			switch ret := ret.(type) {
			case string:
				return ret, nil
			case int:
				return fmt.Sprint(ret), nil
			case bool:
				return fmt.Sprint(ret), nil
			default:
				msg := fmt.Sprintf("Bad type for value %s. Must be string, not %v", key, reflect.TypeOf(ret))
				return "", errors.New(msg)

			}

		}

	} else {
		if def == nil {
			return "", errors.New("Value not found for " + key)
		} else {
			return def.(string), nil
		}
	}

}

// Get a parameter as integer. Returns the default value if not found.
// If the value is not found, and there is no default value, returns an error.
// If the value is not a string, return an error
func (s *Scenario) GetNumber(params map[string]interface{}, key string, def interface{}) (int, error) {

	if params == nil {
		if def == nil {
			return 0, errors.New("Params map empty, and no default value provided for key " + key)
		} else {
			return def.(int), nil
		}
	}

	ret, ok := params[key]

	if ok {
		switch ret := ret.(type) {
		case string:
			iret, err := strconv.Atoi(ret)
			if err != nil {
				return 0, err
			}
			return iret, nil
		case int:
			return ret, nil
		default:
			msg := fmt.Sprintf("Bad type for value %s. Must be int, not %v", key, reflect.TypeOf(ret))
			return 0, errors.New(msg)
		}

	} else {
		if def == nil {
			return 0, errors.New("Value not found for " + key)
		} else {
			return def.(int), nil
		}
	}

}

// Get a parameter as boolean. Returns the default value if not found.
// If the value is not found, and there is no default value, returns an error.
// If the value is not a string, return an error
func (s *Scenario) GetBool(params map[string]interface{}, key string, def interface{}) (bool, error) {

	if params == nil {
		if def == nil {
			return false, errors.New("Params map empty, and no default value provided for key " + key)
		} else {
			return def.(bool), nil
		}
	}

	ret, ok := params[key]

	if ok {
		switch ret := ret.(type) {
		case string:
			bret, err := strconv.ParseBool(ret)
			if err != nil {
				return false, err
			}
			return bret, nil
		case bool:
			return ret, nil
		default:
			msg := fmt.Sprintf("Bad type for value %s. Must be bool, not %v", key, reflect.TypeOf(ret))
			return false, errors.New(msg)
		}

	} else {
		if def == nil {
			return false, errors.New("Value not found for " + key)
		} else {
			return def.(bool), nil
		}
	}

}

func asList(def interface{}) ([]interface{}, error) {
	switch defcast := def.(type) {
	case []interface{}:
		return defcast, nil
	case []string:
		ret := make([]interface{}, len(defcast))
		for i, v := range defcast {
			ret[i] = v
		}
		return ret, nil
	case string:
		return []interface{}{defcast}, nil
	case int:
		return []interface{}{fmt.Sprint(defcast)}, nil
	default:
		return nil, fmt.Errorf("bad type for default value. Must be a string, int, []interface{} or []string, not %T", defcast)
	}
}

// Get a parameter as list. Returns the default value if not found.
// If the value is not found, and there is no default value, returns an error.
// If the value is not a list, return an error
func (s *Scenario) GetList(params map[string]interface{}, key string, def interface{}) ([]interface{}, error) {

	if params == nil {
		if def == nil {
			return nil, errors.New("Params map empty, and no default value provided for key " + key)
		} else {
			return asList(def)
		}
	}

	ret, ok := params[key]

	if ok {
		switch ret := ret.(type) {
		case []interface{}:
			return ret, nil
		case string:
			return []interface{}{ret}, nil
		case int:
			return []interface{}{fmt.Sprint(ret)}, nil
		default:
			msg := fmt.Sprintf("Bad type for value %s. Must be a list, not %v", key, reflect.TypeOf(ret))
			return nil, errors.New(msg)
		}

	} else {
		if def == nil {
			return nil, errors.New("Value not found for " + key)
		} else {
			return asList(def)
		}
	}

}

// Execute a node. Locates the function in the module, and call it.
func (s *Scenario) Exec(val string, params map[string]interface{}) error {

	val2 := strings.ReplaceAll(val, ".", "_")
	val2 = strings.Title(val2)

	var paramsExec = []reflect.Value{
		reflect.ValueOf(params),
		reflect.ValueOf(s),
	}

	meth := reflect.ValueOf(&s.M).MethodByName(val2)

	if !meth.IsValid() {
		return errors.New("Unknown step type: " + val)
	}

	ret := reflect.ValueOf(&s.M).MethodByName(val2).Call(paramsExec)

	if !ret[0].IsNil() {
		return ret[0].Interface().(error)
	} else {
		return nil
	}

}

// Gets the meta-informations on a module
func (s *Scenario) Meta(val string) Meta {

	var ret Meta

	val2 := strings.ReplaceAll(val, ".", "_")
	val2 = strings.Title(val2)

	meth := reflect.ValueOf(&s.M).MethodByName(val2 + "Meta")

	if !meth.IsValid() {
		return ret
	}

	metaret := reflect.ValueOf(&s.M).MethodByName(val2 + "Meta").Call([]reflect.Value{})
	return metaret[0].Interface().(Meta)

}

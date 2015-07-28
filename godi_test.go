package godi

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type I1 interface {
	F1() string
}

type T1 struct {
	s string
}

func (p T1) F1() string {
	return p.s
}

type T2 struct {
	f2 bool
}

func (p T2) F1() string {
	return "t2"
}

func (p T2) F2() string {
	return "t2f2"
}

func TestRegisterType(t *testing.T) {
	Reset()

	RegisterType((*I1)(nil))
	count := len(*getRegisteredTypes())

	if count != 1 {
		t.Error(fmt.Sprintf("Expected 1 types, got %d", count))
	}
}

func TestRegisterDupe(t *testing.T) {
	Reset()
	RegisterType(T1{})

	err := RegisterType(T1{})

	if err == nil {
		t.Error("Expected Error")
	}
}

func TestResolveInstance(t *testing.T) {
	Reset()

	i1 := (*I1)(nil)
	t1 := &T1{s: "foobarx"}

	res, err := RegisterInstanceImplementor(i1, t1)

	if err != nil {
		t.Error("Expected reg")
	}

	t1_r, err2 := Resolve(i1)
	if err2 != nil {
		t.Error("Error resolving: " + err2.Error())
	}

	t1_val := t1_r.(I1)
	str := t1_val.F1()

	if str != t1.s {
		t.Error("Got " + str)
	}

	res.Close()

	_, err3 := Resolve(i1)

	if err3 == nil {
		t.Error("Expected unregistration error")
	}

}

type I2 interface {
	Bar()
}

func TestResolveInstanceFail(t *testing.T) {
	Reset()

	i1 := (*I1)(nil)
	t1 := &T1{s: "foobarx"}

	res, _ := RegisterInstanceImplementor(i1, t1)

	t2, err2 := Resolve((*I2)(nil))

	if t2 != nil {
		t.Error("unexpected value for t2")
	}

	if err2.Error() != ErrorRegistrationNotFound {
		t.Error(fmt.Sprintf("Error resolving: %v", err2))
	}

	res.Close()

}

func TestResolveOverride(t *testing.T) {
	Reset()

	i1 := (*I1)(nil)
	t1 := &T1{s: "foobar1"}
	t2 := &T1{s: "foobar2"}

	res, err := RegisterInstanceImplementor(i1, t1)

	if err != nil {
		t.Error("Expected reg")
	}

	t1_r, err2 := Resolve(i1)
	if err2 != nil {
		t.Error("Error resolving: " + err2.Error())
	}

	t1_val := t1_r.(I1)
	str := t1_val.F1()

	if str != t1.s {
		t.Error("Got " + str)
	}

	res2, err2 := RegisterInstanceImplementor(i1, t2)
	t1_r, _ = Resolve(i1)
	if str != t1.s {
		t.Error("Got " + str)
	}

	res2.Close()

	t1_r, _ = Resolve(i1)
	if str != t1.s {
		t.Error("Got " + str)
	}

	res.Close()

	_, err3 := Resolve(i1)

	if err3 == nil {
		t.Error("Expected unregistration error")
	}

}
func TestResolveType(t *testing.T) {
	Reset()

	i1 := (*I1)(nil)
	t1 := T1{}

	res, err := RegisterTypeImplementor(i1, t1, false, nil)

	if err != nil {
		t.Error("Expected reg")
	}

	t1_r, err2 := Resolve(i1)
	if err2 != nil || t1_r == nil {
		t.Error("Error resolving: " + err2.Error())
	}

	t1_val := t1_r.(I1)

	str := t1_val.F1()

	if str != t1.s {
		t.Error("Got " + str)
	}

	res.Close()

	_, err3 := Resolve(i1)

	if err3 == nil {
		t.Error("Expected unregistration error")
	}

}

type TestInitializer struct {
}

func (p TestInitializer) CanInitialize(instance interface{}, typeName string) bool {
	if typeName == "godi.T1" {
		return true
	}
	return false
}

var initS = "hodor"

func (p TestInitializer) Initialize(instance interface{}, typeName string) (interface{}, error) {

	if typeName == "godi.T1" {
		t1 := instance.(*T1)
		t1.s = initS
		return t1, nil
	}
	return instance, nil
}

func TestInstanceInitializer(t *testing.T) {
	Reset()

	init := TestInitializer{}

	RegisterInstanceInitializer(init)

	i1 := (*I1)(nil)
	RegisterTypeImplementor(i1, T1{}, false, nil)
	t1_r, _ := Resolve(i1)
	t1_c := t1_r.(*T1)

	if t1_c.s != initS {
		t.Error("Expected " + initS)
	}

}

func TestResolvePendingFail(t *testing.T) {

	defer func() {
		if e := recover(); e != nil {
			if strings.Contains(e.(string), "I1") {
				t.Error("Didn't expect I1")
			} else if !strings.Contains(e.(string), "di.T2") {
				t.Error("Expected T2")
			}
		}
	}()

	Reset()

	RegisterByName("godi.I1", "godi.T2", false)

	i1 := (*I1)(nil)

	RegisterType(i1)
	Resolve(i1)
	t.Error("Expected panic.")
}

func TestResolvePending(t *testing.T) {
	Reset()
	RegisterByName("godi.I1", "godi.T2", false)

	i1 := (*I1)(nil)
	if e1 := RegisterType(i1); e1 != nil {
		t.Error(e1)
	}
	if e2 := RegisterType(T2{}); e2 != nil {
		t.Error(e2)
	}

	r1, _ := Resolve(i1)
	r2 := r1.(I1).F1()
	if r2 != "t2" {
		t.Error(fmt.Sprintf("pending resolve fail %v", r2))
	}
}

func TestCreateScope(t *testing.T) {
	Reset()
	i1 := (*I1)(nil)
	t1 := T1{}

	if e1 := RegisterType(i1); e1 != nil {
		t.Error(e1)
	}

	if e2 := RegisterType(T1{}); e2 != nil {
		t.Error(e2)
	}

	RegisterInstanceImplementor(i1, t1)

	r1, _ := Resolve(i1)
	r2 := r1.(I1).F1()
	if r2 != "" {
		t.Error(fmt.Sprintf("pending resolve fail %v", r2))
	}

	// push a scope
	s2 := CreateScope(true)

	t2 := T2{}
	RegisterInstanceImplementor(i1, t2)

	r3, _ := Resolve(i1)
	r2 = r3.(I1).F1()
	if r2 != "t2" {
		t.Error(fmt.Sprintf("pending resolve fail %v", r2))
	}

	s2.Close()

}

func TestFormatType(t *testing.T) {
	typeName := "*list.List"

	typeName = formatType(typeName)

	if typeName != "list.List" {
		t.Error(typeName)
	}
}

type T3 struct {
	n int
}

func (p *T3) Initialize() bool {
	p.n = 42
	return false
}

var _ Initializable = &T3{}

func (p T3) F1() string {
	return strconv.Itoa(p.n)
}

func TestInitializerInterface(t *testing.T) {
	Reset()
	i1 := (*I1)(nil)
	RegisterTypeImplementor(i1, T3{}, true, nil)

	r3, _ := Resolve(i1)
	r2 := r3.(I1).F1()

	if r2 != "42" {
		t.Errorf("Expected 42, got %s", r2)
	}
}

func TestInitializeCallback(t *testing.T) {
	Reset()
	i1 := (*I1)(nil)

	init := func(inst interface{}) bool {
		t3 := inst.(*T3)
		t3.n = 100
		return false
	}

	RegisterTypeImplementor(i1, T3{}, true, init)

	r3, _ := Resolve(i1)
	r2 := r3.(I1).F1()

	if r2 != "100" {
		t.Errorf("Expected 100, got %s", r2)
	}
}

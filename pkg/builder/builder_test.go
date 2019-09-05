package builder

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/staffano/sonbyg/pkg/utils"
)

func doTime(v utils.Variables) *Artifact {
	log.Printf("Doing time %v", v["ID"])
	d := v["SLEEP"].(time.Duration)
	time.Sleep(d)
	return nil
}

// Download downloads all files specified in SOURCES
func Timely(b *Builder, id string, v *utils.Variables, delay time.Duration) *Task {
	t := NewTask(b, id)
	t.Variables = v.Copy()
	t.Variables["SLEEP"] = delay
	t.Variables["ID"] = id
	t.AssignDefaultSignature()
	return b.Add(t, doTime)
}

func TestNoDependencies(t *testing.T) {
	b := NewBuilder(context.Background(), "test", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	t1 := Timely(b, "t1", &b.Variables, 2*time.Second)
	t2 := Timely(b, "t2", &b.Variables, 2*time.Second)
	t3 := Timely(b, "t3", &b.Variables, 2*time.Second)
	t4 := Timely(b, "t4", &b.Variables, 2*time.Second)
	b.Build(t1, t2, t3, t4)

}

func TestDependencies(t *testing.T) {
	b := NewBuilder(context.Background(), "test", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	t1 := Timely(b, "t1", &b.Variables, 2*time.Second)
	t2 := Timely(b, "t2", &b.Variables, 2*time.Second)
	t3 := Timely(b, "t3", &b.Variables, 2*time.Second)
	t4 := Timely(b, "t4", &b.Variables, 2*time.Second)

	t3.DependsOn(t4)
	t2.DependsOn(t3)
	t1.DependsOn(t2)
	b.Build(t2, t4)

}

func TestCircularDependencies(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	b := NewBuilder(context.Background(), "test", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	t1 := Timely(b, "t1", &b.Variables, 2*time.Second)
	t2 := Timely(b, "t2", &b.Variables, 2*time.Second)
	t3 := Timely(b, "t3", &b.Variables, 2*time.Second)
	t4 := Timely(b, "t4", &b.Variables, 2*time.Second)

	t3.DependsOn(t4)
	t2.DependsOn(t3)
	t1.DependsOn(t2)
	t2.DependsOn(t1)
	b.Build(t1, t4)
}

func Test1000Targets(t *testing.T) {
	b := NewBuilder(context.Background(), "test", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	ta := make([]*Task, 1000)
	for i := range ta {
		ta[i] = Timely(b, fmt.Sprintf("t%d", i), &b.Variables, 1*time.Millisecond)
	}
	for i := 0; i < len(ta)-2; i++ {
		ta[i].DependsOn(ta[i+1])
	}
	b.Build(ta[0])
}

func Test10000Targets(t *testing.T) {
	b := NewBuilder(context.Background(), "test", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	ta := make([]*Task, 10000)
	for i := range ta {
		ta[i] = Timely(b, fmt.Sprintf("t%d", i), &b.Variables, 1*time.Millisecond)
	}
	for i := 0; i < len(ta)-2; i++ {
		ta[i].DependsOn(ta[i+1])
	}
	b.Build(ta[0])
}
func Test100000Targets(t *testing.T) {
	b := NewBuilder(context.Background(), "test", 1, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	ta := make([]*Task, 100000)
	for i := range ta {
		ta[i] = Timely(b, fmt.Sprintf("t%d", i), &b.Variables, 1*time.Millisecond)
	}
	for i := 0; i < len(ta)-2; i++ {
		ta[i].DependsOn(ta[i+1])
	}
	b.Build(ta[0])
}
func Test100Targets5Workers(t *testing.T) {
	b := NewBuilder(context.Background(), "test", 5, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	ta := make([]*Task, 100)
	for i := range ta {
		ta[i] = Timely(b, fmt.Sprintf("t%d", i), &b.Variables, 1000*time.Millisecond)
	}
	for i := 0; i < len(ta)-2; i++ {
		ta[i].DependsOn(ta[i+1])
	}
	b.Build(ta...)
}

func Test100Targets5WorkersNoDeps(t *testing.T) {
	b := NewBuilder(context.Background(), "test", 5, true, "x86_64-w64-mingw32",
		"x86_64-w64-mingw32", "x86_64-w64-mingw32", "", "", "", "")
	ta := make([]*Task, 100)
	for i := range ta {
		ta[i] = Timely(b, fmt.Sprintf("t%d", i), &b.Variables, 1000*time.Millisecond)
	}

	b.Build(ta...)
}

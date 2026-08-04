package build

import (
	"reflect"
	"testing"
)

func TestGoTestJSONArgsPlacesJSONBeforePackagesAndProgramArgs(t *testing.T) {
	got := goTestJSONArgs(
		[]string{"test", "-mod=mod", "-count=1"},
		[]string{"./__workspace/suite"},
		[]string{"--config_dir=/project/config", "--config_file=config-local.ini"},
	)
	want := []string{
		"test", "-mod=mod", "-count=1", "-json", "./__workspace/suite",
		"-args", "--config_dir=/project/config", "--config_file=config-local.ini",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("goTestJSONArgs() = %#v, want %#v", got, want)
	}
}

func TestGoTestJSONArgsKeepsJSONBeforeMultiplePackages(t *testing.T) {
	got := goTestJSONArgs(
		[]string{"test", "-mod=readonly"},
		[]string{"./first", "./second"},
		nil,
	)
	want := []string{"test", "-mod=readonly", "-json", "./first", "./second"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("goTestJSONArgs() = %#v, want %#v", got, want)
	}
}

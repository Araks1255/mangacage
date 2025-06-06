package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/volumes"
)

func TestCreateVolume(t *testing.T) {
	scenarios := volumes.GetCreateVolumeScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestEditVolume(t *testing.T) {
	scenarios := volumes.GetEditVolumeScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTitleVolumes(t *testing.T) {
	scenarios := volumes.GetGetTitleVolumesScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetVolume(t *testing.T) {
	scenarios := volumes.GetGetVolumeScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

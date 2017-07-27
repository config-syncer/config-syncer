package version

import "fmt"

func (m *Version) Print() {
	if m.Name != "" {
		fmt.Printf("Name = %v\n", m.Name)
	}
	fmt.Printf("Version = %v\n", m.Version)
	fmt.Printf("VersionStrategy = %v\n", m.VersionStrategy)
	fmt.Printf("Os = %v\n", m.Os)
	fmt.Printf("Arch = %v\n", m.Arch)

	fmt.Printf("CommitHash = %v\n", m.CommitHash)
	fmt.Printf("GitBranch = %v\n", m.GitBranch)
	fmt.Printf("GitTag = %v\n", m.GitTag)
	fmt.Printf("CommitTimestamp = %v\n", m.CommitTimestamp)

	if m.BuildTimestamp != "" {
		fmt.Printf("BuildTimestamp = %v\n", m.BuildTimestamp)
	}
	if m.BuildHost != "" {
		fmt.Printf("BuildHost = %v\n", m.BuildHost)
	}
	if m.BuildHostOs != "" {
		fmt.Printf("BuildHostOs = %v\n", m.BuildHostOs)
	}
	if m.BuildHostArch != "" {
		fmt.Printf("BuildHostArch = %v\n", m.BuildHostArch)
	}
}

package migrate

import (
	"io/ioutil"
	"strings"
)

type Migration struct {
	Up   []string
	Down []string
}

func (m *Migration) ParseSource(fileDirectory string) error {
	file, err := ioutil.ReadFile(fileDirectory)
	if err != nil {
		return err
	}

	queries := strings.Split(string(file), "-- DOWN")

	upQueries := strings.TrimLeft(queries[0], "-- UP")
	m.Up = strings.SplitAfter(upQueries, ";")
	m.Up = m.Up[:len(m.Up) -1]

	for i := 0; i < len(m.Up); i++ {
		m.Up[i] = strings.TrimSpace(m.Up[i])
	}

	if len(queries) > 1 {
		downQueries := strings.TrimSpace(queries[1])
		m.Down = strings.SplitAfter(downQueries, ";")
		m.Down = m.Down[:len(m.Down) -1]

		for i := 0; i < len(m.Down); i++ {
			m.Down[i] = strings.TrimSpace(m.Down[i])
		}
	}

	return nil
}

package actions

import (
	"strconv"
	"testing"
	"time"

	"github.com/gaia-adm/pumba/container"
	"github.com/gaia-adm/pumba/container/mockclient"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeContainersN(n int) ([]string, []container.Container) {
	names := make([]string, n)
	cs := make([]container.Container, n)
	for i := range cs {
		cs[i] = *container.NewContainer(
			&dockerclient.ContainerInfo{
				Name: "c" + strconv.Itoa(i),
			},
			nil,
		)
		names[i] = "c" + strconv.Itoa(i)
	}
	return names, cs
}

func TestPattern_DotRe2Filter(t *testing.T) {
	c1 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "c1",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	c2 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "c2",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	c3 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "c3",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	cf := regexContainerFilter(".")
	assert.True(t, cf(c1))
	assert.True(t, cf(c2))
	assert.True(t, cf(c3))
}

func TestPattern_Re2Filter(t *testing.T) {
	c1 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "AbcEFG",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	c2 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "AbcHKL",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	cc := &dockerclient.ContainerConfig{
		Labels: map[string]string{"com.gaiaadm.pumba": "true"},
	}
	c3 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "AbcPumba",
			Config: cc,
		},
		nil,
	)
	cf := regexContainerFilter("^Abc")
	assert.True(t, cf(c1))
	assert.True(t, cf(c2))
	assert.False(t, cf(c3))
}

func TestNamesFilter(t *testing.T) {
	c1 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "ccc",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	c2 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "ddd",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	cc := &dockerclient.ContainerConfig{
		Labels: map[string]string{"com.gaiaadm.pumba": "true"},
	}
	c3 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "xxx",
			Config: cc,
		},
		nil,
	)
	cf := containerFilter([]string{"ccc", "bbb", "xxx"})
	assert.True(t, cf(c1))
	assert.False(t, cf(c2))
	assert.False(t, cf(c3))
}

func TestAllNamesFilter(t *testing.T) {
	c1 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "ccc",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	c2 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "ddd",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	cc := &dockerclient.ContainerConfig{
		Labels: map[string]string{"com.gaiaadm.pumba": "true"},
	}
	c3 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "xxx",
			Config: cc,
		},
		nil,
	)
	cf := containerFilter([]string{})
	assert.True(t, cf(c1))
	assert.True(t, cf(c2))
	assert.False(t, cf(c3))
}

func TestAllFilter(t *testing.T) {
	c1 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "ccc",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	c2 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "ddd",
			Config: &dockerclient.ContainerConfig{},
		},
		nil,
	)
	cc := &dockerclient.ContainerConfig{
		Labels: map[string]string{"com.gaiaadm.pumba": "true"},
	}
	c3 := *container.NewContainer(
		&dockerclient.ContainerInfo{
			Name:   "xxx",
			Config: cc,
		},
		nil,
	)
	assert.True(t, allContainersFilter(c1))
	assert.True(t, allContainersFilter(c2))
	assert.False(t, allContainersFilter(c3))
}

func TestStopByName(t *testing.T) {
	names, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for _, c := range cs {
		client.On("StopContainer", c, 10).Return(nil)
	}
	err := Pumba{}.StopContainers(client, names, "", 10)
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestStopByNameRandom(t *testing.T) {
	names, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("StopContainer", mock.AnythingOfType("container.Container"), 10).Return(nil)
	RandomMode = true
	err := Pumba{}.StopContainers(client, names, "", 10)
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestStopByPattern(t *testing.T) {
	_, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for _, c := range cs {
		client.On("StopContainer", c, 10).Return(nil)
	}
	err := Pumba{}.StopContainers(client, []string{}, "^c", 10)
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestStopByPatternRandom(t *testing.T) {
	_, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("StopContainer", mock.AnythingOfType("container.Container"), 10).Return(nil)
	RandomMode = true
	err := Pumba{}.StopContainers(client, []string{}, "^c", 10)
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestKillByName(t *testing.T) {
	names, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for _, c := range cs {
		client.On("KillContainer", c, "SIGTEST").Return(nil)
	}
	err := Pumba{}.KillContainers(client, names, "", "SIGTEST")
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestKillByNameRandom(t *testing.T) {
	names, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("KillContainer", mock.AnythingOfType("container.Container"), "SIGTEST").Return(nil)
	RandomMode = true
	err := Pumba{}.KillContainers(client, names, "", "SIGTEST")
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestKillByPattern(t *testing.T) {
	_, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for i := range cs {
		client.On("KillContainer", cs[i], "SIGTEST").Return(nil)
	}
	err := Pumba{}.KillContainers(client, []string{}, "^c", "SIGTEST")
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestKillByPatternRandom(t *testing.T) {
	_, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("KillContainer", mock.AnythingOfType("container.Container"), "SIGTEST").Return(nil)
	RandomMode = true
	err := Pumba{}.KillContainers(client, []string{}, "^c", "SIGTEST")
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestRemoveByName(t *testing.T) {
	names, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for _, c := range cs {
		client.On("RemoveContainer", c, false, "", "").Return(nil)
	}
	err := Pumba{}.RemoveContainers(client, names, "", false, "", "")
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestRemoveByNameRandom(t *testing.T) {
	names, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("RemoveContainer", mock.AnythingOfType("container.Container"), false, "mylink", "myvol").Return(nil)
	RandomMode = true
	err := Pumba{}.RemoveContainers(client, names, "", false, "mylink", "myvol")
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestRemoveByPattern(t *testing.T) {
	_, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for _, c := range cs {
		client.On("RemoveContainer", c, false, "mylink", "myvol").Return(nil)
	}
	err := Pumba{}.RemoveContainers(client, []string{}, "^c", false, "mylink", "myvol")
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestRemoveByPatternRandom(t *testing.T) {
	_, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("RemoveContainer", mock.AnythingOfType("container.Container"), false, "mylink", "myvol").Return(nil)
	RandomMode = true
	err := Pumba{}.RemoveContainers(client, []string{}, "^c", false, "mylink", "myvol")
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestPauseByName(t *testing.T) {
	names, cs := makeContainersN(10)
	d, _ := time.ParseDuration("2ms")
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for _, c := range cs {
		client.On("PauseContainer", c, d).Return(nil)
	}
	err := Pumba{}.PauseContainers(client, names, "", d)
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestPauseByPattern(t *testing.T) {
	_, cs := makeContainersN(10)
	d, _ := time.ParseDuration("2ms")
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for _, c := range cs {
		client.On("PauseContainer", c, d).Return(nil)
	}
	err := Pumba{}.PauseContainers(client, []string{}, "^c", d)
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestPauseByNameRandom(t *testing.T) {
	names, cs := makeContainersN(10)
	d, _ := time.ParseDuration("2ms")
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("PauseContainer", mock.AnythingOfType("container.Container"), d).Return(nil)
	RandomMode = true
	err := Pumba{}.PauseContainers(client, names, "", d)
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestNetemDealyByName(t *testing.T) {
	names, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for _, c := range cs {
		client.On("DisruptContainer", c, "eth1", "delay 120ms 25ms 15%").Return(nil)
	}
	err := Pumba{}.NetemDelayContainers(client, names, "", "eth1", 1*time.Second, 120, 25, 15)
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

/*
func TestNetemByNameRandom(t *testing.T) {
	names, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("DisruptContainer", mock.AnythingOfType("container.Container"), "delay 1000ms").Return(nil)
	RandomMode = true
	err := Pumba{}.NetemContainers(client, names, "", "delay 1000ms")
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestNetemByPattern(t *testing.T) {
	_, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	for i := range cs {
		client.On("DisruptContainer", cs[i], "delay 3000ms:172.19.0.3").Return(nil)
	}
	err := Pumba{}.NetemContainers(client, []string{}, "^c", "delay 3000ms:172.19.0.3")
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestNetemByPatternRandom(t *testing.T) {
	_, cs := makeContainersN(10)
	client := &mockclient.MockClient{}
	client.On("ListContainers", mock.AnythingOfType("container.Filter")).Return(cs, nil)
	client.On("DisruptContainer", mock.AnythingOfType("container.Container"), "172.19.0.3").Return(nil)
	RandomMode = true
	err := Pumba{}.NetemContainers(client, []string{}, "^c", "172.19.0.3")
	RandomMode = false
	assert.NoError(t, err)
	client.AssertExpectations(t)
}
*/

func TestSelectRandomContainer(t *testing.T) {
	_, cs := makeContainersN(10)
	c1 := randomContainer(cs)
	c2 := randomContainer(cs)
	assert.NotNil(t, c1)
	assert.NotNil(t, c2)
	assert.NotEqual(t, c1.Name(), c2.Name())
}

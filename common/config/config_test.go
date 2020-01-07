package config

import (
	"io/ioutil"
	"testing"

	"git.2dfire.net/zerodb/common/utils"
	"github.com/stretchr/testify/assert"
)

func TestConf(t *testing.T) {
	file, err := ioutil.ReadFile("../../proxy/proxy/test-conf/conf_cobar.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, file)

	config := Config{}

	err = utils.LoadYaml(file, &config)
	assert.Nil(t, err)

	assert.Equal(t, "zerodb", config.Basic.User)
	t.Log(config.Basic.ConfigName)
	/*
		exFile, err := os.Create("xxx.yaml")
		assert.Nil(t, err)
		assert.NotNil(t, exFile)
		data, err := utils.UnLoadYaml(&config)
		assert.Nil(t, err)
		assert.NotNil(t, data)
		_, err = exFile.Write(data)
		assert.Nil(t, err)
		exFile.Close()
	*/
}

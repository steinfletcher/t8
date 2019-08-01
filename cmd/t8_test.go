package cmd_test

import (
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/t8/cmd"
	"github.com/steinfletcher/t8/domain"
	"github.com/steinfletcher/t8/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestT8(t *testing.T) {
	t.Skip()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := domain.Config{}
	promptReader := mocks.NewMockPromptReader(ctrl)
	promptReader.EXPECT().
		String(gomock.Eq(conf.Parameters[0])).
		Return("value", nil)
	templateRenderer := mocks.NewMockTemplateRenderer(ctrl)
	templateRenderer.EXPECT().
		Render("", conf).
		Return()
	fetchTemplate := func(source string, target string) error {
		return nil
	}
	args := []string{}

	err := cmd.Run(fetchTemplate, promptReader, templateRenderer, args)

	assert.NoError(t, err)
}

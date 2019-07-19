package domain_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/t8/domain"
	"github.com/steinfletcher/t8/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptForRequiredParameters_PromptsForStringParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := domain.Config{Parameters: []domain.Parameter{{Name: "a", Type: domain.String}}}
	promptReader := mocks.NewMockPromptReader(ctrl)
	promptReader.EXPECT().
		String(gomock.Eq(conf.Parameters[0])).
		Return("value", nil)
	cmd := domain.Command{}

	config, err := domain.PromptForRequiredParameters(promptReader, conf, cmd)

	assert.NoError(t, err)
	assert.Equal(t, domain.Config{
		Parameters: []domain.Parameter{{Name: "a", Type: domain.String, Actual: "value"}},
	}, config)
}

func TestPromptForRequiredParameters_PromptsForStringParamErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := domain.Config{Parameters: []domain.Parameter{{Name: "a", Type: domain.String}}}
	promptReader := mocks.NewMockPromptReader(ctrl)
	promptReader.EXPECT().
		String(gomock.Eq(conf.Parameters[0])).
		Return("", errors.New("some error"))
	cmd := domain.Command{}

	_, err := domain.PromptForRequiredParameters(promptReader, conf, cmd)

	assert.EqualError(t, err, "failed to read user input: some error")
}

func TestPromptForRequiredParameters_PromptsForOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := domain.Config{Parameters: []domain.Parameter{{Name: "a", Type: domain.Option, Default: []string{"x", "y"}}}}
	promptReader := mocks.NewMockPromptReader(ctrl)
	promptReader.EXPECT().
		Options(gomock.Eq(conf.Parameters[0])).
		Return("x", nil)
	cmd := domain.Command{}

	config, err := domain.PromptForRequiredParameters(promptReader, conf, cmd)

	assert.NoError(t, err)
	assert.Equal(t, domain.Config{
		Parameters: []domain.Parameter{{Name: "a", Type: domain.Option, Default: []string{"x", "y"}, Actual: "x"}},
	}, config)
}

func TestPromptForRequiredParameters_PromptsForOptionsErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := domain.Config{Parameters: []domain.Parameter{{Name: "a", Type: domain.Option, Default: []string{"x", "y"}}}}
	promptReader := mocks.NewMockPromptReader(ctrl)
	promptReader.EXPECT().
		Options(gomock.Eq(conf.Parameters[0])).
		Return("", errors.New("some error"))
	cmd := domain.Command{}

	_, err := domain.PromptForRequiredParameters(promptReader, conf, cmd)

	assert.EqualError(t, err, "failed to read user input: some error")
}

func TestPromptForRequiredParameters_DoesNotPromptIfFlagSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := domain.Config{Parameters: []domain.Parameter{{Name: "a", Type: domain.String}}}
	promptReader := mocks.NewMockPromptReader(ctrl)
	cmd := domain.Command{Flags: []domain.Flag{{
		Key:   "a",
		Value: "value",
	}}}

	config, err := domain.PromptForRequiredParameters(promptReader, conf, cmd)

	assert.NoError(t, err)
	assert.Equal(t, domain.Config{
		Parameters: []domain.Parameter{{Name: "a", Type: domain.String, Actual: "value"}},
	}, config)
}

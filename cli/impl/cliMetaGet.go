package impl

import (
	"ffsyncclient/cli"
	"ffsyncclient/consts"
	"ffsyncclient/fferr"
	"ffsyncclient/langext"
	"ffsyncclient/syncclient"
)

type CLIArgumentsMetaGet struct {
}

func NewCLIArgumentsMetaGet() *CLIArgumentsMetaGet {
	return &CLIArgumentsMetaGet{}
}

func (a *CLIArgumentsMetaGet) Mode() cli.Mode {
	return cli.ModeMetaGet
}

func (a *CLIArgumentsMetaGet) PositionArgCount() (*int, *int) {
	return langext.Ptr(0), langext.Ptr(0)
}

func (a *CLIArgumentsMetaGet) AvailableOutputFormats() []cli.OutputFormat {
	return []cli.OutputFormat{cli.OutputFormatText}
}

func (a *CLIArgumentsMetaGet) ShortHelp() [][]string {
	return [][]string{
		{"ffsclient meta", "Get storage metadata"},
	}
}

func (a *CLIArgumentsMetaGet) FullHelp() []string {
	return []string{
		"$> ffsclient meta",
		"",
		"Get storage metadata",
	}
}

func (a *CLIArgumentsMetaGet) Init(positionalArgs []string, optionArgs []cli.ArgumentTuple) error {
	for _, arg := range optionArgs {
		return fferr.DirectOutput.New("Unknown argument: " + arg.Key)
	}

	return nil
}

func (a *CLIArgumentsMetaGet) Execute(ctx *cli.FFSContext) int {
	ctx.PrintVerbose("[Get Meta]")
	ctx.PrintVerbose("")

	// ========================================================================

	cfp, err := ctx.AbsSessionFilePath()
	if err != nil {
		ctx.PrintFatalError(err)
		return consts.ExitcodeError
	}

	if !langext.FileExists(cfp) {
		ctx.PrintFatalMessage("Sessionfile does not exist.")
		ctx.PrintFatalMessage("Use `ffsclient login <email> <password>` first")
		return consts.ExitcodeNoLogin
	}

	// ========================================================================

	client := syncclient.NewFxAClient(ctx.Opt.AuthServerURL)

	ctx.PrintVerbose("Load existing session from " + cfp)
	session, err := syncclient.LoadSession(ctx, cfp)
	if err != nil {
		ctx.PrintFatalError(err)
		return consts.ExitcodeError
	}

	session, err = client.AutoRefreshSession(ctx, session)
	if err != nil {
		ctx.PrintFatalError(err)
		return consts.ExitcodeError
	}

	// ========================================================================

	record, err := client.GetRecord(ctx, session, consts.CollectionMeta, consts.RecordMetaGlobal, false)
	if err != nil {
		ctx.PrintFatalError(err)
		return consts.ExitcodeError
	}

	// ========================================================================

	if langext.Coalesce(ctx.Opt.Format, cli.OutputFormatText) != cli.OutputFormatText {
		ctx.PrintFatalMessage("Unsupported output-format: " + ctx.Opt.Format.String())
		return consts.ExitcodeUnsupportedOutputFormat
	}

	ctx.PrintPrimaryOutput(langext.TryPrettyPrintJson(record.Payload))
	return 0
}

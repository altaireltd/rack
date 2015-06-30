package snapshotcommands

import (
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osSnapshots "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--volume-id <volumeID>"),
	Description: "Creates a volume",
	Action:      actionCreate,
	Flags:       util.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "volume-id",
			Usage: "[required] The volume ID from which to create this snapshot.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] A name for this snapshot.",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] A description for this snapshot.",
		},
	}
}

var keysCreate = []string{"ID", "Name", "Description", "Size", "VolumeType", "SnapshotID", "Attachments", "CreatedAt"}

type paramsCreate struct {
	opts *osSnapshots.CreateOpts
}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"volume-id"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext

	opts := &osSnapshots.CreateOpts{
		VolumeID:    c.String("volume-id"),
		Name:        c.String("name"),
		Description: c.String("description"),
	}

	resource.Params = &paramsCreate{
		opts: opts,
	}

	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	snapshot, err := osSnapshots.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = snapshotSingle(snapshot)
}

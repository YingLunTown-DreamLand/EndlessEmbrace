package BotFunction

import (
	"fmt"
	NBTAssigner "phoenixbuilder/bdump/nbt_assigner"
)

func UpgradeExecuteCommands(command string) (res string) {
	res, failed_blocks, err := NBTAssigner.UpgradeExecuteCommand(command)
	if err != nil {
		res = fmt.Sprintf("Error Occurred: %v", err)
	}
	if len(failed_blocks) > 0 {
		res = fmt.Sprintf("%s\n\nWarning: Failed block states mapping: %#v", res, failed_blocks)
	}
	res = fmt.Sprintf("%s\n\nPowered by GitHub@LNSSPsd/PhoenixBuilder", res)
	return
}

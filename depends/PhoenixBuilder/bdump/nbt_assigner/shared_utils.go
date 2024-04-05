package NBTAssigner

import (
	"fmt"
	"phoenixbuilder/mirror/chunk"
)

// 取得名称为 blockName 且数据值(附加值)为 metaData 的方块的方块状态。
// 特别地，name **不**需要加上命名空间 minecraft
func get_block_states_from_legacy_block(
	blockName string,
	metaData uint16,
) (map[string]interface{}, error) {
	standardRuntimeID, found := chunk.LegacyBlockToRuntimeID(blockName, metaData)
	if !found {
		return nil, fmt.Errorf("get_block_states_from_legacy_block: Failed to get the runtimeID of block %s; metaData = %d", blockName, metaData)
	}
	generalBlock, found := chunk.RuntimeIDToBlock(standardRuntimeID)
	if !found {
		return nil, fmt.Errorf("get_block_states_from_legacy_block: Failed to converse StandardRuntimeID to NEMCRuntimeID; standardRuntimeID = %d, blockName = %s, metaData = %d", standardRuntimeID, blockName, metaData)
	}
	return generalBlock.Properties, nil
}

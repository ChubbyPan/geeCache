package geecache

// 使用一致性hash选择节点 -> 需要远程节点 -> http访问 -> 成功？ 返回 不成功？退回本地节点处理

// depend on key to chose peergetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// from group to get cache
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

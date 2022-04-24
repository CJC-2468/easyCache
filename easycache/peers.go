//peers.go defines how processes find and communicate with their peers.
package easycache

import (
	pb "easyCache/easycache/easycachepb"
)

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	//根据传入的 key 选择相应peer，如果是远程节点，返回节点和true，如果是本机节点，返回nil和false。
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error //从对应 group 查找缓存值,PeerGetter 就对应于上述流程中的 HTTP 客户端。
}

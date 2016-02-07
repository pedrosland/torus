package ring

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/coreos/agro"
	"github.com/coreos/agro/models"

	"github.com/serialx/hashring"
)

type ketama struct {
	version int
	rep     int
	peers   []string
	ring    *hashring.HashRing
}

func init() {
	registerRing(Ketama, "ketama", makeKetama)
}

func makeKetama(r *models.Ring) (agro.Ring, error) {
	rep := int(r.ReplicationFactor)
	if rep == 0 {
		rep = 1
	}
	return &ketama{
		version: int(r.Version),
		peers:   r.UUIDs,
		rep:     rep,
		ring:    hashring.New(r.UUIDs),
	}, nil
}

func (k *ketama) GetPeers(key agro.BlockRef) (agro.PeerPermutation, error) {
	s, ok := k.ring.GetNodes(string(key.ToBytes()), len(k.peers))
	if !ok {
		return agro.PeerPermutation{}, errors.New("couldn't get sufficient nodes")
	}
	return agro.PeerPermutation{
		Peers:       s,
		Replication: k.rep,
	}, nil
}

func (k *ketama) Members() agro.PeerList { return append([]string(nil), k.peers...) }

func (k *ketama) Describe() string {
	s := fmt.Sprintf("Ring: Ketama\nReplication:%d\nPeers:", k.rep)
	for _, x := range k.peers {
		s += fmt.Sprintf("\n\t%s", x)
	}
	return s
}
func (k *ketama) Type() agro.RingType { return Ketama }
func (k *ketama) Version() int        { return k.version }

func (k *ketama) Marshal() ([]byte, error) {
	var out models.Ring

	out.Version = uint32(k.version)
	out.ReplicationFactor = uint32(k.rep)
	out.Type = uint32(k.Type())
	out.UUIDs = k.peers
	return out.Marshal()
}

func (k *ketama) AddPeers(pl agro.PeerList, mods ...agro.RingModification) (agro.Ring, error) {
	newPeers := k.Members().Union(pl)
	if reflect.DeepEqual(newPeers, k.Members()) {
		return nil, errors.New("no difference in membership")
	}
	newk := &ketama{
		version: k.version + 1,
		rep:     k.rep,
		peers:   newPeers,
		ring:    hashring.New(newPeers),
	}
	for _, x := range mods {
		x.ModifyRing(newk)
	}
	return newk, nil
}

func (k *ketama) RemovePeers(pl agro.PeerList, mods ...agro.RingModification) (agro.Ring, error) {
	newPeers := k.Members().AndNot(pl)
	if reflect.DeepEqual(newPeers, k.Members()) {
		return nil, errors.New("no difference in membership")
	}
	newk := &ketama{
		version: k.version + 1,
		rep:     k.rep,
		peers:   newPeers,
		ring:    hashring.New(newPeers),
	}
	for _, x := range mods {
		x.ModifyRing(newk)
	}
	return newk, nil
}

func (k *ketama) ChangeReplication(r int) {
	k.rep = r
}
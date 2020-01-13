package mempool

import (
	"github.com/HcashOrg/hcd/hcutil"
)

type addrPool struct {
	addrPool map[string]interface{}
}

func (mp *TxPool) FetchRouteAddrPoolState() []string{
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()
	return mp.fetchRouteAddrPoolState()
}

func (mp *TxPool)fetchRouteAddrPoolState()[]string {
	addrSlice:=make([]string,0,len(mp.addrPool.addrPool))
	for addr,_:=range mp.addrPool.addrPool{
		addrSlice = append(addrSlice,addr)
	}

	return addrSlice
}

func (mp *TxPool) AddToAddrPool(addrList []string) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	for _, addr := range addrList {
		err:=mp.maybeAddtoAddrPool(addr)
		if err!=nil{
			log.Errorf("AddToAddrPool addr %v ,err: %v",addr,err)
		}
	}
}

func (mp *TxPool) maybeAddtoAddrPool(addr string) error {
	_, err := hcutil.DecodeAddress(addr)
	if err != nil {
		return err
	}
	mp.addrPool.addrPool[addr] = nil

	return nil
}

func (mp *TxPool) RemoveAddr(addr string) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	mp.removeAddr(addr)
}

func (mp *TxPool) removeAddr(addr string) {
	delete(mp.addrPool.addrPool, addr)
}

func (mp *TxPool) GetAddrList() []string {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()
	return mp.getAddrList()
}

func (mp *TxPool) getAddrList() []string {

	addrSlice := make([]string, 0, len(mp.addrPool.addrPool))

	for addr, _ := range mp.addrPool.addrPool {
		addrSlice = append(addrSlice, addr)
	}
	return addrSlice
}

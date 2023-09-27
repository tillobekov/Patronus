package blockchain

type ManagerSet struct {
	ByKey map[string]*Manager
}

var Managers ManagerSet

func NewManagerSet() ManagerSet {
	Managers = ManagerSet{
		ByKey: make(map[string]*Manager),
	}
	return Managers
}

//func GetManager(key string) Manager {
//	if key == "ETH" {
//		//return NewEthereumManager("HTTP://127.0.0.1:8545")
//		return *Managers.ByKey["ETH"]
//	}
//	return nil
//}

func (ms *ManagerSet) AddManager(manager *Manager, key string) {
	Managers.ByKey[key] = manager
}

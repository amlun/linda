package locker

//import (
//	"fmt"
//	neturl "net/url"
//)
//
//type Locker interface {
//	Connect(url *neturl.URL) error
//	Lock(lockName string) error
//	ReleaseLock(lockName string) error
//	Close() error
//}
//
//// registered lockers
//var lockerRegistery = make(map[string]Locker)
//
//// Register locker with its scheme
//func Register(scheme string, l Locker) {
//	lockerRegistery[scheme] = l
//}
//
//// new instance locker with urlString
//func NewLocker(urlString string) (Locker, error) {
//	// get scheme from uri
//	url, err := neturl.Parse(urlString)
//	if err != nil {
//		return nil, err
//	}
//	scheme := url.Scheme
//	if b, ok := lockerRegistery[scheme]; ok {
//		err := b.Connect(url)
//		if err != nil {
//			return nil, err
//		}
//		return b, nil
//	}
//
//	return nil, fmt.Errorf("Unknow locker scheme [%s]", scheme)
//}

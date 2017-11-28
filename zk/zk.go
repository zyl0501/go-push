// Copyright © 2014 Terry Mao, LiuDing All rights reserved.
// This file is part of gopush-cluster.

// gopush-cluster is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// gopush-cluster is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with gopush-cluster.  If not, see <http://www.gnu.org/licenses/>.

// github.com/samuel/go-zookeeper
// Copyright (c) 2013, Samuel Stauffer <samuel@descolada.com>
// All rights reserved.

package zk

import (
	log "github.com/alecthomas/log4go"
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"strings"
	"path"
)

var (
	// error
	ErrNoChild      = errors.New("zk: children is nil")
	ErrNodeNotExist = errors.New("zk: node not exist")
)

// Connect connect to zookeeper, and start a goroutine log the event.
func connect(addr []string, timeout time.Duration) (*zk.Conn, error) {
	conn, session, err := zk.Connect(addr, timeout)
	if err != nil {
		log.Error("zk.Connect(\"%v\", %d) error(%v)", addr, timeout, err)
		return nil, err
	}
	go func() {
		for {
			event := <-session
			log.Debug("zookeeper get a session event: %s", event.State.String())
		}
	}()
	return conn, nil
}

func getData(conn *zk.Conn, path string) []byte {
	data, _, err := conn.Get(path)
	if err != nil {
		log.Error("getData:%s, err=%v", path, err)
	}else{
		log.Debug("getData success. path=%s, data=%s", path, string(data))
	}
	return data
}

/**
* 获取子节点
*
* @param key
* @return
*/
func getNodes(conn *zk.Conn, path string) ([]string) {
	nodes, stat, err := conn.Children(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil
		}
		log.Error("zk.Children(\"%s\") error(%v)", path, err)
		return nil
	}
	if stat == nil {
		return nil
	}
	if len(nodes) == 0 {
		return nil
	}
	return nodes
}

// GetNodesW get all child from zk path with a watch.
func getNodesW(conn *zk.Conn, path string) ([]string, <-chan zk.Event, error) {
	nodes, stat, watch, err := conn.ChildrenW(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, nil, ErrNodeNotExist
		}
		log.Error("zk.ChildrenW(\"%s\") error(%v)", path, err)
		return nil, nil, err
	}
	if stat == nil {
		return nil, nil, ErrNodeNotExist
	}
	if len(nodes) == 0 {
		return nil, nil, ErrNoChild
	}
	return nodes, watch, nil
}

/**
 * 删除节点
 */
func remove(conn *zk.Conn, path string) error {
	err := conn.Delete(path, -1)
	if err != nil {
		log.Error("removeAndClose:%s, err=%v", path, err)
	}
	return err
}

/**
 * 持久化数据
 *
 * @param key
 * @param value
 */
func registerPersist(conn *zk.Conn, key string, value []byte) error {
	if isExisted(conn, key) {
		return update(conn, key, value);
	} else {
		// create zk root path
		tpath := ""
		for _, str := range strings.Split(key, "/")[:] {
			tpath = path.Join(tpath, "/", str)
			if isExisted(conn, tpath) {
				continue
			}
			log.Debug("create zookeeper path: \"%s\"", tpath)
			_, err := conn.Create(tpath, []byte(""), 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				if err == zk.ErrNodeExists {
					log.Warn("zk.create(\"%s\") exists", tpath)
				} else {
					log.Error("zk.create(\"%s\") error(%v)", tpath, err)
					return err
				}
			}
		}
		return nil
	}
}

/**
 * 注册临时数据
 *
 * @param key
 * @param value
 */
func registerEphemeral(conn *zk.Conn, key string, value []byte, cacheNode bool) error {
	if isExisted(conn, key) {
		return update(conn, key, value);
	} else {
		// create zk root path
		tpath := ""
		for _, str := range strings.Split(key, "/")[:] {
			tpath = path.Join(tpath, "/", str)
			//if isExisted(conn, tpath){
			//	continue
			//}
			log.Debug("create zookeeper path: \"%s\"", tpath)
			_, err := conn.Create(tpath, []byte(""), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			if err != nil {
				if err == zk.ErrNodeExists {
					log.Warn("zk.create(\"%s\") exists", tpath)
				} else {
					log.Error("zk.create(\"%s\") error(%v)", tpath, err)
					return err
				}
			}
		}
		return update(conn, key, value);
	}
}

/**
 * 更新数据
 *
 * @param key
 * @param value
 */
func update(conn *zk.Conn, key string, value []byte) error {
	_, err := conn.Set(key, value, -1)
	if err != nil {
		log.Debug("update:%s, err=%v", key, err)
	} else {
		log.Debug("update success, key=%s, value=%s", key, string(value))
	}
	return err
}

/**
 * 判断路径是否存在
 *
 * @param key
 * @return
 */
func isExisted(conn *zk.Conn, key string) bool {
	exist, _, err := conn.Exists(key)
	if err != nil {
		log.Warn("check node %s exist error.", key)
		return false
	} else {
		return exist
	}
}

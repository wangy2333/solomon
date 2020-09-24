package solomon

import (
	"fmt"
	"strings"
)

type Trie struct {
	part   string           //不是最终，就存part
	next   map[string]*Trie //子节点
	path   string           //最终，存整个路径当前路劲
	isWord bool             //当前节点是否是完整路径
	isWild bool             //当前节点是否是动态通配符
}

func NewTrie() *Trie {
	return &Trie{
		next: make(map[string]*Trie),
	}
}

//insert node to trie
func (t *Trie) insertNode(path string) {
	p := t
	pathSlice := strings.Split(path, "/")
	for _, part := range pathSlice {
		if tn, ok := p.next[part]; ok {
			//存在当前节点,向下继续找
			p = tn
		} else {
			//不存在当前节点，创建一个当前节点。存入map中
			p.next[part] = &Trie{
				next: make(map[string]*Trie),
				path: part,
			}
			p = p.next[part]
		}
	}
	//遍历完毕之后， 将其存入一句话
	p.path = path
	p.isWord = true
}

//传入一个string 找到了就返回当前节点的path
//寻找失败包括找不到，就返回空字符串，error返回找不到
type routeNode struct {
	path   string
	params map[string]string
}

func (t *Trie) searchNode(path string) (*routeNode, error) {
	p := t
	pathSlice := strings.Split(path, "/")
	resNode := &routeNode{params: make(map[string]string)}
	for _, part := range pathSlice {
		//精确查询
		if tn, ok := p.next[part]; ok {
			p = tn
		} else {
			//首先进行模糊匹配
			v, parm, matches := t.matches(p, part)
			if len(parm) != 0 {
				resNode.params[parm] = part
			}
			//如果模糊匹配都找不到，就真的找不到了
			if !matches {
				return nil, fmt.Errorf("is not find path ,and the path is %v", path)
			}
			p = v
		}
	}
	if p.isWord {
		resNode.path = p.path
		return resNode, nil
	} else {
		return nil, fmt.Errorf("is not find path ,and the path is %v", path)
	}
}

//模糊匹配的函数 返回的是参数名字
func (t *Trie) matches(p *Trie, part string) (*Trie, string, bool) {
	for k, v := range p.next {
		//能找到的都break
		if k[0] == ':' { //动态路由
			return v, k[1:], true
		}
		index := strings.Index(k, "*")
		if index == 0 { //通配在前面
			parameter := k[1:]
			isHas := strings.HasSuffix(part, parameter)
			if isHas {
				return v, parameter, true
			}
		}
		if index == len(k)-1 {
			parameter := k[:len(k)-1]
			isHas := strings.HasPrefix(part, parameter)
			if isHas {
				return v, parameter, true
			}
		}
		if index > 0 && index < len(k)-1 {
			prePara := k[:index]
			rerPara := k[index+1:]
			if strings.HasPrefix(part, prePara) && strings.HasSuffix(part, rerPara) {
				return v, "", true
			}
		}
	}
	return nil, "", false
}

package indexer

import (
	"fmt"
	"logger"
	"os"
)

var count int

type TreeNode struct {
	word      string
	frequency int
	sources   []string
	color     string
	lchild    *TreeNode
	rchild    *TreeNode
	parent    *TreeNode
}

type IndexCache struct {
	root   *TreeNode
	cur    *TreeNode
	create *TreeNode
	source string
}

func (this *IndexCache) Add(word string, source string) {
	this.create = new(TreeNode)
	this.create.word = word
	this.create.frequency = 1
	this.create.sources = []string{source}
	this.create.color = "red"

	if !this.IsEmpty() {
		this.cur = this.root
		for {
			if word < this.cur.word {
				//如果要插入的值比当前节点的值小，则当前节点指向当前节点的左孩子，如果
				//左孩子为空，就在这个左孩子上插入新值
				if this.cur.lchild == nil {
					this.cur.lchild = this.create
					this.create.parent = this.cur
					break
				} else {
					this.cur = this.cur.lchild
				}

			} else if word > this.cur.word {
				//如果要插入的值比当前节点的值大，则当前节点指向当前节点的右孩子，如果
				//右孩子为空，就在这个右孩子上插入新值
				if this.cur.rchild == nil {
					this.cur.rchild = this.create
					this.create.parent = this.cur
					break
				} else {
					this.cur = this.cur.rchild
				}

			} else {
				//如果要插入的值在树中已经存在
				this.cur.frequency++
				this.cur.sources = append(this.cur.sources, source)
				return
			}
		}

	} else {
		this.root = this.create
		this.root.color = "black"
		this.root.parent = nil
		return
	}

	//插入节点后对红黑性质进行修复
	this.insertBalanceFixup(this.create)
}

func (this *IndexCache) Delete(word string) {
	var (
		deleteNode func(node *TreeNode)
		node       *TreeNode = this.Search(word)
		parent     *TreeNode
		revise     string
	)

	if node == nil {
		return
	} else {
		parent = node.parent
	}

	//下面这小块代码用来判断替代被删节点位置的节点是哪个后代
	if node.lchild == nil && node.rchild == nil {
		revise = "none"
	} else if parent == nil {
		revise = "root"
	} else if node == parent.lchild {
		revise = "left"
	} else if node == parent.rchild {
		revise = "right"
	}

	deleteNode = func(node *TreeNode) {
		if node == nil {
			return
		}

		if node.lchild == nil && node.rchild == nil {
			//如果要删除的节点没有孩子，直接删掉它就可以(毫无挂念~.~!)
			if node == this.root {
				this.root = nil
			} else {
				if node.parent.lchild == node {
					node.parent.lchild = nil
				} else {
					node.parent.rchild = nil
				}
			}

		} else if node.lchild != nil && node.rchild == nil {
			//如果要删除的节点只有左孩子或右孩子，让这个节点的父节点指向它的指针指向它的
			//孩子即可
			if node == this.root {
				node.lchild.parent = nil
				this.root = node.lchild
			} else {
				node.lchild.parent = node.parent
				if node.parent.lchild == node {
					node.parent.lchild = node.lchild
				} else {
					node.parent.rchild = node.lchild
				}
			}

		} else if node.lchild == nil && node.rchild != nil {
			if node == this.root {
				node.rchild.parent = nil
				this.root = node.rchild
			} else {
				node.rchild.parent = node.parent
				if node.parent.lchild == node {
					node.parent.lchild = node.rchild
				} else {
					node.parent.rchild = node.rchild
				}
			}

		} else {
			//如果要删除的节点既有左孩子又有右孩子，就把这个节点的直接后继的值赋给这个节
			//点，然后删除直接后继节点即可
			successor := this.GetSuccessor(node.word)
			node.word = successor.word
			node.color = successor.color
			deleteNode(successor)
		}
	}

	deleteNode(node)
	if node.color == "black" {
		if revise == "root" {
			this.deleteBalanceFixup(this.root)
		} else if revise == "left" {
			this.deleteBalanceFixup(parent.lchild)
		} else if revise == "right" {
			this.deleteBalanceFixup(parent.rchild)
		}
	}
}

//这个函数用于在红黑树执行插入操作后，修复红黑性质
func (this *IndexCache) insertBalanceFixup(insertnode *TreeNode) {
	var uncle *TreeNode

	for insertnode.color == "red" && insertnode.parent.color == "red" {
		//获取新插入的节点的叔叔节点(与父节点同根的另一个节点)
		if insertnode.parent == insertnode.parent.parent.lchild {
			uncle = insertnode.parent.parent.rchild
		} else {
			uncle = insertnode.parent.parent.lchild
		}

		if uncle != nil && uncle.color == "red" {
			uncle.color, insertnode.parent.color = "black", "black"
			insertnode = insertnode.parent.parent
			if insertnode == this.root || insertnode == nil {
				return
			}
			insertnode.color = "red"

		} else {
			if insertnode.parent == insertnode.parent.parent.lchild {
				if insertnode == insertnode.parent.rchild {
					insertnode = insertnode.parent
					this.LeftRotate(insertnode)
				}
				insertnode = insertnode.parent
				insertnode.color = "black"
				insertnode = insertnode.parent
				insertnode.color = "red"
				this.RightRotate(insertnode)

			} else {
				if insertnode == insertnode.parent.lchild {
					insertnode = insertnode.parent
					this.RightRotate(insertnode)
				}
				insertnode = insertnode.parent
				insertnode.color = "black"
				insertnode = insertnode.parent
				insertnode.color = "red"
				this.LeftRotate(insertnode)
			}
			return
		}
	}
}

//这个函数用于在红黑树执行删除操作后，修复红黑性质
func (this *IndexCache) deleteBalanceFixup(node *TreeNode) {
	var brother *TreeNode

	for node != this.root && node.color == "black" {
		if node.parent.lchild == node && node.parent.rchild != nil {
			brother = node.parent.rchild
			if brother.color == "red" {
				brother.color = "black"
				node.parent.color = "red"
				this.LeftRotate(node.parent)
			} else if brother.color == "black" && brother.lchild != nil && brother.lchild.color == "black" && brother.rchild != nil && brother.rchild.color == "black" {
				brother.color = "red"
				node = node.parent
			} else if brother.color == "black" && brother.lchild != nil && brother.lchild.color == "red" && brother.rchild != nil && brother.rchild.color == "black" {
				brother.color = "red"
				brother.lchild.color = "black"
				this.RightRotate(brother)
			} else if brother.color == "black" && brother.rchild != nil && brother.rchild.color == "red" {
				brother.color = "red"
				brother.rchild.color = "black"
				brother.parent.color = "black"
				this.LeftRotate(brother.parent)
				node = this.root
			}

		} else if node.parent.rchild == node && node.parent.lchild != nil {
			brother = node.parent.lchild
			if brother.color == "red" {
				brother.color = "black"
				node.parent.color = "red"
				this.RightRotate(node.parent)
			} else if brother.color == "black" && brother.lchild != nil && brother.lchild.color == "black" && brother.rchild != nil && brother.rchild.color == "black" {
				brother.color = "red"
				node = node.parent
			} else if brother.color == "black" && brother.lchild != nil && brother.lchild.color == "black" && brother.rchild != nil && brother.rchild.color == "red" {
				brother.color = "red"
				brother.rchild.color = "black"
				this.LeftRotate(brother)
			} else if brother.color == "black" && brother.lchild != nil && brother.lchild.color == "red" {
				brother.color = "red"
				brother.lchild.color = "black"
				brother.parent.color = "black"
				this.RightRotate(brother.parent)
				node = this.root
			}
		} else {
			return
		}
	}
}

func (this IndexCache) GetRoot() *TreeNode {
	if this.root != nil {
		return this.root
	}
	return nil
}

func (this IndexCache) IsEmpty() bool {
	if this.root == nil {
		return true
	}
	return false
}

func (this IndexCache) InOrderTravel() {
	var inOrderTravel func(node *TreeNode)

	inOrderTravel = func(node *TreeNode) {
		if node != nil {
			inOrderTravel(node.lchild)
			fmt.Printf("%s %d\r\n", node.word, node.frequency)
			inOrderTravel(node.rchild)
		}
	}

	inOrderTravel(this.root)
}

func (this IndexCache) writeToFile(node *TreeNode, index_file *os.File) {
	if node != nil {
		if node.lchild != nil {
			this.writeToFile(node.lchild, index_file)
		}
		fmt.Fprintf(index_file, "%s %d %s\r\n", node.word, node.frequency, this.source)
		if node.rchild != nil {
			this.writeToFile(node.rchild, index_file)
		}
	}
}

func (this IndexCache) WriteToFile(file string) error {
	index_file, err := os.Create(file)
	if err != nil {
		logger.Log("生成临时索引文件失败")
		return err
	}
	defer func() {
		if err = index_file.Close(); err != nil {
			logger.Log("关闭临时索引文件:'" + file + "'失败")
		}
	}()
	this.writeToFile(this.root, index_file)
	return nil
}

func (this IndexCache) Search(word string) *TreeNode {
	//和Add操作类似，只要按照比当前节点小就往左孩子上拐，比当前节点大就往右孩子上拐的思路
	//一路找下去，知道找到要查找的值返回即可
	this.cur = this.root
	for {
		if this.cur == nil {
			return nil
		}

		if word < this.cur.word {
			this.cur = this.cur.lchild
		} else if word > this.cur.word {
			this.cur = this.cur.rchild
		} else {
			return this.cur
		}
	}
}

func (this IndexCache) GetDeepth() int {
	var getDeepth func(node *TreeNode) int

	getDeepth = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		if node.lchild == nil && node.rchild == nil {
			return 1
		}
		var (
			ldeepth int = getDeepth(node.lchild)
			rdeepth int = getDeepth(node.rchild)
		)
		if ldeepth > rdeepth {
			return ldeepth + 1
		} else {
			return rdeepth + 1
		}
	}

	return getDeepth(this.root)
}

func (this IndexCache) GetMin() string {
	//根据二叉查找树的性质，树中最左边的节点就是值最小的节点
	if this.root == nil {
		return ""
	}
	this.cur = this.root
	for {
		if this.cur.lchild != nil {
			this.cur = this.cur.lchild
		} else {
			return this.cur.word
		}
	}
}

func (this IndexCache) GetMax() string {
	//根据二叉查找树的性质，树中最右边的节点就是值最大的节点
	if this.root == nil {
		return ""
	}
	this.cur = this.root
	for {
		if this.cur.rchild != nil {
			this.cur = this.cur.rchild
		} else {
			return this.cur.word
		}
	}
}

func (this IndexCache) GetPredecessor(word string) *TreeNode {
	getMax := func(node *TreeNode) *TreeNode {
		if node == nil {
			return nil
		}
		for {
			if node.rchild != nil {
				node = node.rchild
			} else {
				return node
			}
		}
	}

	node := this.Search(word)
	if node != nil {
		if node.lchild != nil {
			//如果这个节点有左孩子，那么它的直接前驱就是它左子树的最右边的节点，因为比这
			//个节点值小的节点都在左子树，而这些节点中值最大的就是这个最右边的节点
			return getMax(node.lchild)
		} else {
			//如果这个节点没有左孩子，那么就沿着它的父节点找，知道某个父节点的父节点的右
			//孩子就是这个父节点，那么这个父节点的父节点就是直接前驱
			for {
				if node == nil || node.parent == nil {
					break
				}
				if node == node.parent.rchild {
					return node.parent
				}
				node = node.parent
			}
		}
	}

	return nil
}

func (this IndexCache) GetSuccessor(word string) *TreeNode {
	getMin := func(node *TreeNode) *TreeNode {
		if node == nil {
			return nil
		}
		for {
			if node.lchild != nil {
				node = node.lchild
			} else {
				return node
			}
		}
	}

	//参照寻找直接前驱的函数对比着看
	node := this.Search(word)
	if node != nil {
		if node.rchild != nil {
			return getMin(node.rchild)
		} else {
			for {
				if node == nil || node.parent == nil {
					break
				}
				if node == node.parent.lchild {
					return node.parent
				}
				node = node.parent
			}
		}
	}

	return nil
}

func (this *IndexCache) Clear() {
	this.root = nil
	this.cur = nil
	this.create = nil
}

func (this *IndexCache) LeftRotate(node *TreeNode) {
	if node.rchild == nil {
		return
	}

	right_child := node.rchild
	//将要旋转的节点的右孩子的左孩子赋给这个节点的右孩子，这里最好按如下3行代码的顺序写，
	//否则该节点的右孩子的左孩子为nil时，很容易忘记把这个节点的右孩子也置为nil.
	node.rchild = right_child.lchild
	if node.rchild != nil {
		node.rchild.parent = node
	}

	//让要旋转的节点的右孩子的父节点指针指向当前节点父节点。如果父节点是根节点要特别处理
	right_child.parent = node.parent
	if node.parent == nil {
		this.root = right_child
	} else {
		if node.parent.lchild == node {
			node.parent.lchild = right_child
		} else {
			node.parent.rchild = right_child
		}
	}

	//上面的准备工作完毕了，就可以开始旋转了，让要旋转的节点的右孩子的左孩子指向该节点，
	//别忘了把这个被旋转的节点的父节点指针指向新的父节点
	right_child.lchild = node
	node.parent = right_child
}

func (this *IndexCache) RightRotate(node *TreeNode) {
	//向右旋转的过程与LeftRotate()正好相反
	if node.lchild == nil {
		return
	}

	left_child := node.lchild
	node.lchild = left_child.rchild
	if node.lchild != nil {
		node.lchild.parent = node
	}

	left_child.parent = node.parent
	if node.parent == nil {
		this.root = left_child
	} else {
		if node.parent.lchild == node {
			node.parent.lchild = left_child
		} else {
			node.parent.rchild = left_child
		}
	}
	left_child.rchild = node
	node.parent = left_child
}

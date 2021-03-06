package answer

import (
	"regexp"
	"errors"
	"strings"
)

const BLANK_STRING = "___"

type Find struct {
	BlankId     int
	PreString   string // 前导诗句
	BlankString string // 本空
	PostString  string // 后置诗句
	BlankNum    int
	BlankFinish bool //是否是完全块
}

//find the answer by contents.
//findSrc : blank is ___
func FindTheAnswer(contents string, findSrc string) (string, error) {
	finds := make([]Find, 0, 0)
	findSplits, findPun := SplitByPunctuation(findSrc)

	if len(findSplits) != len(findPun) {
		return "", errors.New("标点符号要输入完整，才是一句话。")
	}

	var j = 0
	for i := range findSplits {
		if findSplits[i] == BLANK_STRING {
			var newFind Find
			newFind.BlankId = j
			newFind.BlankNum = i
			//前一句
			if i > 0 {
				newFind.PreString = findSplits[i-1]
			}
			//后一句
			if i < len(findSplits)-1 {
				newFind.PostString = findSplits[i+1]
			}

			finds = append(finds, newFind)
			j++
		}
	}

	//构造答案
	makeAnswer(contents, finds)

	//Return答案
	var returnStr string
	for i := range findSplits {
		if findSplits[i] == BLANK_STRING {
			ans := GetFindsByBlankNum(finds, i).BlankString
			if len(ans) == 0 {
				ans = "未找到本空答案"
			}
			returnStr += "<ans>" + ans + "</ans>"
		} else {
			returnStr += findSplits[i]
		}
		//标点
		returnStr += findPun[i]
	}
	return returnStr, nil
}

func GetFindsByBlankNum(finds []Find, j int) (Find) {
	for i := range finds {
		if finds[i].BlankNum == j {
			return finds[i]
		}
	}
	return Find{}
}

func makeAnswer(contents string, finds []Find) {
	//先分割contents
	contentsSplit, _ := SplitByPunctuation(contents)

	for i := range finds {
		//自主块：可自主推导出答案的Find
		//完全块：PreString BlankString PostString 均有的Find
		//推导所有自主块的答案 - 让它们变成完全块
		if checkStringNotBlank(finds[i].PreString) || checkStringNotBlank(finds[i].PostString) {
			rootFind := &finds[i]
			if checkStringNotBlank(rootFind.PreString) {
				makeWithPreContent(contentsSplit, rootFind)
			} else {
				makeWithPostContent(contentsSplit, rootFind)
			}
		}
	}

	// 遍历完全块，找它们左右是否有非自主块。给它们添加答案
	for i := range finds {
		if finds[i].BlankFinish {
			root := i
			//前向搜索
			if root > 0 {
				for i := root - 1; i >= 0; i-- {
					if finds[i].BlankFinish {
						break
					}
					finds[i].BlankString = finds[i+1].PreString
					finds[i].PostString = finds[i+1].BlankString

					//find finds[i].PreString
					makeWithPostContent(contentsSplit, &finds[i])
				}
			}
			//后向搜索
			if root < len(finds)-1 {
				for i := root + 1; i < len(finds); i++ {
					if finds[i].BlankFinish {
						break
					}
					finds[i].BlankString = finds[i-1].PostString
					finds[i].PreString = finds[i-1].BlankString

					//find finds[i].PostString
					makeWithPreContent(contentsSplit, &finds[i])
				}
			}
		}
	}
}

func checkStringNotBlank(check string) bool {
	return len(check) > 0 && check != BLANK_STRING
}

//已知newFind的PreString，求BlankString和PostString
func makeWithPreContent(contentsSplit []string, newFind *Find) {
	for l := range contentsSplit {
		if isEqual(contentsSplit[l], newFind.PreString) && l < len(contentsSplit)-1 {
			newFind.BlankString = contentsSplit[l+1]
			if l < len(contentsSplit)-2 {
				newFind.PostString = contentsSplit[l+2]
			}
			newFind.BlankFinish = true
		}
	}
}

//已知newFind的PostString，求BlankString和PreString
func makeWithPostContent(contentsSplit []string, newFind *Find) {
	for l := range contentsSplit {
		if isEqual(contentsSplit[l], newFind.PostString) && l > 0 {
			newFind.BlankString = contentsSplit[l-1]
			if l-1 > 0 {
				newFind.PreString = contentsSplit[l-2]
			}
			newFind.BlankFinish = true
		}
	}
}

// 判断要找的句子，和诗文内的句子是否相等。
func isEqual(content, find string) bool {
	if content == find {
		return true
	}

	//去掉书名号再判断，也就是说：《齐谐》者 == 齐谐者
	if strings.Contains(content,"《") {
		regX := regexp.MustCompile("[《》]")
		newContent := regX.ReplaceAllString(content,"")

		if newContent == find{
			return true
		}
	}

	return false
}

// 按标点符号分隔句子
func SplitByPunctuation(s string) ([]string, []string) {
	regPunctuation, _ := regexp.Compile(`[,，。.?？！!;；：:]`)
	//匹配标点符号，保存下来。 然后分割字符串
	toPun := regPunctuation.FindAllString(s, -1)
	result := regPunctuation.Split(s, -1)

	if len(result[len(result)-1]) == 0 {
		result = result[:len(result)-1]
	}

	//去掉前后空格，去掉引号
	for i := range result {
		result[i] = strings.TrimSpace(result[i])
		regQuoting := regexp.MustCompile("[“”‘’']")
		result[i] = regQuoting.ReplaceAllString(result[i], "")
	}

	return result, toPun
}

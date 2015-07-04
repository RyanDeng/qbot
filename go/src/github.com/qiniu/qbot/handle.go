package qbot

import (
	"fmt"
	"strings"
)

type Handle interface {
	KeyWords() []string
	ThinkOut(srv *Service, msg string) string
	GroupThinkOut(srv *Service, msg string) string
}

//------

type ReminderHandle struct {
}

func (r ReminderHandle) KeyWords() []string {
	return []string{"提醒", "会议"}
}

func (r ReminderHandle) ThinkOut(srv *Service, msg string) string {
	// message := strings.Trim(msg, "提醒")
	// commands := strings.Split(msg, " ")
	// date := commands[0]
	// time := commands[1]
	// event := commands[2]
	return ""
}

func (r ReminderHandle) GroupThinkOut(srv *Service, msg string) string {
	return ""
}

//--------------
type ContactHandle struct {
}

func (c ContactHandle) KeyWords() []string {
	return []string{"联系", "联系方式", "电话", "信息", "部门", "查找", "是谁", "找"}
}

func (c ContactHandle) ThinkOut(srv *Service, msg string) string {
	message := strings.TrimSpace(msg)
	message = strings.Trim(message, "，")
	message = strings.Trim(message, "。")
	message = strings.Trim(message, "联系方式")
	message = strings.Trim(message, "电话")
	message = strings.Trim(message, "信息")
	message = strings.Trim(message, "部门")
	message = strings.Trim(message, "查找")
	message = strings.Trim(message, "是谁")
	message = strings.Trim(message, "找")
	message = strings.TrimPrefix(message, "求")
	message = strings.TrimSuffix(message, "的")

	fmt.Println("message going to be searhed", message)
	if len(message) > 0 {
		contacts, err := srv.contactTbl.SearchByAllName(message)
		if err != nil {
			return "Ooops, 我好像短路了。。。能否稍候来找我"
		}
		if len(contacts) == 0 {
			return "sorry, 貌似查不到你要找的人, 请再输入些线索"
		}
		contact := contacts[0]
		if strings.Contains(msg, "电话") || strings.Contains(msg, "联系") {
			return fmt.Sprintf("%v的电话是%d", message, contact.Phone)
		} else if strings.Contains(msg, "部门") {
			return fmt.Sprintf("%v的部门是%v", message, contact.Department)
		} else {
			return fmt.Sprintf("%v, 人称: %v, 电话%v, 部门:%v, 邮箱:%v", contact.Name, contact.NickName, contact.Phone, contact.Department, contact.Email)
		}

	}
	return "不好意思,我不太懂你的意思, 请再输入些线索"
}

func (c ContactHandle) GroupThinkOut(srv *Service, msg string) string {
	message := strings.TrimSpace(msg)
	message = strings.Trim(message, "@QBot")
	message = strings.TrimSpace(message)
	message = strings.Trim(message, "，")
	message = strings.Trim(message, "。")
	message = strings.Trim(message, "联系方式")
	message = strings.Trim(message, "电话")
	message = strings.Trim(message, "信息")
	message = strings.Trim(message, "部门")
	message = strings.Trim(message, "查找")
	message = strings.Trim(message, "是谁")
	message = strings.Trim(message, "找")
	message = strings.TrimPrefix(message, "求")
	message = strings.TrimSuffix(message, "的")
	fmt.Println("message going to be searhed", message)
	if len(message) > 0 {
		contacts, err := srv.contactTbl.SearchByAllName(message)
		if err != nil {
			return "Ooops, 我好像短路了。。。能否稍候来找我"
		}
		if len(contacts) == 0 {
			return "sorry, 貌似查不到你要找的人, 请再输入些线索"
		}
		contact := contacts[0]
		if strings.Contains(msg, "电话") || strings.Contains(msg, "联系") {
			return fmt.Sprintf("%v的电话是%d", message, contact.Phone)
		} else if strings.Contains(msg, "部门") {
			return fmt.Sprintf("%v的部门是%v", message, contact.Department)
		} else {
			return fmt.Sprintf("%v, 人称: %v, 电话%v, 部门:%v, 邮箱:%v", contact.Name, contact.NickName, contact.Phone, contact.Department, contact.Email)
		}

	}
	return "不好意思,我不太懂你的意思, 请再输入些线索"
}

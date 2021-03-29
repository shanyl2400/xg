package service

import (
	"context"
	"fmt"
	"strings"
	"xg/entity"

	"github.com/tealeg/xlsx"
)

func studentStatusName(status int) string {
	switch status {
	case entity.StudentCreated:
		return "新订单"
	case entity.StudentConflictFailed:
		return "无效订单"
	case entity.StudentConflictSuccess:
		return "有效订单"
	}
	return "无效状态"
}

// 序号，学生姓名，联系方式，咨询课程（录单显示也改成咨询课程），居住地址，录单时间（24小时制），状态（有效，无效）
func StudentsToXlsx(ctx context.Context, studentList []*entity.StudentInfo) (*xlsx.File, error) {
	var cell *xlsx.Cell
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("名单列表")
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	row := sheet.AddRow()
	row.SetHeightCM(0.5)
	cell = row.AddCell()
	cell.Value = "录单人"
	cell = row.AddCell()
	cell.Value = "录单时间"
	cell = row.AddCell()
	cell.Value = "学生姓名"
	cell = row.AddCell()
	cell.Value = "联系电话"
	cell = row.AddCell()
	cell.Value = "推荐科目"
	cell = row.AddCell()
	cell.Value = "订单来源"
	cell = row.AddCell()
	cell.Value = "状态"

	for i := range studentList {
		// 判断是否为-9999，是的变为0.0
		var row1 *xlsx.Row

		row1 = sheet.AddRow()
		row1.SetHeightCM(0.5)

		cell = row1.AddCell()
		cell.Value = studentList[i].AuthorName
		cell = row1.AddCell()
		cell.Value = studentList[i].CreatedAt.Format(TIME_FORMAT)
		cell = row1.AddCell()
		cell.Value = studentList[i].Name
		cell = row1.AddCell()
		cell.Value = studentList[i].Telephone
		cell = row1.AddCell()
		cell.Value = strings.Join(studentList[i].IntentSubject, ",")

		cell = row1.AddCell()
		cell.Value = studentList[i].OrderSourceName
		cell = row1.AddCell()
		cell.Value = studentStatusName(studentList[i].Status)
	}

	return file, nil
}

package service

import (
	"context"
	"fmt"
	"strings"
	"xg/entity"

	"github.com/tealeg/xlsx"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

func orderStatusName(status int) string {
	switch status {
	case entity.OrderStatusCreated:
		return "未报名"
	case entity.OrderStatusSigned:
		return "已报名"
	case entity.OrderStatusRevoked:
		return "已退费"
	case entity.OrderStatusInvalid:
		return "无效"
	case entity.OrderStatusDeposit:
		return "付定金"
	}
	return "无效状态"
}

func OrdersToXlsx(ctx context.Context, orderList []*entity.OrderInfoWithRecords) (*xlsx.File, error) {
	var cell *xlsx.Cell
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("派单列表")
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
	cell.Value = "派单时间"
	cell = row.AddCell()
	cell.Value = "学生姓名"
	cell = row.AddCell()
	cell.Value = "推荐机构"
	cell = row.AddCell()
	cell.Value = "联系电话"
	cell = row.AddCell()
	cell.Value = "推荐科目"
	cell = row.AddCell()
	cell.Value = "订单来源"
	cell = row.AddCell()
	cell.Value = "状态"
	cell = row.AddCell()
	cell.Value = "报名费用"

	for i := range orderList {
		// 判断是否为-9999，是的变为0.0
		var row1 *xlsx.Row

		row1 = sheet.AddRow()
		row1.SetHeightCM(0.5)

		cell = row1.AddCell()
		cell.Value = orderList[i].AuthorName
		cell = row1.AddCell()
		cell.Value = orderList[i].StudentSummary.CreatedAt.Format(TIME_FORMAT)
		cell = row1.AddCell()
		cell.Value = orderList[i].CreatedAt.Format(TIME_FORMAT)
		cell = row1.AddCell()
		cell.Value = orderList[i].StudentSummary.Name
		cell = row1.AddCell()
		cell.Value = orderList[i].OrgName

		cell = row1.AddCell()
		cell.Value = orderList[i].StudentSummary.Telephone
		cell = row1.AddCell()
		cell.Value = strings.Join(orderList[i].IntentSubject, ",")
		cell = row1.AddCell()
		cell.Value = orderList[i].OrderSourceName
		cell = row1.AddCell()
		cell.Value = orderStatusName(orderList[i].Status)

		amount := float64(0)
		for j := range orderList[i].PaymentInfo {
			if orderList[i].PaymentInfo[j].Status != entity.OrderPayStatusChecked {
				continue
			}
			if orderList[i].PaymentInfo[j].Mode == entity.OrderPayModePay {
				amount = amount + orderList[i].PaymentInfo[j].Amount
			} else {
				amount = amount - orderList[i].PaymentInfo[j].Amount
			}
		}
		cell = row1.AddCell()
		cell.Value = fmt.Sprintf("%v", amount)

	}

	return file, nil
}

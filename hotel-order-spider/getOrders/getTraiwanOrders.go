package getOrders

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GetTraiwanOrderResponseBody struct {
	XMLName xml.Name `xml:"response"`
	Orders  struct {
		Order []struct {
			ID     string `xml:"id"`
			Person struct {
				Name string `xml:"name"`
			} `xml:"person"`
			Source            string `xml:"source"`
			Transaction_price string `xml:"transaction_price"`
			Room_reservations []struct {
				Room_type struct {
					Id   string `xml:"id"`
					Name string `xml:"name"`
				} `xml:"room_type"`
				Date string `xml:"date"`
			} ` xml:"room_reservations>room_reservation"`
			Delete_status  int    `xml:"delete_status"`
			Generated_time string `xml:"generated_time"`
		} `xml:"order"`
	} `xml:"orders"`
}

type RoomInfo struct {
	RoomType  string
	StartDate string
	Days      int
}

func GetTraiwan(platform map[string]interface{}, dateFrom, dateTo string) {
	var url string
	var result string

	hotels, ok := platform["hotel"].([]interface{})
	if !ok || hotels == nil {
		fmt.Println("無法取得 hotel")
	}

	for _, hotelRaw := range hotels {
		hotel, ok := hotelRaw.(map[string]interface{})
		if !ok || hotel == nil {
			fmt.Println("無法取得 hotel")
			continue
		}
		hotelName, _ := hotel["name"].(string)
		cookie, _ := hotel["cookie"].(string)
		hotelId, _ := hotel["hotelid"].(string)
		fmt.Printf("Hotel Name: %s, Hotel ID: %s\n", hotelName, hotelId)

		url = "https://traiwan.com/place/accommodation/butler/order/search/ajax/search.php"
		rawbody := `criteria_xml=%3Ccriteria%3E%3Cid%3E%3C%2Fid%3E%3Cname%3E%3C%2Fname%3E%3Cphone%3E%3C%2Fphone%3E%3Cemail%3E%3C%2Femail%3E%3Cbirthday%3E%3C%2Fbirthday%3E%3Cssn%3E%3C%2Fssn%3E%3Cnotice%3E%3C%2Fnotice%3E%3Cstay_date%3E%3Cbeginning_date%3E` + dateFrom + `%3C%2Fbeginning_date%3E%3Cending_date%3E` + dateTo + `%3C%2Fending_date%3E%3C%2Fstay_date%3E%3Creservation_date%3E%3Cbeginning_date%3E%3C%2Fbeginning_date%3E%3Cending_date%3E%3C%2Fending_date%3E%3C%2Freservation_date%3E%3Croom_types%3E%3C%2Froom_types%3E%3Cstatus%3E%3C%2Fstatus%3E%3Corder_filter%3ENON%3C%2Forder_filter%3E%3Cprice%3E%3Cbeginning_price%3E%3C%2Fbeginning_price%3E%3Cending_price%3E%3C%2Fending_price%3E%3C%2Fprice%3E%3Cdown_payment_type%3E%3C%2Fdown_payment_type%3E%3Csource%3E%3C%2Fsource%3E%3C%2Fcriteria%3E&page=1&rows_per_page=1000`
		if err := DoRequestAndGetResponse_trai("POST", url, strings.NewReader(rawbody), cookie, &result); err != nil {
			fmt.Println("DoRequestAndGetResponse failed!")
			fmt.Println("err", err)
			return
		}

		var ordersData GetTraiwanOrderResponseBody
		err := xml.Unmarshal([]byte(result), &ordersData)
		if err != nil {
			fmt.Println("xml解碼錯誤:", err)
			return
		}

		fmt.Println("ordersData", ordersData)

		var resultData []ReservationsDB
		var data ReservationsDB
		for _, reservation := range ordersData.Orders.Order {
			data.BookingId = reservation.ID
			data.GuestName = reservation.Person.Name

			roomInfoData := make(map[string]*RoomInfo)
			for _, roomReservation := range reservation.Room_reservations {
				roomType := roomReservation.Room_type.Name
				date := roomReservation.Date

				roomInfo, ok := roomInfoData[roomType]
				if !ok {
					roomInfo = &RoomInfo{
						RoomType:  roomType,
						StartDate: date,
						Days:      1,
					}
					roomInfoData[roomType] = roomInfo
				} else {
					roomInfo.Days++
				}
			}
			var combinedRoomInfo string
			for _, roomInfo := range roomInfoData {
				arrivalTime, err := time.Parse("2006-01-02", roomInfo.StartDate)
				if err != nil {
					fmt.Println("Error parsing arrival time:", err)
				}
				if combinedRoomInfo != "" {
					combinedRoomInfo += "+"
				}
				combinedRoomInfo += roomInfo.RoomType

				departureTime := arrivalTime.Add(time.Duration(roomInfo.Days) * 24 * time.Hour)
				checkOutTime := departureTime
				checkInTime := arrivalTime
				data.CheckOutDate = checkOutTime.Format("2006-01-02")
				data.CheckInDate = checkInTime.Format("2006-01-02")
				data.RoomNights = int64(roomInfo.Days)

			}
			numRooms := len(roomInfoData)
			data.NumOfRooms = int64(numRooms)

			parsedTime, err := time.Parse("2006-01-02 15:04:05", reservation.Generated_time)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}

			resultTimeStr := parsedTime.Format("2006-01-02")
			data.BookDate = resultTimeStr

			floatNum, _ := strconv.ParseFloat(reservation.Transaction_price, 64)
			data.Price = floatNum
			if reservation.Delete_status == 1 {
				data.ReservationStatus = "已取消"
			} else {
				data.ReservationStatus = "已成立"
			}

			data.Platform = reservation.Source
			data.Currency = "TWD"
			data.HotelId = hotelId

			if data.Platform != "BACK_END" && data.Platform != "CTRIP_CM" && data.Platform != "BOOKING" && data.Platform != "EXPEDIA" {
				resultData = append(resultData, data)
			}
		}
		fmt.Println("resultdata", resultData)
		time.Sleep(5 * time.Second)

		resultDataJSON, err := json.Marshal(resultData)
		if err != nil {
			fmt.Println("JSON 轉換錯誤:", err)
			return
		}

		var resultDB string
		// 將資料存入DB
		apiurl := "http://149.28.24.90:8893/revenue_reservation/setParseHtmlToDB"
		if err := DoRequestAndGetResponse("POST", apiurl, bytes.NewBuffer(resultDataJSON), cookie, &resultDB); err != nil {
			fmt.Println("setParseHtmlToDB failed!")
			return
		}

	}
}

func DoRequestAndGetResponse_trai(method, postUrl string, reqBody io.Reader, cookie string, resBody any) error {
	req, err := http.NewRequest(method, postUrl, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := &http.Client{Timeout: 40 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	// resBody of type *string is for html
	switch resBody := resBody.(type) {
	case *string:
		// If resBody is a string
		resBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		*resBody = string(resBytes)
	default:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, resBody); err != nil {
			return err
		}
	}

	defer resp.Body.Close()

	return nil
}
